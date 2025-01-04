package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
//import "github.com/gookit/goutil/dump" // DEBUG - dump variables with this
var db *sql.DB

// Initialisiere die Datenbank
func initDatabase() {
	var err error
	var result *sql.Rows
	var versionAfterMigration int
	var currentVersion int
	errorsInInit := false
	db, err = sql.Open("sqlite3", "./datenbank.db")
	if err != nil {
		log.Fatal(err)
	}

	// Tabelle erstellen, falls sie nicht existiert
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	/*addAdmin := `
	INSERT INTO users VALUES (
		'1','admin', 'test'
	);`*/
	_, err = db.Exec(createTable)

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	request := `SELECT data FROM system WHERE "key" = 'dbsystem';`
	result, err = db.Query(request)
	if (err != nil) {
		currentVersion = 0
		versionAfterMigration = migrate(currentVersion)
		if (versionAfterMigration >= 1) {
			dbupdate := updateDBVersion(versionAfterMigration)
			errorsInInit = errorsInInit || dbupdate
		} else {
			errorsInInit = true
		}
	} else
	{
		defer result.Close()
		var jsonData string
		result.Next()
		err = result.Scan(&jsonData)
		if (err != nil) {
			log.Fatalf("Error retrieving the version from the database: %v", err)
		}
		currentVersion, err = extractVersionFromJSON(jsonData)
		if (err != nil) {
			log.Fatalf("Unable to extract version from JSON! The database might be corrupt. %v", err)
		} else {
			versionAfterMigration = migrate(currentVersion)
			if (versionAfterMigration > currentVersion) {
				dbupdate := updateDBVersion(versionAfterMigration)
				errorsInInit = errorsInInit || dbupdate
			}
		}
	}
	if (versionAfterMigration > currentVersion) {
		fmt.Printf("Database was migrated from version %v to %v\n", currentVersion, versionAfterMigration)
	}
	if (!errorsInInit) {
		fmt.Println("Database initialized successfully.")
	} else {
		fmt.Println("Database was not initialized successfully!\n Please consult the logs for more information.")
	}
}

func migrate(version int) int {
	if (version < 1) {
		request:= `
		-- Create the table "system"
		CREATE TABLE system (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			"key" NVARCHAR(128) NOT NULL,
			data NVARCHAR
		);

		-- Insert the first row
		INSERT INTO system ("key", data)
		VALUES ('dbsystem', '{"version":1}');
		`
		if (tryExecute(request, true)) {
			version = 1
		}
	}
	if (version <= 1) {
		// etc. Future db migrations go here. Always execute db migrations from here; never by hand!
		// version = 2 // make sure to use tryExecute
	}
	fmt.Println("System database version:", version)
	return version
}

func tryExecute(request string, logfatal ...bool) bool { // Todo: Maybe move this to some kind of db utils file
	var err error
	_, err = db.Exec(request)
	if (err != nil) {
		if (logfatal != nil) {
			log.Fatalf("Failed to fulfill request: %v due to error: %v", request, err) // Todo: circle back on whether this formatting is correct
		}
		return false
	}
	return true
}

func extractVersionFromJSON(jsonStr string) (int, error) {
	// Define a map to hold the parsed JSON
	var data map[string]interface{}

	// Parse the JSON string into the map
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}

	// Extract the "version" field and assert it's a float64 (default JSON number type in Go)
	if version, ok := data["version"].(float64); ok {
		return int(version), nil // Convert float64 to int and return it
	}

	return 0, fmt.Errorf("key 'version' not found or not a number")
}

func updateAttributeInJSON(jsonStr string, attributeName string, value interface{}) (interface{}, error) { // Todo: Move this functions to some kind of util
	// Define a map to hold the parsed JSON
	var data map[string]interface{}

	// Parse the JSON string into the map
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err // Return an empty string and the error if parsing fails
	}

	// Set the "version" field to 2
	data[attributeName] = value

	// Convert the updated map back to a JSON string
	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return "", err // Return an empty string and the error if marshaling fails
	}

	return string(updatedJSON), nil
}


func updateDBVersion(version int) bool {
	var result *sql.Rows
	var err error
	request := `SELECT data FROM system WHERE "key" = 'dbsystem';`
	result, err = db.Query(request)
	if (err != nil) {
		log.Fatalf("Error retrieving the version from the database")
	} else {
		var data string
		result.Scan(&data)
		updatedData, err := updateAttributeInJSON(data, "version", version)
		if (err != nil) {
			fmt.Errorf("Error updating the version in the database json: %w", err)
			return false
		} else {
			request = `UPDATE system SET data = ? WHERE "key" = 'dbsystem'`
			_, err := db.Exec(request, updatedData)
			if err != nil {
				fmt.Errorf("Error updating version in database: %w", err)
				return false
			}
			return true	
		}
	}
	return false
}