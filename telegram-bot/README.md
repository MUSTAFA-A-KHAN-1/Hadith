# Hadith Portal Telegram Bot

A Telegram bot for browsing and searching hadith collections, built with Go using clean architecture.

## Features

- **Browse Collections**: Explore the six major hadith collections
- **Search Hadiths**: Search hadiths by keyword with pagination
- **Random Hadith**: Get a random hadith for daily inspiration
- **Inline Keyboards**: User-friendly navigation with inline buttons
- **Pagination**: Browse through books and hadiths with next/previous buttons
- **MarkdownV2**: Properly formatted messages with Markdown support
- **Rate Limiting**: Basic spam protection
- **Graceful Error Handling**: Never crashes, logs errors properly

## Supported Collections

- Sahih al-Bukhari
- Sahih Muslim
- Sunan Abu Dawood
- Jami' at-Tirmidhi
- Sunan an-Nasa'i
- Sunan Ibn Majah

## Commands

| Command | Description |
|---------|-------------|
| `/start` | Welcome message and main menu |
| `/help` | Show help information |
| `/collections` | Browse hadith collections |
| `/search <keyword>` | Search hadiths |
| `/random` | Get a random hadith |

## Project Structure

```
telegram-bot/
├── cmd/
│   └── bot/
│       └── main.go           # Application entry point
├── internal/
│   ├── bot/
│   │   ├── handlers.go       # Command and callback handlers
│   │   └── ratelimiter.go    # Rate limiting implementation
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── data/
│   │   └── loader.go         # Hadith data loader
│   ├── logger/
│   │   └── logger.go         # Structured logging
│   ├── models/
│   │   └── models.go         # Data models
│   └── services/
│       └── hadith.go         # Hadith business logic
├── data/                      # Hadith JSON data files
├── Dockerfile
├── .env.example
└── README.md
```

## Prerequisites

- Go 1.21 or later
- Telegram Bot Token (get from @BotFather)

## Installation

### Local Development

1. Clone the repository and navigate to the bot directory:
```bash
cd telegram-bot
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and add your Telegram Bot Token:
```
TELEGRAM_BOT_TOKEN=your_actual_token_here
```

4. Run the bot:
```bash
go run ./cmd/bot
```

### Docker

1. Build the Docker image:
```bash
docker build -t hadith-bot .
```

2. Run the container:
```bash
docker run -d \
  --name hadith-bot \
  -e TELEGRAM_BOT_TOKEN=your_bot_token \
  hadith-bot
```

### Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  hadith-bot:
    build: .
    container_name: hadith-bot
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - LOG_LEVEL=info
    restart: unless-stopped
```

Run with:
```bash
docker-compose up -d
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEGRAM_BOT_TOKEN` | Your Telegram bot token | (required) |
| `API_URL` | External API URL | `https://api.sunnah.com/v1` |
| `API_KEY` | API key for external service | `""` |
| `API_TIMEOUT` | API request timeout | `10s` |
| `RATE_LIMIT_REQUESTS` | Max requests per window | `10` |
| `RATE_LIMIT_WINDOW` | Rate limit window | `1m` |
| `LOG_LEVEL` | Logging level | `info` |

## Architecture

The bot follows clean architecture principles:

1. **Handler Layer** (`internal/bot/`): Handles Telegram commands and callbacks
2. **Service Layer** (`internal/services/`): Business logic for hadith operations
3. **Data Layer** (`internal/data/`): Data loading and parsing
4. **Model Layer** (`internal/models/`): Data structures

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

- No global state - all dependencies injected
- Environment-based configuration
- Structured logging
- Clean separation of concerns

## License

MIT License

