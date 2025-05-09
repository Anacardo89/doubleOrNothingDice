package api

import (
	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
	"github.com/Anacardo89/doubleOrNothingDice/internal/email"
)

type AuthHandler struct {
	DB          *db.Manager
	EmailSender *email.EmailSender
}

func NewAuthHandler(dbManager *db.Manager, mailSender *email.EmailSender) *AuthHandler {
	return &AuthHandler{
		DB:          dbManager,
		EmailSender: mailSender,
	}
}
