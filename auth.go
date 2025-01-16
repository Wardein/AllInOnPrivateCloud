package main

import (
	//"AllInOnPrivateCloud/plugininterface"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"userid"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

func generateToken(w http.ResponseWriter, claims *Claims) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Fehler beim Signieren des Tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  claims.ExpiresAt.Time,
		HttpOnly: true,                    // Verhindert JavaScript-Zugriff auf Cookies
		Secure:   false,                   // Nur über HTTPS verfügbar
		SameSite: http.SameSiteStrictMode, // Verhindert Cross-Site Cookie-Zugriffe
	})
	return nil
}

func checkToken(w http.ResponseWriter, r *http.Request) { // error)
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
	tokenStr := cookie.Value
	claims := &Claims{}

	// Token parsen und Claims extrahieren
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unerwartete Signaturmethode: %v", token.Header["alg"])
		}
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
	response := map[string]interface{}{
		"message":  fmt.Sprintf("Welcome %s!", claims.Username),
		"username": claims.Username,
		"userid":   claims.UserID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
