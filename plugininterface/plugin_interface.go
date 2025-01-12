package plugininterface

import (
	"net/http"

	"gorm.io/gorm"
)

// PluginMetadata enthält Metadaten für ein Plugin.
type PluginMetadata struct {
	Name          string `json:"Name"`
	Description   string `json:"Description"`
	Path          string `json:"Path"`
	MenuButton    bool   `json:"MenuButton"`    //TODO: implement the Option to not use the menu
	UsingDatabase bool   `json:"UsingDatabase"` //TODO: Implement automaigrate functions for plugins
}

// Plugin definiert die Schnittstelle, die jedes Plugin implementieren muss.
type Plugin interface {
	Metadata() PluginMetadata
	//Migrate(db *gorm.DB) error
	//Initialize() error
}

type PluginDatabase interface {
	Migrate(db *gorm.DB) error
	Init(Api) error
}

type Api struct {
	Metadata          PluginMetadata
	Mux               *http.ServeMux // shared resource
	RegisterWidget    func() error
	RegisterMenuEntry func() error
}
