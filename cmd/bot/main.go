package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/WhoYa/subscription-manager/internal/bot"
)

func main() {
	// Получаем токен бота из переменной окружения
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable is required")
	}

	// Получаем базовый URL API (по умолчанию localhost)
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8080"
	}

	// Получаем список админов из переменной окружения
	adminsEnv := os.Getenv("ADMINS")
	if adminsEnv == "" {
		log.Fatal("ADMINS environment variable is required (comma-separated list of Telegram user IDs)")
	}

	// Парсим список админов
	adminStrings := strings.Split(adminsEnv, ",")
	var adminUserIDs []int64
	for _, adminStr := range adminStrings {
		adminStr = strings.TrimSpace(adminStr)
		if adminStr == "" {
			continue
		}
		adminID, err := strconv.ParseInt(adminStr, 10, 64)
		if err != nil {
			log.Fatalf("Invalid admin ID '%s': %v", adminStr, err)
		}
		adminUserIDs = append(adminUserIDs, adminID)
	}

	if len(adminUserIDs) == 0 {
		log.Fatal("At least one admin user ID must be specified in ADMINS environment variable")
	}

	log.Printf("Configured %d admin users", len(adminUserIDs))
	log.Printf("API Base URL: %s", apiBaseURL)

	// Создаем и запускаем бота
	botInstance, err := bot.NewBot(token, apiBaseURL, adminUserIDs)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Println("Starting Telegram bot...")
	if err := botInstance.Start(); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
