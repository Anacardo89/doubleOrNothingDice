package api

import (
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
)

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req APIMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Type != "login" {
		http.Error(w, "Invalid request type", http.StatusBadRequest)
		return
	}
	var reqPayload LoginRequest
	if err := json.Unmarshal(req.Payload, &reqPayload); err != nil {
		http.Error(w, "Error processing request data", http.StatusInternalServerError)
		return
	}
	if reqPayload.Username == "" || reqPayload.Password == "" {
		http.Error(w, "username/email and password are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user, err := h.DB.GetUserByName(ctx, reqPayload.Username)
	if err != nil {
		http.Error(w, "error running query", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if !user.IsActive {
		http.Error(w, "user not active", http.StatusUnauthorized)
		return
	}
	if err := auth.CheckPasswordHash(reqPayload.Password, user.PasswordHash); err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	loginResponse := LoginResponse{
		UserID: user.ID,
		Token:  token,
	}
	writeJSON(w, http.StatusOK, "login", loginResponse)
}
