package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hadith-bot/internal/logger"
	"hadith-bot/internal/models"
	"hadith-bot/internal/services"
	"hadith-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// strPtr returns a pointer to a string
func strPtr(s string) *string {
	return &s
}

// Handler handles telegram bot commands and callbacks
type Handler struct {
	bot           *tgbotapi.BotAPI
	hadithService *services.HadithService
	log           *logger.Logger
	rateLimiter   *RateLimiter
}

// NewHandler creates a new handler
func NewHandler(bot *tgbotapi.BotAPI, hadithService *services.HadithService, log *logger.Logger, rateLimitRequests int, rateLimitWindow time.Duration) *Handler {
	return &Handler{
		bot:           bot,
		hadithService: hadithService,
		log:           log,
		rateLimiter:   NewRateLimiter(rateLimitRequests, rateLimitWindow),
	}
}

// HandleUpdates handles incoming updates from the channel
func (h *Handler) HandleUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			h.log.Info("Stopping update handler")
			return
		case update, ok := <-updates:
			if !ok {
				h.log.Info("Updates channel closed")
				return
			}
			h.handleUpdate(update)
		}
	}
}

// handleUpdate processes a single update
func (h *Handler) handleUpdate(update tgbotapi.Update) {
	// Handle inline queries
	if update.InlineQuery != nil {
		h.handleInlineQuery(update.InlineQuery)
		return
	}

	// Handle callback queries
	if update.CallbackQuery != nil {
		h.handleCallback(update.CallbackQuery)
		return
	}

	// Handle regular messages
	if update.Message != nil {
		h.handleMessage(update.Message)
	}
}

// handleMessage processes incoming messages
func (h *Handler) handleMessage(m *tgbotapi.Message) {
	if !m.IsCommand() {
		return
	}

	chatID := m.Chat.ID
	userID := m.From.ID

	if !h.rateLimiter.Allow(int64(userID)) {
		h.sendMessage(chatID, "Please wait a moment before sending another command.")
		return
	}

	switch m.Command() {
	case "start":
		h.handleStart(m)
	case "help":
		h.handleHelp(m)
	case "random":
		h.handleRandom(m)
	case "search":
		h.handleSearch(m)
	case "collections":
		h.handleCollections(m)
	}
}

// handleStart handles /start command
func (h *Handler) handleStart(m *tgbotapi.Message) {
	welcomeText := `*Welcome to Hadith Portal Bot* 🕌

This bot provides access to authentic hadith collections from the six major books of hadith.

*Available Commands:*

• /start - Start the bot
• /collections - Browse hadith collections
• /search <keyword> - Search hadiths
• /random - Get a random hadith
• /help - Get help

Use the inline keyboard below to navigate:`

	// Create inline keyboard
	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "📚 Browse Collections", CallbackData: strPtr("collections:1")},
			{Text: "🔍 Search Hadith", CallbackData: strPtr("search")},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🎲 Random Hadith", CallbackData: strPtr("random")},
			{Text: "❓ Help", CallbackData: strPtr("help")},
		},
	)

	h.sendMessageWithKeyboard(m.Chat.ID, welcomeText, &menu)
}

// handleHelp handles /help command
func (h *Handler) handleHelp(m *tgbotapi.Message) {
	helpText := `*Hadith Portal Bot - Help* ❓

*Commands:*

/start - Welcome message and main menu
/collections - Browse all hadith collections
/search <keyword> - Search for hadiths
/random - Get a random hadith
/help - Show this help message

*How to Search:*
Use /search followed by your keyword
Example: /search prayer

*Tips:*
• Search results are limited to first 5 matches
• Use pagination buttons to see more results
• Click on a book to view hadiths

*Collections Available:*
• Sahih al-Bukhari
• Sahih Muslim
• Sunan Abu Dawood
• Jami' at-Tirmidhi
• Sunan an-Nasa'i
• Sunan Ibn Majah`

	h.sendMessage(m.Chat.ID, helpText)
}

// handleRandom handles /random command
func (h *Handler) handleRandom(m *tgbotapi.Message) {
	chatID := m.Chat.ID
	userID := m.From.ID

	h.log.LogRequest(int64(m.MessageID), int64(userID), "/random")

	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		h.sendMessage(chatID, "Sorry, couldn't fetch a random hadith. Please try again.")
		h.log.LogResponse(int64(userID), "random", false)
		return
	}

	hadithText := h.formatHadithDisplay(result.Hadith, result.Collection, result.Book)
	h.sendMessage(chatID, hadithText)
	h.log.LogResponse(int64(userID), "random", true)
}

// handleSearch handles /search command
func (h *Handler) handleSearch(m *tgbotapi.Message) {
	chatID := m.Chat.ID
	userID := m.From.ID

	// Extract query from command arguments
	args := m.CommandArguments()

	if args == "" {
		h.sendMessage(chatID, "Please provide a search keyword.\nUsage: /search <keyword>")
		return
	}

	h.log.LogRequest(int64(m.MessageID), int64(userID), "/search "+args)

	results := h.hadithService.SearchHadiths(args, 1, 10)

	if len(results.Hadiths) == 0 {
		h.sendMessage(chatID, fmt.Sprintf("No results found for: *%s*", utils.EscapeMarkdownV2(args)))
		return
	}

	// Create inline keyboard with search results
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for _, hadith := range results.Hadiths {
		// Find which collection this hadith belongs to
		collectionName := h.findCollectionForHadith(hadith)
		data := fmt.Sprintf("hadith_search:%s:%d", collectionName, hadith.HadithNumber)
		grade := hadith.Grade
		if grade == "" {
			grade = "Sahih"
		}
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), CallbackData: strPtr(data)},
		})
	}

	// Add pagination if needed
	if results.TotalPages > 1 {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: "➡️ Next", CallbackData: strPtr(fmt.Sprintf("search_next:%s:2", args))},
		})
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	resultText := fmt.Sprintf("🔍 *Search Results for:* %s\n\nFound *%d* results. Showing first 5:\n\n_Click on a hadith to view full details_",
		utils.EscapeMarkdownV2(args), results.Total)

	h.sendMessageWithKeyboard(chatID, resultText, &menu)
	h.log.LogResponse(int64(userID), "search", true)
}

// handleCollections handles /collections command
func (h *Handler) handleCollections(m *tgbotapi.Message) {
	chatID := m.Chat.ID
	userID := m.From.ID

	h.log.LogRequest(int64(m.MessageID), int64(userID), "/collections")

	collections := h.hadithService.GetCollections()
	h.sendCollectionsMenu(chatID, collections, 1)
	h.log.LogResponse(int64(userID), "collections", true)
}

// sendCollectionsMenu sends the collections menu
func (h *Handler) sendCollectionsMenu(chatID int64, collections []models.Collection, page int) {
	const perPage = 6

	start := (page - 1) * perPage
	end := start + perPage
	if end > len(collections) {
		end = len(collections)
	}

	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	for _, c := range collections[start:end] {
		data := fmt.Sprintf("books:%s:1", c.Name)
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: c.Title, CallbackData: strPtr(data)},
		})
	}

	// Add navigation buttons
	if len(collections) > perPage {
		var navRow []tgbotapi.InlineKeyboardButton
		if page > 1 {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{Text: "⬅️ Previous", CallbackData: strPtr(fmt.Sprintf("collections:%d", page-1))})
		}
		if end < len(collections) {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{Text: "Next ➡️", CallbackData: strPtr(fmt.Sprintf("collections:%d", page+1))})
		}
		if len(navRow) > 0 {
			keyboardRows = append(keyboardRows, navRow)
		}
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	text := "📚 *Hadith Collections*\n\nSelect a collection to browse:"
	h.sendMessageWithKeyboard(chatID, text, &menu)
}

// handleCallback handles all callback queries
func (h *Handler) handleCallback(c *tgbotapi.CallbackQuery) {
	userID := c.From.ID

	if !h.rateLimiter.Allow(int64(userID)) {
		h.answerCallback(c, "Please wait a moment...", true)
		return
	}

	h.log.LogCallback(c.ID, int64(userID), c.Data)

	// Parse callback data
	parts := strings.Split(c.Data, ":")
	if len(parts) < 1 {
		return
	}

	// Answer callback to remove loading state
	h.answerCallback(c, "", false)

	switch parts[0] {
	case "collections":
		h.handleCollectionsCallback(c, parts)
	case "books":
		h.handleBooksCallback(c, parts)
	case "hadiths":
		h.handleHadithsCallback(c, parts)
	case "hadith_detail":
		h.handleHadithDetailCallback(c, parts)
	case "hadith_search":
		h.handleHadithSearchCallback(c, parts)
	case "random":
		h.handleRandomCallback(c)
	case "search":
		h.handleSearchPaginationCallback(c, parts)
	case "search_next":
		h.handleSearchNextCallback(c, parts)
	case "search_prev":
		h.handleSearchPrevCallback(c, parts)
	case "help":
		h.handleHelpCallback(c)
	}
}

// handleInlineQuery handles inline queries
// Format: @MyHadithBot random - returns random hadith
//
//	@MyHadithBot search <keyword> - returns search results
func (h *Handler) handleInlineQuery(q *tgbotapi.InlineQuery) {
	// Parse the query
	query := strings.TrimSpace(q.Query)
	query = strings.ToLower(query)

	var results []interface{}

	if query == "random" {
		// Handle random hadith inline query
		results = h.handleInlineRandom(q)
	} else if strings.HasPrefix(query, "search ") {
		// Handle search inline query
		keyword := strings.TrimPrefix(query, "search ")
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			results = h.handleInlineSearch(q, keyword)
		} else {
			// Empty search - show help
			results = h.handleInlineHelp(q)
		}
	} else if query == "" {
		// Empty query - show help
		results = h.handleInlineHelp(q)
	} else {
		// Unknown command - show help
		results = h.handleInlineHelp(q)
	}

	// Send response with cache time of 10 seconds
	resp := tgbotapi.InlineConfig{
		InlineQueryID: q.ID,
		Results:       results,
		CacheTime:     10,
	}

	h.bot.Request(resp)
}

// handleInlineRandom handles inline query for random hadith
func (h *Handler) handleInlineRandom(q *tgbotapi.InlineQuery) []interface{} {
	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		return []interface{}{}
	}

	// Format the hadith for inline display
	hadithText := h.formatHadithForInline(result.Hadith, result.Collection, result.Book)

	// Create unique result ID
	resultID := fmt.Sprintf("random_%d_%d", result.Hadith.HadithNumber, time.Now().UnixNano())

	article := tgbotapi.NewInlineQueryResultArticle(resultID, "🎲 Random Hadith", hadithText)
	article.Description = fmt.Sprintf("Hadith #%d from %s", result.Hadith.HadithNumber, getCollectionTitle(result.Collection))

	return []interface{}{article}
}

// handleInlineSearch handles inline query for search
func (h *Handler) handleInlineSearch(q *tgbotapi.InlineQuery, keyword string) []interface{} {
	results := h.hadithService.SearchHadiths(keyword, 1, 5)

	if len(results.Hadiths) == 0 {
		// No results - return a single result indicating no matches
		resultID := fmt.Sprintf("no_results_%d", time.Now().UnixNano())
		article := tgbotapi.NewInlineQueryResultArticle(resultID, "🔍 No Results Found", fmt.Sprintf("No results found for *%s*", utils.EscapeMarkdownV2(keyword)))
		article.Description = fmt.Sprintf("No hadiths found for: %s", keyword)
		return []interface{}{article}
	}

	var inlineResults []interface{}

	for i, hadith := range results.Hadiths {
		// Find the collection for this hadith
		collectionName := h.findCollectionForHadith(hadith)
		collection := h.hadithService.GetCollection(collectionName)

		// Format hadith text
		hadithText := h.formatHadithForInline(&hadith, collection, nil)

		// Create unique result ID
		resultID := fmt.Sprintf("search_%s_%d_%d_%d", keyword, hadith.HadithNumber, i, time.Now().UnixNano())

		grade := hadith.Grade
		if grade == "" {
			grade = "Sahih"
		}

		description := truncate(hadith.English, 50)

		article := tgbotapi.NewInlineQueryResultArticle(resultID, fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), hadithText)
		article.Description = description

		inlineResults = append(inlineResults, article)
	}

	return inlineResults
}

// handleInlineHelp handles inline query help/empty query
func (h *Handler) handleInlineHelp(q *tgbotapi.InlineQuery) []interface{} {
	resultID := fmt.Sprintf("help_%d", time.Now().UnixNano())
	text := "🕌 *Hadith Portal Bot*\\n\\nUse inline mode:\\n\\n• @MyHadithBot *random* - Get a random hadith\\n• @MyHadithBot *search <keyword>* - Search hadiths\\n\\nExample: @MyHadithBot search prayer"

	article := tgbotapi.NewInlineQueryResultArticle(resultID, "🕌 Hadith Portal Bot", text)
	article.Description = "Use @MyHadithBot random or @MyHadithBot search <keyword>"

	return []interface{}{article}
}

// formatHadithForInline formats a hadith for inline query display with MarkdownV2
func (h *Handler) formatHadithForInline(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	collectionName := "Unknown Collection"
	if collection != nil {
		collectionName = collection.Title
	}

	bookNumber := 0
	if book != nil {
		bookNumber = book.BookNumber
	}

	grade := hadith.Grade
	if grade == "" {
		grade = "Sahih"
	}

	// Build the formatted message using MarkdownV2
	var sb strings.Builder
	sb.WriteString("🵿 *Hadith*\\n\\n")
	sb.WriteString("*Arabic:*\\n")
	sb.WriteString(utils.EscapeMarkdownV2(hadith.Arabic))
	sb.WriteString("\\n\\n")
	sb.WriteString("*English:*\\n")
	sb.WriteString(utils.EscapeMarkdownV2(hadith.English))
	sb.WriteString("\\n\\n")
	sb.WriteString("*Reference:* ")
	sb.WriteString(utils.EscapeMarkdownV2(collectionName))
	sb.WriteString(", Book ")
	sb.WriteString(fmt.Sprintf("%d", bookNumber))
	sb.WriteString(", Hadith #")
	sb.WriteString(fmt.Sprintf("%d", hadith.HadithNumber))
	sb.WriteString("\\n")
	sb.WriteString("*Grade:* ")
	sb.WriteString(utils.EscapeMarkdownV2(grade))
	sb.WriteString("\\n")

	return sb.String()
}

// getCollectionTitle returns the title of a collection
func getCollectionTitle(c *models.Collection) string {
	if c == nil {
		return "Unknown"
	}
	return c.Title
}

// answerCallback answers a callback query
func (h *Handler) answerCallback(c *tgbotapi.CallbackQuery, text string, showAlert bool) {
	resp := tgbotapi.CallbackConfig{
		CallbackQueryID: c.ID,
		Text:            text,
		ShowAlert:       showAlert,
	}
	h.bot.Request(resp)
}

// handleHelpCallback handles help from inline keyboard
func (h *Handler) handleHelpCallback(c *tgbotapi.CallbackQuery) {
	helpText := `*Hadith Portal Bot - Help* ❓

*Commands:*

/start - Welcome message and main menu
/collections - Browse all hadith collections
/search <keyword> - Search for hadiths
/random - Get a random hadith
/help - Show this help message

*How to Search:*
Use /search followed by your keyword
Example: /search prayer

*Tips:*
• Search results are limited to first 5 matches
• Use pagination buttons to see more results
• Click on a book to view hadiths`

	msg := tgbotapi.NewEditMessageText(c.Message.Chat.ID, c.Message.MessageID, helpText)
	msg.ParseMode = "MarkdownV2"
	h.bot.Request(msg)
}

// handleCollectionsCallback handles collection navigation
func (h *Handler) handleCollectionsCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	page, err := strconv.Atoi(parts[1])
	if err != nil {
		page = 1
	}

	collections := h.hadithService.GetCollections()
	h.sendCollectionsMenu(c.Message.Chat.ID, collections, page)
}

// handleBooksCallback handles books display
func (h *Handler) handleBooksCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	collectionName := parts[1]
	page, _ := strconv.Atoi(parts[2])

	books := h.hadithService.GetBooks(collectionName)
	if len(books) == 0 {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "No books found in this collection.")
		return
	}

	h.sendBooksMenu(c.Message.Chat.ID, collectionName, books, page)
}

// sendBooksMenu sends the books menu for a collection
func (h *Handler) sendBooksMenu(chatID int64, collection string, books []models.Book, page int) {
	const perPage = 10

	start := (page - 1) * perPage
	end := start + perPage
	if end > len(books) {
		end = len(books)
	}

	collectionName := services.GetCollectionDisplayName(collection)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	for _, b := range books[start:end] {
		data := fmt.Sprintf("hadiths:%s:%d:1", collection, b.BookNumber)
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: fmt.Sprintf("📖 %s", truncate(b.Title, 40)), CallbackData: strPtr(data)},
		})
	}

	// Add navigation buttons
	if len(books) > perPage {
		var navRow []tgbotapi.InlineKeyboardButton
		if page > 1 {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{Text: "⬅️ Prev", CallbackData: strPtr(fmt.Sprintf("books:%s:%d", collection, page-1))})
		}
		if end < len(books) {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{Text: "Next ➡️", CallbackData: strPtr(fmt.Sprintf("books:%s:%d", collection, page+1))})
		}
		if len(navRow) > 0 {
			keyboardRows = append(keyboardRows, navRow)
		}
	}

	// Back to collections button
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		{Text: "⬅️ Back to Collections", CallbackData: strPtr("collections:1")},
	})

	menu := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	text := fmt.Sprintf("📚 *%s*\n\nSelect a book:", collectionName)
	h.sendMessageWithKeyboard(chatID, text, &menu)
}

// handleHadithsCallback handles hadiths display
func (h *Handler) handleHadithsCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 4 {
		return
	}

	collectionName := parts[1]
	bookNumber, _ := strconv.Atoi(parts[2])
	page, _ := strconv.Atoi(parts[3])

	result := h.hadithService.GetHadiths(collectionName, bookNumber, page, 10)
	h.sendHadithsMenu(c.Message.Chat.ID, collectionName, bookNumber, result)
}

// sendHadithsMenu sends the hadiths menu for a book
func (h *Handler) sendHadithsMenu(chatID int64, collection string, bookNumber int, result models.HadithResponse) {
	collectionName := services.GetCollectionDisplayName(collection)
	book := h.hadithService.GetBook(collection, bookNumber)
	bookTitle := "Unknown"
	if book != nil {
		bookTitle = book.Title
	}

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for i, hadith := range result.Hadiths {
		data := fmt.Sprintf("hadith_detail:%s:%d:%d:%d", collection, bookNumber, result.Page, i)
		grade := hadith.Grade
		if grade == "" {
			grade = "Sahih"
		}
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), CallbackData: strPtr(data)},
		})
	}

	// Add navigation buttons
	if result.TotalPages > 1 {
		var navRow []tgbotapi.InlineKeyboardButton
		if result.Page > 1 {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{
				Text:         "⬅️ Prev",
				CallbackData: strPtr(fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page-1)),
			})
		}
		if result.Page < result.TotalPages {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: strPtr(fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page+1)),
			})
		}
		if len(navRow) > 0 {
			keyboardRows = append(keyboardRows, navRow)
		}
	}

	// Back button
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		{
			Text:         "⬅️ Back",
			CallbackData: strPtr(fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page)),
		},
	})

	menu := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	text := fmt.Sprintf("📖 *%s*\n📑 %s\n\nPage %d/%d - Showing %d hadiths:",
		collectionName, truncate(bookTitle, 30), result.Page, result.TotalPages, len(result.Hadiths))
	h.sendMessageWithKeyboard(chatID, text, &menu)
}

// handleHadithDetailCallback handles hadith detail display
func (h *Handler) handleHadithDetailCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 5 {
		return
	}

	collectionName := parts[1]
	bookNumber, _ := strconv.Atoi(parts[2])
	page, _ := strconv.Atoi(parts[3])
	index, _ := strconv.Atoi(parts[4])

	result := h.hadithService.GetHadiths(collectionName, bookNumber, page, 10)

	if index >= len(result.Hadiths) {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "Hadith not found.")
		return
	}

	hadith := result.Hadiths[index]
	collection := h.hadithService.GetCollection(collectionName)
	book := h.hadithService.GetBook(collectionName, bookNumber)

	hadithText := h.formatHadithDisplay(&hadith, collection, book)

	// Create menu for navigation
	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🎲 Random Hadith", CallbackData: strPtr("random")},
			{Text: "🔍 Search", CallbackData: strPtr("search")},
		},
		[]tgbotapi.InlineKeyboardButton{
			{Text: "⬅️ Back", CallbackData: strPtr(fmt.Sprintf("hadiths:%s:%d:%d", collectionName, bookNumber, page))},
		},
	)

	h.editMessageTextWithKeyboard(c.Message.Chat.ID, c.Message.MessageID, hadithText, &menu)
}

// handleHadithSearchCallback handles hadith search detail display
func (h *Handler) handleHadithSearchCallback(c *tgbotapi.CallbackQuery, parts []string) {
	// Format: hadith_search:{collectionName}:{hadithNumber}
	if len(parts) < 3 {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "Invalid hadith selection.")
		return
	}

	collectionName := parts[1]
	hadithNumber, err := strconv.Atoi(parts[2])
	if err != nil {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "Invalid hadith selection.")
		return
	}

	// Find the hadith - we need to search through books to find it
	books := h.hadithService.GetBooks(collectionName)
	var foundHadith *models.Hadith
	var foundBook *models.Book

	for _, book := range books {
		// Get hadiths for this book
		result := h.hadithService.GetHadiths(collectionName, book.BookNumber, 1, 100)
		for i := range result.Hadiths {
			if result.Hadiths[i].HadithNumber == hadithNumber {
				foundHadith = &result.Hadiths[i]
				foundBook = &book
				break
			}
		}
		if foundHadith != nil {
			break
		}
	}

	if foundHadith == nil {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "Hadith not found.")
		return
	}

	collection := h.hadithService.GetCollection(collectionName)
	hadithText := h.formatHadithDisplay(foundHadith, collection, foundBook)

	// Create menu for navigation
	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🎲 Random Hadith", CallbackData: strPtr("random")},
			{Text: "🔍 Search", CallbackData: strPtr("search")},
		},
	)

	h.editMessageTextWithKeyboard(c.Message.Chat.ID, c.Message.MessageID, hadithText, &menu)
}

// handleRandomCallback handles random hadith from callback
func (h *Handler) handleRandomCallback(c *tgbotapi.CallbackQuery) {
	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		h.editMessageText(c.Message.Chat.ID, c.Message.MessageID, "Sorry, couldn't fetch a random hadith.")
		return
	}

	hadithText := h.formatHadithDisplay(result.Hadith, result.Collection, result.Book)

	menu := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "🎲 Another Random", CallbackData: strPtr("random")},
			{Text: "🔍 Search", CallbackData: strPtr("search")},
		},
	)

	h.editMessageTextWithKeyboard(c.Message.Chat.ID, c.Message.MessageID, hadithText, &menu)
}

// handleSearchPaginationCallback handles search pagination
func (h *Handler) handleSearchPaginationCallback(c *tgbotapi.CallbackQuery, parts []string) {
	h.sendMessage(c.Message.Chat.ID, "Please use /search <keyword> command to search.")
}

// handleSearchNextCallback handles search next page
func (h *Handler) handleSearchNextCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	query := parts[1]
	page, _ := strconv.Atoi(parts[2])

	results := h.hadithService.SearchHadiths(query, page, 5)
	h.sendSearchResults(c.Message.Chat.ID, query, results)
}

// handleSearchPrevCallback handles search previous page
func (h *Handler) handleSearchPrevCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	query := parts[1]
	page, _ := strconv.Atoi(parts[2])

	if page > 1 {
		page--
	}

	results := h.hadithService.SearchHadiths(query, page, 5)
	h.sendSearchResults(c.Message.Chat.ID, query, results)
}

// sendSearchResults sends search results
func (h *Handler) sendSearchResults(chatID int64, query string, results models.SearchResult) {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for i, hadith := range results.Hadiths {
		data := fmt.Sprintf("hadith_search:%d:%d", results.Page, i)
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			{Text: fmt.Sprintf("Hadith #%d", hadith.HadithNumber), CallbackData: strPtr(data)},
		})
	}

	// Add pagination
	if results.TotalPages > 1 {
		var navRow []tgbotapi.InlineKeyboardButton
		if results.Page > 1 {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{
				Text:         "⬅️ Prev",
				CallbackData: strPtr(fmt.Sprintf("search_prev:%s:%d", query, results.Page-1)),
			})
		}
		if results.Page < results.TotalPages {
			navRow = append(navRow, tgbotapi.InlineKeyboardButton{
				Text:         "Next ➡️",
				CallbackData: strPtr(fmt.Sprintf("search_next:%s:%d", query, results.Page+1)),
			})
		}
		if len(navRow) > 0 {
			keyboardRows = append(keyboardRows, navRow)
		}
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	resultText := fmt.Sprintf("🔍 *Search Results for:* %s\n\nPage %d/%d - Showing %d of %d results:",
		utils.EscapeMarkdownV2(query), results.Page, results.TotalPages, len(results.Hadiths), results.Total)

	h.sendMessageWithKeyboard(chatID, resultText, &menu)
}

// formatHadithDisplay formats a hadith for display
func (h *Handler) formatHadithDisplay(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	collectionName := "Unknown Collection"
	if collection != nil {
		collectionName = collection.Title
	}

	bookNumber := 0
	if book != nil {
		bookNumber = book.BookNumber
	}

	grade := hadith.Grade
	if grade == "" {
		grade = "Sahih"
	}

	// Build the formatted message
	var sb strings.Builder
	sb.WriteString("🵿 *Hadith*\n\n")
	sb.WriteString("*Arabic:*\n")
	sb.WriteString(utils.EscapeMarkdownV2(hadith.Arabic))
	sb.WriteString("\n\n*English:*\n")
	sb.WriteString(utils.EscapeMarkdownV2(hadith.English))
	sb.WriteString("\n\n")
	sb.WriteString("*Reference:* ")
	sb.WriteString(utils.EscapeMarkdownV2(collectionName))
	sb.WriteString(", Book ")
	sb.WriteString(fmt.Sprintf("%d", bookNumber))
	sb.WriteString(", Hadith #")
	sb.WriteString(fmt.Sprintf("%d", hadith.HadithNumber))
	sb.WriteString("\n")
	sb.WriteString("*Grade:* ")
	sb.WriteString(utils.EscapeMarkdownV2(grade))
	sb.WriteString("\n")

	return sb.String()
}

// sendMessage sends a text message
func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	_, err := h.bot.Send(msg)
	if err != nil {
		h.log.LogError(err, "sendMessage")
	}
}

// sendMessageWithKeyboard sends a text message with keyboard
func (h *Handler) sendMessageWithKeyboard(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = keyboard
	_, err := h.bot.Send(msg)
	if err != nil {
		h.log.LogError(err, "sendMessageWithKeyboard")
	}
}

// editMessageText edits a message text
func (h *Handler) editMessageText(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "MarkdownV2"
	h.bot.Request(msg)
}

// editMessageTextWithKeyboard edits a message text with keyboard
func (h *Handler) editMessageTextWithKeyboard(chatID int64, messageID int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = keyboard
	h.bot.Request(msg)
}

// truncate truncates a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// findCollectionForHadith finds which collection a hadith belongs to based on its chapterId
func (h *Handler) findCollectionForHadith(hadith models.Hadith) string {
	collections := h.hadithService.GetCollections()
	for _, collection := range collections {
		books := h.hadithService.GetBooks(collection.Name)
		for _, book := range books {
			if book.BookNumber == hadith.ChapterID {
				return collection.Name
			}
		}
	}
	// Default to bukhari if not found
	return "bukhari"
}
