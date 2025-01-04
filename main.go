package main

import (
	//"AllInOnPrivateCloud/plugininterface"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/plugininterface"
	"net/http"
	"github.com/rs/cors"
	"os"
	"path/filepath"
	"plugin"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var pluginList []plugininterface.PluginMetadata

func main() {
	initDatabase()
	// Serve static files from the "static" directory
	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)
	mux := http.NewServeMux()

	handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowOriginFunc: allowOrigins,
		AllowCredentials: true,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Todo: Move this to config file
		OptionsPassthrough: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true, // Todo: add environment config / debug flag or whatever to automatically distinguish between debug environment and production environment
	})
	handler = c.Handler(handler)

	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./frontend/styles"))))
	mux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./frontend/scripts"))))
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./frontend/html"))))

	//http.ListenAndServe(":8080", mux)
	// API routes
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/welcome", welcomeHandler)

	loadPlugins()
	log.Println(pluginList)

	// API-Endpunkt für die Plugin-Liste
	mux.HandleFunc("/api/plugins", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pluginList)
	})

	// Start the server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler) // Todo: Port auslagern in Config
}

func allowOrigins(origin string) bool {
	// Todo: Add proper handling here & configuration & handle the configuration in a setup script.
	return true
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", creds.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", creds.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,                    // Verhindert JavaScript-Zugriff auf Cookies
		Secure:   true,                    // Nur über HTTPS verfügbar
		SameSite: http.SameSiteStrictMode, // Verhindert Cross-Site Cookie-Zugriffe
	})
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Überprüfen des Cookies
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Kein Cookie gefunden, also unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token fehlt"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fehler beim Abrufen des Tokens"))
		return
	}

	// Token aus dem Cookie extrahieren
	tokenStr := cookie.Value
	claims := &Claims{}

	// Token parsen und Claims extrahieren
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		// Ungültiges Token, Unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Ungültiges Token"))
		return
	}

	// Erfolgreich authentifiziert, User willkommen heißen
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome %s!", claims.Username)})
}

func loadPlugins() {
	pluginDir := "./plugins"

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
			log.Println(sym)
			if err != nil {
				log.Printf("Fehler beim Suchen nach Symbol 'Plugin' in %s: %v", path, err)
				return nil
			}

			// Typprüfung und Registrierung
			if plg, ok := sym.(plugininterface.Plugin); ok {
				//plg.Register(mux)
				log.Printf("Plugin registriert: %s", plg.Metadata().Name)

				// Füge Plugin-Metadaten zur globalen Liste hinzu
				pluginList = append(pluginList, plg.Metadata())
			} else {
				log.Printf("Ungültiger Plugin-Typ in %s", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Fehler beim Laden der Plugins: %v", err)
	}
}
