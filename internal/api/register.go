package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anacardo89/doubleOrNothingDice/config"
	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
)

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req APIMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Type != "register" {
		http.Error(w, "Invalid request type", http.StatusBadRequest)
		return
	}
	var reqPayload RegisterRequest
	if err := json.Unmarshal(req.Payload, &reqPayload); err != nil {
		http.Error(w, "Error processing request data", http.StatusInternalServerError)
		return
	}
	if reqPayload.Username == "" ||
		reqPayload.Email == "" ||
		reqPayload.Password == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	exists, err := h.DB.IsUsernameTaken(ctx, reqPayload.Username)
	if err != nil {
		http.Error(w, "error running query", http.StatusBadRequest)
		return
	}
	if exists {
		http.Error(w, "username already exists", http.StatusBadRequest)
		return
	}
	exists, err = h.DB.IsEmailTaken(ctx, reqPayload.Email)
	if err != nil {
		http.Error(w, "error running query", http.StatusBadRequest)
		return
	}
	if exists {
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(reqPayload.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	user := db.User{
		Username:     reqPayload.Username,
		Email:        reqPayload.Email,
		PasswordHash: hashedPassword,
	}
	err = h.DB.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to create user in the database", http.StatusInternalServerError)
		return
	}
	activationToken, err := auth.GenerateToken(reqPayload.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	activationLink := fmt.Sprintf("http://%s:%s/activate?token=%s", config.AppConfig.Server.Host, strconv.Itoa(config.AppConfig.Server.Port), activationToken)
	subject := "Double or Nothing Dice - Activate your account"
	body := fmt.Sprintf("Click the following link to verify your account: %s", activationLink)
	if err := h.EmailSender.Send(reqPayload.Email, subject, body); err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}
	registerResponse := RegisterResponse{
		UserID:  user.ID,
		Message: "User created, please verify your email",
	}
	writeJSON(w, http.StatusCreated, "register", registerResponse)
}
