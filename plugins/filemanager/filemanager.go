package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/plugininterface"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

const basePath = "./files" // TODO ./files/{user} and Auth

type FileInfo struct {
	Name string `json:"name"`
	//Type string `json:"type"`
	IsDir    bool   `json:"isDir"`
	FullPath string `json:"fullPath"`
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
		{Path: "/", Handler: fileHandler},
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PROPFIND":
		handleListFiles(w, r) //TODO: Ausgabeformat beachten!
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
	}
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := listFiles(basePath) //TOdo: handle empty file folder
	if err != nil {
		log.Println("error read files" + err.Error())
		http.Error(w, "Failed to read files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "OPTIONS, GET, PUT, DELETE, MKCOL, PROPFIND, MOVE, COPY")

	// DAV-Level angeben (1 = Basis-WebDAV, 2 = erweiterte Features)
	w.Header().Set("DAV", "1, 2")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "WebDAV-OPTIONS: Supported")
}
func handleDownload(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(basePath, filepath.Clean(r.URL.Path))

	if !strings.HasPrefix(filePath, basePath) {
		http.Error(w, "Ungültiger Pfad", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, filePath)
	log.Println("Served file: " + filePath)
}
func handleUpload(w http.ResponseWriter, r *http.Request) {
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
}
func handleMkDir(w http.ResponseWriter, r *http.Request) {
	// mkdir
}
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// delete file
}

func listFiles(path string) ([]FileInfo, error) {
	log.Println("listFiles called")
	var fileList []FileInfo
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, _ := filepath.Rel(basePath, filePath)
		fileList = append(fileList, FileInfo{
			Name:     info.Name(),
			IsDir:    info.IsDir(),
			FullPath: relativePath,
		})
		return nil
	})

	return fileList, err
}

func (p FileManagerPlugin) Migrate(db *gorm.DB) error {
	log.Println("migrate called")
	return nil
}

/*func Initialize() error {
	return nil
}*/

// Exportiertes Plugin-Objekt
var Plugin FileManagerPlugin
