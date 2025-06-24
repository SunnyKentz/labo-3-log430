package auth

import (
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"
	"net/http"
)

const pw = "password" // TODO: Replace with env/config in production

// CheckPassword checks if the provided password matches the secret password.
// Returns true if the password matches, false otherwise.
func IsUserPasswordValid(nom string, password string) bool {
	println(nom)
	return password == pw
}

func IsUsernameValid(nom string) bool {
	resp, err := http.PostForm(API_MERE()+"/api/v1/login", map[string][]string{
		"username": {nom},
		"role":     {"commis"},
	})
	if err != nil {
		log.Println("Erreur lors de la connexion: " + err.Error())
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
