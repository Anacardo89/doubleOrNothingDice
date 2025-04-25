package http

import (
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)
	if creds.Username != "alice" || creds.Password != "password123" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
