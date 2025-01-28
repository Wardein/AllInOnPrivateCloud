package main

import (
	"log"
	"main/plugininterface"
)

//go build -buildmode=plugin -o plugins/notes/notes.so plugins/notes/notes.go

// settingsPlugin implementiert das Plugin-Interface.
type settingsPlugin struct {
	//PluginMetadata
}

/*func (p NotesPlugin) Register(mux *http.ServeMux) {
	mux.HandleFunc("/plugins/notes", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("plugins", "notes", "index.html"))
	})
}*/

func (p settingsPlugin) Metadata() plugininterface.PluginMetadata {
	log.Println("DEBUG@settings") //Plugin.mux.HandleFunc("/welcome", welcomeHandler)
	log.Println(Plugin)
	return plugininterface.PluginMetadata{
		Name:        "Settings",
		Description: "Ein Plugin zum Verwalten der Einstellungen",
		Path:        "/plugins/settings.html",
	}
}

func (p settingsPlugin) Routes() []plugininterface.Route {
	return nil
}

func (p settingsPlugin) Init(api plugininterface.Api) error {
	log.Println("DEBUG@settings - Init")
	log.Println(api)
	return nil
}

// Exportiertes Plugin-Objekt
var Plugin settingsPlugin
