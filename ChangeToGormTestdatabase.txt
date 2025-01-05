package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type System struct {
	ID   uint   `gorm:"primaryKey"`
	Key  string `gorm:"not null"`
	Data string
}

var db *gorm.DB

func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("datenbank.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatische Migration der Tabellen
	err = db.AutoMigrate(&User{}, &System{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialisierung der System-Tabelle
	initializeSystemTable()
}

func initializeSystemTable() {
	var system System
	err := db.FirstOrCreate(&system, System{Key: "dbsystem", Data: `{"version":1}`}).Error
	if err != nil {
		log.Fatalf("Failed to initialize system table: %v", err)
	}
}

func getCurrentVersion() (int, error) {
	var system System
	if err := db.Where("key = ?", "dbsystem").First(&system).Error; err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(system.Data), &data); err != nil {
		return 0, err
	}

	if version, ok := data["version"].(float64); ok {
		return int(version), nil
	}
	return 0, fmt.Errorf("Version not found in system table")
}

func updateVersion(newVersion int) error {
	var system System
	if err := db.Where("key = ?", "dbsystem").First(&system).Error; err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(system.Data), &data); err != nil {
		return err
	}
	data["version"] = newVersion

	updatedData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	system.Data = string(updatedData)
	return db.Save(&system).Error
}

func migrate() {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		log.Fatalf("Failed to get current version: %v", err)
	}

	if currentVersion < 1 {
		// Beispiel: Migration für Version 1
		if err := migrateToVersion1(); err != nil {
			log.Fatalf("Migration to version 1 failed: %v", err)
		}
		currentVersion = 1
		if err := updateVersion(currentVersion); err != nil {
			log.Fatalf("Failed to update version: %v", err)
		}
	}

	// Weitere Migrationen können hier hinzugefügt werden
	fmt.Printf("Database migrated to version %d\n", currentVersion)
}

func migrateToVersion1() error {
	// Beispielhafte Migration für Version 1
	return db.Create(&System{Key: "dbsystem", Data: `{"version":1}`}).Error
}
