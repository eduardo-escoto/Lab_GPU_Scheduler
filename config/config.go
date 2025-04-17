package config

import (
    "log"
    "os"
)

type Config struct {
    DatabaseDSN string
    SlackToken  string
    SMTPHost    string
    SMTPPort    string
    SMTPUser    string
    SMTPPass    string
}

func LoadConfig() Config {
    cfg := Config{
		DatabaseDSN: getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname"),
        SlackToken:  getEnv("SLACK_TOKEN", ""),
        SMTPHost:    getEnv("SMTP_HOST", "smtp.example.com"),
        SMTPPort:    getEnv("SMTP_PORT", "587"),
        SMTPUser:    getEnv("SMTP_USER", "your-email@example.com"),
        SMTPPass:    getEnv("SMTP_PASS", "your-email-password"),
    }

    log.Println("Configuration loaded successfully")
    return cfg
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}