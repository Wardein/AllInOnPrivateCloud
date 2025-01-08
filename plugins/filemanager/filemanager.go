package main

import (
	"log"
	"main/plugininterface"

	"gorm.io/gorm"
)

//go build -buildmode=plugin -o plugins/filemanager/filemanager.so plugins/filemanager/filemanager.go

// FileManagerPlugin implementiert das Plugin-Interface.
type FileManagerPlugin struct{}

/*func (p FileManagerPlugin) Register(mux *http.ServeMux) {
	mux.HandleFunc("/plugins/filemanager", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("plugins", "filemanager", "index.html"))
	})
}*/

func (p FileManagerPlugin) Metadata() plugininterface.PluginMetadata {
	return plugininterface.PluginMetadata{
		Name:          "File Manager",
		Description:   "Ein Plugin zum Verwalten von Dateien",
		Path:          "/plugins/filemanager.html",
		MenuButton:    true,
		UsingDatabase: true,
	}
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
