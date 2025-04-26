package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	JWT struct {
		Secret        string `yaml:"secret"`
		ExpiryMinutes int    `yaml:"expiry_minutes"`
	} `yaml:"jwt"`

	Email struct {
		SMTPHost       string `yaml:"smtp_host"`
		SMTPPort       int    `yaml:"smtp_port"`
		SenderEmail    string `yaml:"sender_email"`
		SenderPassword string `yaml:"sender_password"`
	} `yaml:"email"`
}

var AppConfig *Config

func LoadConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("failed to decode config: %v", err)
	}

	AppConfig = &cfg
}
