package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	botpkg "hadith-bot/internal/bot"
	"hadith-bot/internal/config"
	"hadith-bot/internal/logger"
	"hadith-bot/internal/services"

	telebot "github.com/tucnak/telebot"
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

	// Create Telegram bot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 60},
	})
	if err != nil {
		log.Fatal("Failed to create bot: %v", err)
	}

	log.Info("Bot created successfully")

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

	// Create handler
	handler := botpkg.NewHandler(
		bot,
		hadithService,
		log,
		cfg.RateLimitRequests,
		cfg.RateLimitWindow,
	)
	handler.HandleCommands()

	log.Info("Bot is ready to handle commands")

	// Handle graceful shutdown
	// Start cleanup goroutine for rate limiter
	go func() {
		ticker := make(chan os.Signal, 1)
		signal.Notify(ticker, syscall.SIGINT, syscall.SIGTERM)
		<-ticker
		log.Info("Shutting down bot...")
		os.Exit(0)
	}()

	// Start bot (this blocks)
	bot.Start()
}
