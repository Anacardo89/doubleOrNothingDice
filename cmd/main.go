package main

import (
	"log"
	"strconv"

	"github.com/Anacardo89/doubleOrNothingDice/config"
	"github.com/Anacardo89/doubleOrNothingDice/internal/api"
	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
	"github.com/Anacardo89/doubleOrNothingDice/internal/email"
	"github.com/Anacardo89/doubleOrNothingDice/internal/redis"
	"github.com/Anacardo89/doubleOrNothingDice/internal/server"
)

func main() {
	config.LoadConfig("config/config.yaml")
	db.Connect(config.AppConfig.Database.User, config.AppConfig.Database.Password, config.AppConfig.Database.Host, config.AppConfig.Database.Port, config.AppConfig.Database.DBName)
	dbManager := db.NewManager(db.DB)
	emailSender := email.NewEmailSender(config.AppConfig.Email.SMTPHost, config.AppConfig.Email.SMTPPort, config.AppConfig.Email.SenderEmail, config.AppConfig.Email.SenderPassword)
	authHandler := api.NewAuthHandler(dbManager, emailSender)
	redisClient := redis.NewClient(config.AppConfig.Redis.Host, config.AppConfig.Redis.Port, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)
	port := config.AppConfig.Server.Port
	server := server.NewServer(authHandler, redis.NewManager(redisClient))
	log.Println("Starting the server...")
	server.Run(":" + strconv.Itoa(port))
}
