package api

import (
	"net/http"

	"github.com/Anacardo89/doubleOrNothingDice/internal/auth"
)

func (h *AuthHandler) ActivateHandler(w http.ResponseWriter, r *http.Request) {
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
	ctx := r.Context()
	user, err := h.DB.GetUserByName(ctx, claims.ClientID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user.IsActive = true
	err = h.DB.ActivateUser(ctx, user.ID)
	if err != nil {
		http.Error(w, "Failed to activate user", http.StatusInternalServerError)
		return
	}
	activateResponse := ActivateResponse{
		UserID:  user.ID,
		Message: "Account activated successfully",
	}
	writeJSON(w, http.StatusOK, "activate", activateResponse)
}
