package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/config"
	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
)

func (h *AuthHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req APIMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Type != "forgot-password" {
		http.Error(w, "Invalid request type", http.StatusBadRequest)
		return
	}
	var reqPayload ForgotPasswordRequest
	if err := json.Unmarshal(req.Payload, &reqPayload); err != nil {
		http.Error(w, "Error processing request data", http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	user, err := h.DB.GetUserByEmail(ctx, reqPayload.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	resetToken, err := auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}
	resetLink := fmt.Sprintf("http://%s:%d/recover-password?token=%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port, resetToken)
	subject := "Double or Nothing Dice - Password Reset Request"
	body := fmt.Sprintf("Click the link to reset your password:\n%s\nSend password in JSON body with the above link.", resetLink)
	if err := h.EmailSender.Send(user.Email, subject, body); err != nil {
		http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
		return
	}
	forgotPasswordRsponse := ForgotPasswordRsponse{
		Message: "Password reset email sent",
	}
	writeJSON(w, http.StatusOK, "forgot-password", forgotPasswordRsponse)
}
