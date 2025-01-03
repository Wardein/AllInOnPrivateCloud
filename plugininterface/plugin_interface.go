package plugininterface

// PluginMetadata enthält Metadaten für ein Plugin.
type PluginMetadata struct {
	Name        string
	Description string
	Path        string
}

// Plugin definiert die Schnittstelle, die jedes Plugin implementieren muss.
type Plugin interface {
	Metadata() PluginMetadata
}
