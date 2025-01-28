package main

import (
	"main/plugininterface"
)

//go build -buildmode=plugin -o plugins/notes/notes.so plugins/notes/notes.go

// NotesPlugin implementiert das Plugin-Interface.
type NotesPlugin struct{}

/*func (p NotesPlugin) Register(mux *http.ServeMux) {
	mux.HandleFunc("/plugins/notes", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("plugins", "notes", "index.html"))
	})
}*/

func (p NotesPlugin) Metadata() plugininterface.PluginMetadata {
	return plugininterface.PluginMetadata{
		Name:          "Notes",
		Description:   "Ein Plugin zum Erstellen und Verwalten von Notizen",
		Path:          "/plugins/notes.html",
		MenuButton:    true,
		UsingDatabase: true,
	}
}

func (p NotesPlugin) Routes() []plugininterface.Route {
	return nil
}

// Exportiertes Plugin-Objekt
var Plugin NotesPlugin
