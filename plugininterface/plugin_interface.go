package plugininterface

import "gorm.io/gorm"

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
}
