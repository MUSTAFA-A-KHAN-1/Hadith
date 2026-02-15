# Migration Task: go-telegram-bot-api/v5

## Steps to Complete:

- [ ] 1. Update go.mod - Replace telebot with go-telegram-bot-api/v5
- [ ] 2. Rewrite main.go - Use tgbotapi.NewBotAPI, tgbotapi.NewUpdate, GetUpdatesChan
- [ ] 3. Rewrite handlers.go - Implement command, inline query, callback handlers
- [ ] 4. Create utils/markdown.go - Add MarkdownV2 escape helper
- [ ] 5. Update .env.example - Add example environment variables
- [ ] 6. Update Dockerfile - Ensure proper configuration

## Notes:
- Keep clean architecture (main.go, handlers/, services/, utils/)
- No global mutable state
- Use context timeouts for API calls
- Ensure inline queries respond in under 2 seconds

