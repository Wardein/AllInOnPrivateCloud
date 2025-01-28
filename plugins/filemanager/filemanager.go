package main

import (
	"encoding/json"
	"log"
	"main/plugininterface"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

//go build -buildmode=plugin -o plugins/filemanager/filemanager.so plugins/filemanager/filemanager.go

// FileManagerPlugin implementiert das Plugin-Interface.
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
		{Path: "/files", Handler: listFilesHandler},
	}
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listFiles(basePath)
	if err != nil {
		log.Panicln("error read files")
		http.Error(w, "Failed to read files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
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
