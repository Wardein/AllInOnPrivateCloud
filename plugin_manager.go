package main

import (
	"log"
	"main/plugininterface"
	"os"
	"path/filepath"
	"plugin"
)

type PluginManager struct {
	plugins []plugininterface.Plugin
}

// RegisterPlugin registriert ein Plugin
func (pm *PluginManager) RegisterPlugin(plugin plugininterface.Plugin) {
	pm.plugins = append(pm.plugins, plugin)
}

// Migrate führt die Migrationen aller registrierten Plugins aus
/*func (pm *PluginManager) Migrate(db *gorm.DB) error {
	for _, plugin := range pm.plugins {
		fmt.Printf("Migration für Plugin %s starten...\n", plugin.Metadata().Name)
		if err := plugin.Migrate(db); err != nil {
			return fmt.Errorf("Migration für Plugin %s fehlgeschlagen: %w", plugin.Metadata().Name, err)
		}
	}
	return nil
}*/

// Initialize initialisiert alle registrierten Plugins
/*func (pm *PluginManager) Initialize() error {
	for _, plugin := range pm.plugins {
		fmt.Printf("Initialisierung von Plugin %s starten...\n", plugin.Metadata().Name)
		if err := plugin.Initialize(); err != nil {
			return fmt.Errorf("Initialisierung von Plugin %s fehlgeschlagen: %w", plugin.Metadata().Name, err)
		}
	}
	return nil
}*/

func loadPlugins() []plugininterface.PluginMetadata {
	pluginDir := "./plugins"
	var pl []plugininterface.PluginMetadata
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Nur .so-Dateien laden
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			log.Printf("Lade Plugin: %s", path)

			// Plugin öffnen
			p, err := plugin.Open(path)
			if err != nil {
				log.Printf("Fehler beim Laden des Plugins %s: %v", path, err)
				return nil
			}

			// Symbol "Plugin" laden
			sym, err := p.Lookup("Plugin")
			if err != nil {
				log.Printf("Fehler beim Suchen nach Symbol 'Plugin' in %s: %v", path, err)
				return nil
			}

			// Typprüfung und Registrierung
			if plg, ok := sym.(plugininterface.Plugin); ok {
				//plg.Register(mux)
				log.Printf("Plugin registriert: %s", plg.Metadata().Name)

				// Füge Plugin-Metadaten zur globalen Liste hinzu
				pl = append(pl, plg.Metadata())
			} else {
				log.Printf("Ungültiger Plugin-Typ in %s", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Fehler beim Laden der Plugins: %v", err)
	} else if pl == nil {
		log.Println("Keine Plugins gefunden")
	}
	return pl
}
