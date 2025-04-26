package api

import (
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
)

func (h *AuthHandler) RecoverPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}
	claims, err := auth.ParseToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	var req APIMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Type != "recover-password" {
		http.Error(w, "Invalid request type", http.StatusBadRequest)
		return
	}
	var reqPayload RecoverPasswordRequest
	if err := json.Unmarshal(req.Payload, &reqPayload); err != nil {
		http.Error(w, "Error processing request data", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	user, err := h.DB.GetUserByName(ctx, claims.ClientID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	hashedPassword, err := auth.HashPassword(reqPayload.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	err = h.DB.UpdateUserPassword(ctx, user.ID, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	recoverPasswordResponse := RecoverPasswordResponse{
		Message: "Password has been reset successfully",
	}
	writeJSON(w, http.StatusOK, "recover-password", recoverPasswordResponse)
}
