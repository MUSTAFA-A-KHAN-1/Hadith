package main

import (
	"context"
	"fmt"

	// "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"hadith-bot/internal/bot"
	"hadith-bot/internal/config"
	"hadith-bot/internal/logger"
	"hadith-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load .env file if exists
	godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	logLevel := logger.InfoLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = logger.DebugLevel
	case "warn":
		logLevel = logger.WarnLevel
	case "error":
		logLevel = logger.ErrorLevel
	}
	log := logger.New(os.Stdout, logLevel, true)
	log.WithPrefix("hadith-bot")

	log.Info("Starting Hadith Portal Bot...")

	// Create Telegram bot using go-telegram-bot-api/v5
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal("Failed to create bot: %v", err)
	}

	// Set debug mode based on log level
	botAPI.Debug = cfg.LogLevel == "debug"

	log.Info("Bot created successfully (using go-telegram-bot-api/v5)")

	// Create hadith service
	// Try to load from data directory - check multiple possible locations
	dataDir := "../../src/data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = "../src/data"
	}
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = "./data"
	}
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		dataDir = ""
	}

	hadithService := services.NewHadithService(dataDir, cfg.APIURL, cfg.APIKey, cfg.APITimeout, log)
	log.Info("Hadith service initialized")

	// Create handler (no global mutable state - botAPI is passed as dependency)
	handler := bot.NewHandler(botAPI, hadithService, log, cfg.RateLimitRequests, cfg.RateLimitWindow)

	log.Info("Bot is ready to handle commands")

	// Set up update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60 // Timeout for long polling

	// Get update channel
	updates := botAPI.GetUpdatesChan(updateConfig)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle updates in a goroutine
	go func() {
		handler.HandleUpdates(ctx, updates)
	}()

	log.Info("Bot is now running...")

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down bot...")

	// Clear the updates channel before exiting
	for range updates {
		// Drain the channel
	}

	log.Info("Bot shutdown complete")
}

// sendMessageWithTimeout sends a message with context timeout
func sendMessageWithTimeout(botAPI *tgbotapi.BotAPI, chatID int64, text string, parseMode string, replyMarkup interface{}) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := tgbotapi.NewMessage(chatID, text)
	if parseMode != "" {
		msg.ParseMode = parseMode
	}
	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	_, err := botAPI.Send(msg)

	return err
}
