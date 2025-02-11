package main

import (
	//"AllInOnPrivateCloud/plugininterface"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	initDatabase()
	// Serve static files from the "static" directory
	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)
	mux := http.NewServeMux()

	handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowOriginFunc:    allowOrigins,
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"}, // TODO: Move this to config file
		OptionsPassthrough: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true, // TODO: add environment config / debug flag or whatever to automatically distinguish between debug environment and production environment
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
	mux.HandleFunc("/check-token", checkToken)

	// API-Endpunkt für die Plugin-Liste

	mux.HandleFunc("/api/plugins", func(w http.ResponseWriter, r *http.Request) {
		pl := loadPlugins(mux)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pl)
	})

	// Start the server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler) // TODO: Port auslagern in Config
}

func allowOrigins(origin string) bool {
	// TODO: Add proper handling here & configuration & handle the configuration in a setup script.
	return true
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Registrierung soll von Admin genehmigt werden, Anzahl unbestätigter Registrierungen beschränkt
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

	user := User{
		Username: creds.Username,
		Password: string(hashedPassword),
	}
	result := db.Create(&user)

	// Fehlerbehandlung
	if result.Error != nil {
		http.Error(w, "Failed to create user: "+result.Error.Error(), http.StatusInternalServerError)
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

	var user User
	result := db.Where("username = ?", creds.Username).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Vergleiche das Passwort nur, wenn der Benutzer existiert
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	//expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		Username: creds.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	generateToken(w, claims)

}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: funktionalität gleicht checkToken(), Handler eventuell nicht mehr benötigt
	// Überprüfen des Cookies
	log.Println("welcomeHandler")
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Kein Cookie vorhanden
			http.Error(w, "Token fehlt", http.StatusUnauthorized)
			return
		}
		// Fehler beim Abrufen des Cookies
		http.Error(w, "Fehler beim Abrufen des Tokens", http.StatusBadRequest)
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
		// Token ist ungültig
		log.Println("Token ungültig:", err)
		http.Error(w, "Ungültiges Token", http.StatusUnauthorized)
		return
	}

	log.Printf("Token gültig für Benutzer: %s (ID: %d)\n", claims.Username, claims.UserID)
	// Erfolgreich authentifiziert, User willkommen heißen
	response := map[string]string{
		"message":  fmt.Sprintf("Welcome %s!", claims.Username),
		"username": claims.Username,
		"id":       claims.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//TODO: funktion in LoadPlugin aus pluginmanager integrieren
/*func loadPlugins() {
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
			if err != nil {
				log.Printf("Fehler beim Suchen nach Symbol 'Plugin' in %s: %v", path, err)
				return nil
			}

			// Typprüfung und Registrierung
			if plg, ok := sym.(plugininterface.Plugin); ok {
				//plg.Register(mux)
				metadata := plg.Metadata()

				log.Printf("Plugin registriert: %s", metadata.Name)

				// Füge Plugin-Metadaten zur globalen Liste hinzu
				pluginList = append(pluginList, metadata)

				api := plugininterface.Api{
					Metadata: metadata,
					Mux:      http.NewServeMux(),
					RegisterWidget: func() error {
						fmt.Println("RegisterWidget called") // TODO: implement RegisterWidget functionality here or preferrably in a API scripts file
						return nil
					},
					RegisterMenuEntry: func() error {
						fmt.Println("RegisterMenuEntry called") // TODO: implement RegisterMenu functionality here or preferrably in a API scripts file
						return nil
					},
				}

				err = plg.Init(api)
				if err != nil {
					log.Printf("Fehler beim Laden des Plugins %s: %v", path, err)
				}
			} else {
				log.Printf("Ungültiger Plugin-Typ in %s", path) // TODO: maybe remove this
				return fmt.Errorf("ungültiger plugin-typ in %s", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Fehler beim Laden der Plugins: %v", err)
	} else if pluginList == nil {
		log.Println("Keine Plugins gefunden")
	}
}*/
