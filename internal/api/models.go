package api

import (
	"encoding/json"
)

type APIMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type ActivateResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordRsponse struct {
	Message string `json:"message"`
}

type RecoverPasswordRequest struct {
	Password string `json:"password"`
}

type RecoverPasswordResponse struct {
	Message string `json:"message"`
}
