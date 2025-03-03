package main

import (
	"fmt"
	"io"
	"log"
	"main/plugininterface"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

//go build -buildmode=plugin -o plugins/filemanager/filemanager.so plugins/filemanager/filemanager.go

// FileManagerPlugin implementiert das Plugin-Interface.

//WEBDAV operations:
/*
GET	Datei herunterladen.
PUT	Datei hochladen oder überschreiben.
DELETE	Datei oder Verzeichnis löschen.
MKCOL	Neues Verzeichnis erstellen.
PROPFIND	Metadaten oder Verzeichnisstruktur abrufen.
PROPPATCH	Metadaten ändern (optional, z. B. Benutzerrechte).
MOVE	Datei oder Verzeichnis verschieben (optional).
COPY	Datei oder Verzeichnis kopieren (optional).
OPTIONS	Prüfen, welche Methoden unterstützt werden.
LOCK	Datei sperren (für gleichzeitige Zugriffe, optional).
UNLOCK	Datei- oder Verzeichnissperre aufheben.
*/

type FileManagerPlugin struct{}

const basePath = "./webdav" // TODO ./files/{user} and Auth

// XML STruktur für Propfind:
/*type Multistatus struct {
	XMLName  xml.Name   `xml:"D:multistatus"`
	XMLNS    string     `xml:"xmlns:D,attr"`
	Response []Response `xml:"D:response"`
}*/

type Response struct {
	Href     string   `xml:"D:href"`
	Propstat Propstat `xml:"D:propstat"`
}

type Propstat struct {
	Prop   Prop   `xml:"D:prop"`
	Status string `xml:"D:status"`
}

type Prop struct {
	DisplayName   string `xml:"D:displayname,omitempty"`
	ResourceType  string `xml:"D:resourcetype,omitempty"`
	ContentLength int64  `xml:"D:getcontentlength,omitempty"`
	LastModified  string `xml:"D:getlastmodified,omitempty"`
}

func (p FileManagerPlugin) Metadata() plugininterface.PluginMetadata {
	return plugininterface.PluginMetadata{
		Name:          "File Manager",
		Description:   "Ein Plugin zum Verwalten von Dateien",
		Path:          "/plugins/filemanager.html",
		MenuButton:    true,
		UsingDatabase: true,
	}
}

func (p FileManagerPlugin) Routes() []plugininterface.Route {
	log.Println("try to register routes")
	//mux.HandleFunc("/files", listFilesHandler)
	return []plugininterface.Route{
		//{Path: "/files", Handler: listFilesHandler},
		{Path: "/webdav", Handler: fileHandler},
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PROPFIND":
		handlePropfind(w, r) //TODO: Ausgabeformat beachten!
	case "OPTIONS":
		handleOptions(w, r)
	case "GET":
		handleDownload(w, r)
	case "PUT":
		handleUpload(w, r) //TODO: limit Datasize, check http statuscodes, atomic write, (locking?), check streaming and caching
	case "DELETE":
		handleDelete(w, r)
	case "MKCOL":
		handleMkDir(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		return
	}
}

func handlePropfind(w http.ResponseWriter, r *http.Request) {
	log.Println("Propfind was called")
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))
	/*if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Zugriff verweigert", http.StatusForbidden)
		return
	}*/
	log.Println(filePath)
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		http.Error(w, "Nicht gefunden", http.StatusNotFound)
		return
	}

	depth := r.Header.Get("Depth")
	if depth == "" {
		depth = "1"
	}

	var responses []Response
	if info.IsDir() && (depth == "1" || depth == "infinity") {
		entries, _ := os.ReadDir(filePath)
		for _, entry := range entries {
			entryPath := filepath.Join(r.URL.Path, entry.Name())
			entryInfo, _ := entry.Info()
			responses = append(responses, createResponse(entryPath, entryInfo))
		}
	}
	responses = append(responses, createResponse(r.URL.Path, info))
	//w.Header().Set("Content-Encoding", "identity")
	//w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	//w.WriteHeader(http.StatusMultiStatus)
	/*xmlOut, err := xml.MarshalIndent(Multistatus{
		XMLNS:    "DAV:",
		Response: responses,
	}, "", "  ")
	if err != nil {
		http.Error(w, "Fehler beim Erstellen der XML", http.StatusInternalServerError)
		return
	}
	/* xml.NewEncoder(w).Encode(Multistatus{
		XMLNS:    "DAV:",
		Response: responses,
	})*/
	xmlOutput := fmt.Sprintf(`
	<D:multistatus xmlns:D="DAV:">
	  %s
	</D:multistatus>`, joinResponses(responses))

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write([]byte(xmlOutput))
}

func joinResponses(responses []Response) string {
	var joinedResponses string
	for _, response := range responses {
		joinedResponses += createResponseXML(response)
	}
	return joinedResponses
}

func createResponseXML(response Response) string {
	return fmt.Sprintf(`
    <D:response>
      <D:href>%s</D:href>
      <D:propstat>
        <D:prop>
          <D:displayname>%s</D:displayname>
          <D:getcontentlength>%d</D:getcontentlength>
          <D:getlastmodified>%s</D:getlastmodified>
          <D:resourcetype>%s</D:resourcetype>
        </D:prop>
        <D:status>%s</D:status>
      </D:propstat>
    </D:response>`, response.Href, response.Propstat.Prop.DisplayName, response.Propstat.Prop.ContentLength, response.Propstat.Prop.LastModified, response.Propstat.Prop.ResourceType, response.Propstat.Status)
}

func createResponse(path string, info os.FileInfo) Response {
	return Response{
		Href: path,
		Propstat: Propstat{
			Prop: Prop{
				DisplayName:   info.Name(),
				ResourceType:  GetResourceType(info),
				ContentLength: info.Size(),
				LastModified:  info.ModTime().Format(time.RFC1123),
			},
			Status: "HTTP/1.1 200 OK",
		},
	}
}

func GetResourceType(info os.FileInfo) string {
	if info.IsDir() {
		return "<D:collection/>"
	}
	return ""
}

/*func handleOptions(w http.ResponseWriter, r *http.Request) {
	log.Println("OPTIONS was called")
	w.Header().Set("Allow", "OPTIONS, GET, PUT, DELETE, MKCOL, PROPFIND, MOVE, COPY")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, PUT, DELETE, MKCOL, PROPFIND, MOVE, COPY")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth")

	// DAV-Level angeben (1 = Basis-WebDAV, 2 = erweiterte Features)
	w.Header().Set("DAV", "1, 2")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "WebDAV-OPTIONS: Supported")
}*/

func handleOptions(w http.ResponseWriter, r *http.Request) {
	// Setzt die Header für die CORS-Unterstützung und WebDAV-Methoden
	log.Println("OPTIONS called")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, PUT, DELETE, MKCOL, PROPFIND, MOVE, COPY")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Allow", "OPTIONS, GET, PUT, DELETE, MKCOL, PROPFIND, MOVE, COPY")
	w.Header().Set("Dav", "1, 2") // WebDAV Versionen
	w.Header().Set("Vary", "Origin")
	//w.Header().Set("Date", "Sun, 02 Mar 2025 23:09:12 GMT") // Kann dynamisch gesetzt werden

	// XML Antwortkörper
	/*xmlOutput := `<?xml version="1.0" encoding="UTF-8" ?>
	<D:options-response xmlns:D="DAV:">
		<D:supported-method-set>
			<D:method name="OPTIONS"/>
			<D:method name="GET"/>
			<D:method name="PUT"/>
			<D:method name="DELETE"/>
			<D:method name="MKCOL"/>
			<D:method name="PROPFIND"/>
			<D:method name="MOVE"/>
			<D:method name="COPY"/>
		</D:supported-method-set>
		<D:supported-live-property-set>
			<D:live-property name="resourcetype"/>
			<D:live-property name="displayname"/>
			<D:live-property name="getcontentlength"/>
			<D:live-property name="getlastmodified"/>
		</D:supported-live-property-set>
	</D:options-response>`*/

	// Setzt den Content-Type auf XML
	//w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	// Setzt den Content-Length
	//w.Header().Set("Content-Length", fmt.Sprintf("%d", len(xmlOutput)))
	// Setzt den Statuscode auf 200 OK
	w.WriteHeader(http.StatusOK)
	// Schreibt den XML-Inhalt in die Antwort
	//w.Write([]byte(xmlOutput))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("get was called")
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))

	if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Ungültiger Pfad", http.StatusForbidden)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Datei nicht gefunden", http.StatusNotFound)
		} else {
			http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		http.Error(w, "Verzeichnisse können nicht heruntergeladen werden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(filePath)+"\"")
	http.ServeFile(w, r, filePath)
	log.Printf("File downloaded: %s", filePath)
}
func handleUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("put was called")
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))

	if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Zugriff verweigert", http.StatusForbidden)
		return
	}
	//make dir if needed
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		http.Error(w, "Fehler beim Erstellen des Verzeichnisses", http.StatusInternalServerError)
		return
	}
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Fehler beim Erstellen der Datei", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, r.Body); err != nil {
		http.Error(w, "Fehler beim Schreiben der Datei", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Datei gespeichert: %s\n", filePath)
	log.Printf("Datei hochgeladen: %s", filePath)
}
func handleMkDir(w http.ResponseWriter, r *http.Request) {
	log.Println("mkdir was called")
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))

	if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Zugriff verweigert", http.StatusForbidden)
		return
	}

	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "Verzeichnis existiert bereits", http.StatusMethodNotAllowed)
		return
	}
	//Webdav standard: check if body is empty
	if r.ContentLength > 0 {
		http.Error(w, "MKCOL akzeptiert keinen Request-Body", http.StatusUnsupportedMediaType)
		return
	}

	if err := os.Mkdir(filePath, 0755); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Elternverzeichnis fehlt", http.StatusConflict)
			return
		}
		http.Error(w, "Fehler beim Erstellen des Verzeichnisses", http.StatusInternalServerError)
		return
	}

	// Erfolgsantwort (201 Created)
	w.WriteHeader(http.StatusCreated)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("delete was called")
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))

	if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Zugriff verweigert", http.StatusForbidden)
		return
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "Datei oder Verzeichnis nicht gefunden", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Interner Fehler", http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		err = os.RemoveAll(filePath)
	} else {
		err = os.Remove(filePath)
	}

	if err != nil {
		http.Error(w, "Fehler beim Löschen der Ressource", http.StatusInternalServerError)
		return
	}

	// Erfolgreiche Antwort (204: Kein Inhalt)
	w.WriteHeader(http.StatusNoContent)
}

func (p FileManagerPlugin) Migrate(db *gorm.DB) error {
	log.Println("migrate called")
	return nil
}

// Exportiertes Plugin-Objekt
var Plugin FileManagerPlugin
