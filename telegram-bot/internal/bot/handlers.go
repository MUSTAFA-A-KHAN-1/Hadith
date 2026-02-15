package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"hadith-bot/internal/logger"
	"hadith-bot/internal/models"
	"hadith-bot/internal/services"

	telebot "github.com/tucnak/telebot"
)

// Handler handles telegram bot commands and callbacks
type Handler struct {
	bot           *telebot.Bot
	hadithService *services.HadithService
	log           *logger.Logger
	rateLimiter   *RateLimiter
}

// NewHandler creates a new handler
func NewHandler(bot *telebot.Bot, hadithService *services.HadithService, log *logger.Logger, rateLimitRequests int, rateLimitWindow time.Duration) *Handler {
	return &Handler{
		bot:           bot,
		hadithService: hadithService,
		log:           log,
		rateLimiter:   NewRateLimiter(rateLimitRequests, rateLimitWindow),
	}
}

// HandleCommands sets up all command handlers
func (h *Handler) HandleCommands() {
	// Command: /start
	h.bot.Handle("/start", h.handleStart)

	// Command: /help
	h.bot.Handle("/help", h.handleHelp)

	// Command: /random
	h.bot.Handle("/random", h.handleRandom)

	// Command: /search
	h.bot.Handle("/search", h.handleSearch)

	// Command: /collections
	h.bot.Handle("/collections", h.handleCollections)

	// Inline query handler
	h.bot.Handle(telebot.OnQuery, h.handleInlineQuery)

	// Callback queries
	h.bot.Handle(telebot.OnCallback, h.handleCallback)
}

// handleStart handles /start command
func (h *Handler) handleStart(m *telebot.Message) {
	if !h.rateLimiter.Allow(int64(m.Sender.ID)) {
		h.sendMessage(m.Sender, "Please wait a moment before sending another command.")
		return
	}

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
	menu := &telebot.ReplyMarkup{}

	menu.InlineKeyboard = [][]telebot.InlineButton{
		{
			{Text: "📚 Browse Collections", Data: "collections:1"},
			{Text: "🔍 Search Hadith", Data: "search"},
		},
		{
			{Text: "🎲 Random Hadith", Data: "random"},
			{Text: "❓ Help", Data: "help"},
		},
	}

	h.sendMessageWithKeyboard(m.Sender, welcomeText, menu)
}

// handleHelp handles /help command
func (h *Handler) handleHelp(m *telebot.Message) {
	if !h.rateLimiter.Allow(int64(m.Sender.ID)) {
		h.sendMessage(m.Sender, "Please wait a moment before sending another command.")
		return
	}

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

	h.sendMessage(m.Sender, helpText)
}

// handleRandom handles /random command
func (h *Handler) handleRandom(m *telebot.Message) {
	if !h.rateLimiter.Allow(int64(m.Sender.ID)) {
		h.sendMessage(m.Sender, "Please wait a moment before sending another command.")
		return
	}

	h.log.LogRequest(int64(m.ID), int64(m.Sender.ID), "/random")

	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		h.sendMessage(m.Sender, "Sorry, couldn't fetch a random hadith. Please try again.")
		h.log.LogResponse(int64(m.Sender.ID), "random", false)
		return
	}

	hadithText := h.formatHadithDisplay(result.Hadith, result.Collection, result.Book)
	h.sendMessage(m.Sender, hadithText)
	h.log.LogResponse(int64(m.Sender.ID), "random", true)
}

// handleSearch handles /search command
func (h *Handler) handleSearch(m *telebot.Message) {
	if !h.rateLimiter.Allow(int64(m.Sender.ID)) {
		h.sendMessage(m.Sender, "Please wait a moment before sending another command.")
		return
	}

	// Extract query from command
	args := strings.TrimPrefix(m.Text, "/search")
	args = strings.TrimSpace(args)

	if args == "" {
		h.sendMessage(m.Sender, "Please provide a search keyword.\nUsage: /search <keyword>")
		return
	}

	h.log.LogRequest(int64(m.ID), int64(m.Sender.ID), "/search "+args)

	results := h.hadithService.SearchHadiths(args, 1, 10)

	if len(results.Hadiths) == 0 {
		h.sendMessage(m.Sender, fmt.Sprintf("No results found for: *%s*", args))
		return
	}

	// Create search results message
	menu := &telebot.ReplyMarkup{}
	var keyboard [][]telebot.InlineButton

	// Search across all collections and track which collection each hadith belongs to
	// Format: hadith_search:{collectionName}:{hadithNumber}
	for _, hadith := range results.Hadiths {
		// Find which collection this hadith belongs to
		collectionName := h.findCollectionForHadith(hadith)
		data := fmt.Sprintf("hadith_search:%s:%d", collectionName, hadith.HadithNumber)
		grade := hadith.Grade
		if grade == "" {
			grade = "Sahih"
		}
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), Data: data},
		})
	}

	// Add pagination if needed
	if results.TotalPages > 1 {
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: "➡️ Next", Data: fmt.Sprintf("search_next:%s:2", args)},
		})
	}

	menu.InlineKeyboard = keyboard

	resultText := fmt.Sprintf("🔍 *Search Results for:* %s\n\nFound *%d* results. Showing first 5:\n\n_Click on a hadith to view full details_",
		args, results.Total)

	h.sendMessageWithKeyboard(m.Sender, resultText, menu)
	h.log.LogResponse(int64(m.Sender.ID), "search", true)
}

// handleCollections handles /collections command
func (h *Handler) handleCollections(m *telebot.Message) {
	if !h.rateLimiter.Allow(int64(m.Sender.ID)) {
		h.sendMessage(m.Sender, "Please wait a moment before sending another command.")
		return
	}

	h.log.LogRequest(int64(m.ID), int64(m.Sender.ID), "/collections")

	collections := h.hadithService.GetCollections()
	h.sendCollectionsMenu(m.Sender, collections, 1)
	h.log.LogResponse(int64(m.Sender.ID), "collections", true)
}

// sendCollectionsMenu sends the collections menu
func (h *Handler) sendCollectionsMenu(user *telebot.User, collections []models.Collection, page int) {
	menu := &telebot.ReplyMarkup{}
	const perPage = 6

	start := (page - 1) * perPage
	end := start + perPage
	if end > len(collections) {
		end = len(collections)
	}

	var keyboard [][]telebot.InlineButton
	for _, c := range collections[start:end] {
		data := fmt.Sprintf("books:%s:1", c.Name)
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: c.Title, Data: data},
		})
	}

	// Add navigation buttons
	if len(collections) > perPage {
		var navRow []telebot.InlineButton
		if page > 1 {
			navRow = append(navRow, telebot.InlineButton{Text: "⬅️ Previous", Data: fmt.Sprintf("collections:%d", page-1)})
		}
		if end < len(collections) {
			navRow = append(navRow, telebot.InlineButton{Text: "Next ➡️", Data: fmt.Sprintf("collections:%d", page+1)})
		}
		if len(navRow) > 0 {
			keyboard = append(keyboard, navRow)
		}
	}

	menu.InlineKeyboard = keyboard

	text := "📚 *Hadith Collections*\n\nSelect a collection to browse:"
	h.sendMessageWithKeyboard(user, text, menu)
}

// handleCallback handles all callback queries
func (h *Handler) handleCallback(c *telebot.Callback) {
	if !h.rateLimiter.Allow(int64(c.Sender.ID)) {
		h.bot.Respond(c, &telebot.CallbackResponse{Text: "Please wait a moment...", ShowAlert: true})
		return
	}

	h.log.LogCallback(c.ID, int64(c.Sender.ID), c.Data)

	// Parse callback data
	parts := strings.Split(c.Data, ":")
	if len(parts) < 1 {
		return
	}

	// Delete the original message if it's from a command (to clean up the chat)
	// This applies to navigation from inline keyboards
	isNavigation := false
	switch parts[0] {
	case "collections", "books", "hadiths", "search", "search_next", "search_prev", "help":
		isNavigation = true
	}

	if isNavigation && c.Message != nil {
		// Try to delete the original command message
		h.bot.Delete(c.Message)
	}

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

	// Answer callback to remove loading state
	h.bot.Respond(c, nil)
}

// handleInlineQuery handles inline queries
// Format: @MyHadithBot random - returns random hadith
//
//	@MyHadithBot search <keyword> - returns search results
func (h *Handler) handleInlineQuery(q *telebot.Query) {
	// Parse the query
	query := strings.TrimSpace(q.Text)
	query = strings.ToLower(query)

	var results telebot.Results

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
	h.bot.Answer(q, &telebot.QueryResponse{
		Results:   results,
		CacheTime: 10, // 10 seconds cache
	})
}

// handleInlineRandom handles inline query for random hadith
func (h *Handler) handleInlineRandom(q *telebot.Query) telebot.Results {
	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		return telebot.Results{}
	}

	// Format the hadith for inline display
	hadithText := h.formatHadithForInline(result.Hadith, result.Collection, result.Book)

	// Create unique result ID
	resultID := fmt.Sprintf("random_%d_%d", result.Hadith.HadithNumber, time.Now().UnixNano())

	article := &telebot.ArticleResult{
		ResultBase: telebot.ResultBase{
			ID: resultID,
		},
		Title:       "🎲 Random Hadith",
		Description: fmt.Sprintf("Hadith #%d from %s", result.Hadith.HadithNumber, getCollectionTitle(result.Collection)),
		Text:        hadithText,
	}

	return telebot.Results{article}
}

// handleInlineSearch handles inline query for search
func (h *Handler) handleInlineSearch(q *telebot.Query, keyword string) telebot.Results {
	results := h.hadithService.SearchHadiths(keyword, 1, 5)

	if len(results.Hadiths) == 0 {
		// No results - return a single result indicating no matches
		article := &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID: fmt.Sprintf("no_results_%d", time.Now().UnixNano()),
			},
			Title:       "🔍 No Results Found",
			Description: fmt.Sprintf("No hadiths found for: %s", keyword),
			Text:        fmt.Sprintf("No results found for \\*%s\\*", h.escapeMarkdown(keyword)),
		}
		return telebot.Results{article}
	}

	var inlineResults telebot.Results

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

		article := &telebot.ArticleResult{
			ResultBase: telebot.ResultBase{
				ID: resultID,
			},
			Title:       fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade),
			Description: description,
			Text:        hadithText,
		}

		inlineResults = append(inlineResults, article)
	}

	return inlineResults
}

// handleInlineHelp handles inline query help/empty query
func (h *Handler) handleInlineHelp(q *telebot.Query) telebot.Results {
	article := &telebot.ArticleResult{
		ResultBase: telebot.ResultBase{
			ID: fmt.Sprintf("help_%d", time.Now().UnixNano()),
		},
		Title:       "🕌 Hadith Portal Bot",
		Description: "Use @MyHadithBot random or @MyHadithBot search <keyword>",
		Text:        "🕌 *Hadith Portal Bot*\\n\\nUse inline mode:\\n\\n• @MyHadithBot \\*random\\* - Get a random hadith\\n• @MyHadithBot \\*search <keyword>\\* - Search hadiths\\n\\nExample: @MyHadithBot search prayer",
	}

	return telebot.Results{article}
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
	sb.WriteString(h.escapeMarkdown(hadith.Arabic))
	sb.WriteString("\\n\\n")
	sb.WriteString("*English:*\\n")
	sb.WriteString(h.escapeMarkdown(hadith.English))
	sb.WriteString("\\n\\n")
	sb.WriteString("*Reference:* ")
	sb.WriteString(h.escapeMarkdown(collectionName))
	sb.WriteString(", Book ")
	sb.WriteString(fmt.Sprintf("%d", bookNumber))
	sb.WriteString(", Hadith #")
	sb.WriteString(fmt.Sprintf("%d", hadith.HadithNumber))
	sb.WriteString("\\n")
	sb.WriteString("*Grade:* ")
	sb.WriteString(h.escapeMarkdown(grade))
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

// handleHelpCallback handles help from inline keyboard
func (h *Handler) handleHelpCallback(c *telebot.Callback) {
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

	h.bot.Edit(c.Message, helpText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}

// handleCollectionsCallback handles collection navigation
func (h *Handler) handleCollectionsCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 2 {
		return
	}

	page, err := strconv.Atoi(parts[1])
	if err != nil {
		page = 1
	}

	collections := h.hadithService.GetCollections()
	h.sendCollectionsMenu(c.Sender, collections, page)
}

// handleBooksCallback handles books display
func (h *Handler) handleBooksCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 3 {
		return
	}

	collectionName := parts[1]
	page, _ := strconv.Atoi(parts[2])

	books := h.hadithService.GetBooks(collectionName)
	if len(books) == 0 {
		h.bot.Edit(c.Message, &telebot.Message{Text: "No books found in this collection."})
		return
	}

	h.sendBooksMenu(c.Sender, collectionName, books, page)
}

// sendBooksMenu sends the books menu for a collection
func (h *Handler) sendBooksMenu(user *telebot.User, collection string, books []models.Book, page int) {
	menu := &telebot.ReplyMarkup{}
	const perPage = 10

	start := (page - 1) * perPage
	end := start + perPage
	if end > len(books) {
		end = len(books)
	}

	collectionName := services.GetCollectionDisplayName(collection)

	var keyboard [][]telebot.InlineButton
	for _, b := range books[start:end] {
		data := fmt.Sprintf("hadiths:%s:%d:1", collection, b.BookNumber)
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: fmt.Sprintf("📖 %s", truncate(b.Title, 40)), Data: data},
		})
	}

	// Add navigation buttons
	if len(books) > perPage {
		var navRow []telebot.InlineButton
		if page > 1 {
			navRow = append(navRow, telebot.InlineButton{Text: "⬅️ Prev", Data: fmt.Sprintf("books:%s:%d", collection, page-1)})
		}
		if end < len(books) {
			navRow = append(navRow, telebot.InlineButton{Text: "Next ➡️", Data: fmt.Sprintf("books:%s:%d", collection, page+1)})
		}
		if len(navRow) > 0 {
			keyboard = append(keyboard, navRow)
		}
	}

	// Back to collections button
	keyboard = append(keyboard, []telebot.InlineButton{
		{Text: "⬅️ Back to Collections", Data: "collections:1"},
	})

	menu.InlineKeyboard = keyboard

	text := fmt.Sprintf("📚 *%s*\n\nSelect a book:", collectionName)
	h.sendMessageWithKeyboard(user, text, menu)
}

// handleHadithsCallback handles hadiths display
func (h *Handler) handleHadithsCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 4 {
		return
	}

	collectionName := parts[1]
	bookNumber, _ := strconv.Atoi(parts[2])
	page, _ := strconv.Atoi(parts[3])

	result := h.hadithService.GetHadiths(collectionName, bookNumber, page, 10)
	h.sendHadithsMenu(c.Sender, collectionName, bookNumber, result)
}

// sendHadithsMenu sends the hadiths menu for a book
func (h *Handler) sendHadithsMenu(user *telebot.User, collection string, bookNumber int, result models.HadithResponse) {
	collectionName := services.GetCollectionDisplayName(collection)
	book := h.hadithService.GetBook(collection, bookNumber)
	bookTitle := "Unknown"
	if book != nil {
		bookTitle = book.Title
	}

	menu := &telebot.ReplyMarkup{}
	var keyboard [][]telebot.InlineButton

	for i, hadith := range result.Hadiths {
		data := fmt.Sprintf("hadith_detail:%s:%d:%d:%d", collection, bookNumber, result.Page, i)
		grade := hadith.Grade
		if grade == "" {
			grade = "Sahih"
		}
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), Data: data},
		})
	}

	// Add navigation buttons
	if result.TotalPages > 1 {
		var navRow []telebot.InlineButton
		if result.Page > 1 {
			navRow = append(navRow, telebot.InlineButton{Text: "⬅️ Prev", Data: fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page-1)})
		}
		if result.Page < result.TotalPages {
			navRow = append(navRow, telebot.InlineButton{Text: "Next ➡️", Data: fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page+1)})
		}
		if len(navRow) > 0 {
			keyboard = append(keyboard, navRow)
		}
	}

	// Back button
	keyboard = append(keyboard, []telebot.InlineButton{
		{Text: "⬅️ Back to Books", Data: fmt.Sprintf("books:%s:1", collection)},
	})

	menu.InlineKeyboard = keyboard

	text := fmt.Sprintf("📖 *%s*\n📑 %s\n\nPage %d/%d - Showing %d hadiths:",
		collectionName, truncate(bookTitle, 30), result.Page, result.TotalPages, len(result.Hadiths))
	h.sendMessageWithKeyboard(user, text, menu)
}

// handleHadithDetailCallback handles hadith detail display
func (h *Handler) handleHadithDetailCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 5 {
		return
	}

	collectionName := parts[1]
	bookNumber, _ := strconv.Atoi(parts[2])
	page, _ := strconv.Atoi(parts[3])
	index, _ := strconv.Atoi(parts[4])

	result := h.hadithService.GetHadiths(collectionName, bookNumber, page, 10)

	if index >= len(result.Hadiths) {
		h.bot.Edit(c.Message, &telebot.Message{Text: "Hadith not found."})
		return
	}

	hadith := result.Hadiths[index]
	collection := h.hadithService.GetCollection(collectionName)
	book := h.hadithService.GetBook(collectionName, bookNumber)

	hadithText := h.formatHadithDisplay(&hadith, collection, book)

	// Create menu for navigation
	menu := &telebot.ReplyMarkup{}
	menu.InlineKeyboard = [][]telebot.InlineButton{
		{
			{Text: "🎲 Random Hadith", Data: "random"},
			{Text: "🔍 Search", Data: "search"},
		},
		{
			{Text: "⬅️ Back", Data: fmt.Sprintf("hadiths:%s:%d:%d", collectionName, bookNumber, page)},
		},
	}

	h.bot.Edit(c.Message, hadithText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown, ReplyMarkup: menu})
}

// handleHadithSearchCallback handles hadith search detail display
func (h *Handler) handleHadithSearchCallback(c *telebot.Callback, parts []string) {
	// Format: hadith_search:{collectionName}:{hadithNumber}
	if len(parts) < 3 {
		h.bot.Edit(c.Message, &telebot.Message{Text: "Invalid hadith selection."})
		return
	}

	collectionName := parts[1]
	hadithNumber, err := strconv.Atoi(parts[2])
	if err != nil {
		h.bot.Edit(c.Message, &telebot.Message{Text: "Invalid hadith selection."})
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
		h.bot.Edit(c.Message, &telebot.Message{Text: "Hadith not found."})
		return
	}

	collection := h.hadithService.GetCollection(collectionName)
	hadithText := h.formatHadithDisplay(foundHadith, collection, foundBook)

	// Create menu for navigation
	menu := &telebot.ReplyMarkup{}
	menu.InlineKeyboard = [][]telebot.InlineButton{
		{
			{Text: "🎲 Random Hadith", Data: "random"},
			{Text: "🔍 Search", Data: "search"},
		},
	}

	h.bot.Edit(c.Message, hadithText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown, ReplyMarkup: menu})
}

// handleRandomCallback handles random hadith from callback
func (h *Handler) handleRandomCallback(c *telebot.Callback) {
	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		h.bot.Edit(c.Message, &telebot.Message{Text: "Sorry, couldn't fetch a random hadith."})
		return
	}

	hadithText := h.formatHadithDisplay(result.Hadith, result.Collection, result.Book)

	menu := &telebot.ReplyMarkup{}
	menu.InlineKeyboard = [][]telebot.InlineButton{
		{
			{Text: "🎲 Another Random", Data: "random"},
			{Text: "🔍 Search", Data: "search"},
		},
	}

	h.bot.Edit(c.Message, hadithText, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown, ReplyMarkup: menu})
}

// handleSearchPaginationCallback handles search pagination
func (h *Handler) handleSearchPaginationCallback(c *telebot.Callback, parts []string) {
	h.sendMessage(c.Sender, "Please use /search <keyword> command to search.")
}

// handleSearchNextCallback handles search next page
func (h *Handler) handleSearchNextCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 3 {
		return
	}

	query := parts[1]
	page, _ := strconv.Atoi(parts[2])

	results := h.hadithService.SearchHadiths(query, page, 5)
	h.sendSearchResults(c.Sender, query, results)
}

// handleSearchPrevCallback handles search previous page
func (h *Handler) handleSearchPrevCallback(c *telebot.Callback, parts []string) {
	if len(parts) < 3 {
		return
	}

	query := parts[1]
	page, _ := strconv.Atoi(parts[2])

	if page > 1 {
		page--
	}

	results := h.hadithService.SearchHadiths(query, page, 5)
	h.sendSearchResults(c.Sender, query, results)
}

// sendSearchResults sends search results
func (h *Handler) sendSearchResults(user *telebot.User, query string, results models.SearchResult) {
	menu := &telebot.ReplyMarkup{}
	var keyboard [][]telebot.InlineButton

	for i, hadith := range results.Hadiths {
		data := fmt.Sprintf("hadith_search:%d:%d", results.Page, i)
		keyboard = append(keyboard, []telebot.InlineButton{
			{Text: fmt.Sprintf("Hadith #%d", hadith.HadithNumber), Data: data},
		})
	}

	// Add pagination
	if results.TotalPages > 1 {
		var navRow []telebot.InlineButton
		if results.Page > 1 {
			navRow = append(navRow, telebot.InlineButton{Text: "⬅️ Prev", Data: fmt.Sprintf("search_prev:%s:%d", query, results.Page)})
		}
		if results.Page < results.TotalPages {
			navRow = append(navRow, telebot.InlineButton{Text: "Next ➡️", Data: fmt.Sprintf("search_next:%s:%d", query, results.Page)})
		}
		if len(navRow) > 0 {
			keyboard = append(keyboard, navRow)
		}
	}

	menu.InlineKeyboard = keyboard

	resultText := fmt.Sprintf("🔍 *Search Results for:* %s\n\nPage %d/%d - Showing %d of %d results:",
		query, results.Page, results.TotalPages, len(results.Hadiths), results.Total)

	h.sendMessageWithKeyboard(user, resultText, menu)
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
	sb.WriteString(hadith.Arabic)
	sb.WriteString("\n\n*English:*\n")
	sb.WriteString(hadith.English)
	sb.WriteString("\n\n")
	sb.WriteString("*Reference:* ")
	sb.WriteString(collectionName)
	sb.WriteString(", Book ")
	sb.WriteString(fmt.Sprintf("%d", bookNumber))
	sb.WriteString(", Hadith #")
	sb.WriteString(fmt.Sprintf("%d", hadith.HadithNumber))
	sb.WriteString("\n")
	sb.WriteString("*Grade:* ")
	sb.WriteString(grade)
	sb.WriteString("\n")

	return sb.String()
}

// sendMessage sends a text message
func (h *Handler) sendMessage(user *telebot.User, text string) {
	_, err := h.bot.Send(user, text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		h.log.LogError(err, "sendMessage")
	}
}

// sendMessageWithKeyboard sends a text message with keyboard
func (h *Handler) sendMessageWithKeyboard(user *telebot.User, text string, keyboard *telebot.ReplyMarkup) {
	_, err := h.bot.Send(user, text, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown, ReplyMarkup: keyboard})
	if err != nil {
		h.log.LogError(err, "sendMessageWithKeyboard")
	}
}

// escapeMarkdown escapes special characters for MarkdownV2
func (h *Handler) escapeMarkdown(text string) string {
	// Escape special MarkdownV2 characters
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
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
