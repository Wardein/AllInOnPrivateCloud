package main

import (
	"main/plugininterface"
)

//This Plugin adds a URL as a Button in the Menu.
//For a different URL, this file and the HTML file must be copied and saved under a new name, with the name and path adjusted accordingly

//go build -buildmode=plugin -o plugins/anotherURL/anotherURL.so plugins/anotherURL/anotherURL.go

type URLPlugin struct{}

func (p URLPlugin) Metadata() plugininterface.PluginMetadata {

	return plugininterface.PluginMetadata{
		Name:          "example.com",
		Description:   "Only another URL",
		Path:          "/plugins/anotherURL.html",
		MenuButton:    true,
		UsingDatabase: false,
	}
}

func (p URLPlugin) Routes() []plugininterface.Route {
	return nil
}

// Exportiertes Plugin-Objekt
var Plugin URLPlugin
