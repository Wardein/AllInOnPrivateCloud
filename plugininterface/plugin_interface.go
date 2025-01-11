package plugininterface

import (
	"net/http"
)

// PluginMetadata enthält Metadaten für ein Plugin.
type PluginMetadata struct {
	Name        string
	Description string
	Path        string
}

// Plugin definiert die Schnittstelle, die jedes Plugin implementieren muss.
type Plugin interface {
	Metadata() PluginMetadata
	Init(Api) error
}

type Api struct {
	Metadata          PluginMetadata
	Mux               *http.ServeMux // shared resource
	RegisterWidget    func() error
	RegisterMenuEntry func() error
}
