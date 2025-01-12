#!/bin/bash

# Pfad zum Plugins-Ordner
PLUGIN_DIR="plugins"

# Für jeden Unterordner im Plugins-Verzeichnis
for folder in "$PLUGIN_DIR"/*; do
  if [ -d "$folder" ]; then # Prüfen, ob es ein Ordner ist
    folder_name=$(basename "$folder")
    go_file="$folder/$folder_name.go"
    output_file="$folder/$folder_name.so"
    
    if [ -f "$go_file" ]; then # Prüfen, ob die .go-Datei existiert
      echo "Building plugin for $folder_name"
      go build -buildmode=plugin -o "$output_file" "$go_file"
    else
      echo "Skipping $folder_name: $go_file not found"
    fi
  fi
done