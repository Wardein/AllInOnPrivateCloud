package main

import (
	"fmt"
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
var pluginsLoaded bool

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

			addPlugin(path, &pl)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Fehler beim Laden der Plugins: %v", err)
	} else if pl == nil {
		log.Println("Keine Plugins gefunden")
	}
	pluginsLoaded = true
	return pl
}

func addPlugin(path string, pl *[]plugininterface.PluginMetadata) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("fehler beim Laden des Plugins %s: %v", path, err)
	}
	sym, err := p.Lookup("Plugin") // Symbol "Plugin" laden
	if err != nil {
		return fmt.Errorf("fehler beim Suchen nach Symbol 'Plugin' in %s: %v", path, err)
	}
	// Typprüfung und Registrierung
	plg, ok := sym.(plugininterface.Plugin)
	if ok {
		//plg.Register(mux)
		log.Printf("Plugin registriert: %s", plg.Metadata().Name)

		// Füge Plugin-Metadaten zur globalen Liste hinzu
		*pl = append(*pl, plg.Metadata())
	} else {
		return fmt.Errorf("ungültiger Plugin-Typ in %s", path)
	}
	if dbPlugin, ok := plg.(plugininterface.PluginDatabase); ok {
		log.Printf("Migration wird für Plugin '%s' gestartet", plg.Metadata().Name)

		// Migration aufrufen
		if err := dbPlugin.Migrate(db); err != nil {
			return fmt.Errorf("fehler bei der Migration für Plugin '%s': %v", plg.Metadata().Name, err)
		} else {
			log.Printf("Migration für Plugin '%s' erfolgreich abgeschlossen", plg.Metadata().Name)
		}
	}
	return nil
}
