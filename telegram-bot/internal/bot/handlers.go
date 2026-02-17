package bot

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"hadith-bot/internal/logger"
	"hadith-bot/internal/models"
	"hadith-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

// StartListening replaces HandleCommands to process the update stream
func (h *Handler) StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			h.handleIncomingMessage(update.Message)
		} else if update.CallbackQuery != nil {
			h.handleCallback(update.CallbackQuery)
		} else if update.InlineQuery != nil {
			h.handleInlineQuery(update.InlineQuery)
		}
	}
}

func (h *Handler) handleIncomingMessage(m *tgbotapi.Message) {
	if !h.rateLimiter.Allow(m.From.ID) {
		h.sendMessage(m.Chat.ID, "Please wait a moment before sending another command.")
		return
	}

	if m.IsCommand() {
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
}

// --- COMMAND HANDLERS ---

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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Browse Collections", "collections:1"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search Hadith", "search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Random Hadith", "random"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
		),
	)

	h.sendMessageWithKeyboard(m.Chat.ID, welcomeText, keyboard)
}

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
Example: /search prayer`

	h.sendMessage(m.Chat.ID, helpText)
}

func (h *Handler) handleRandom(m *tgbotapi.Message) {
	h.log.LogRequest(int64(m.MessageID), m.From.ID, "/random")
	result := h.hadithService.GetRandomHadith()

	if result.Hadith == nil {
		h.sendMessage(m.Chat.ID, "Sorry, couldn't fetch a random hadith. Please try again.")
		h.log.LogResponse(m.From.ID, "random", false)
		return
	}

	hadithText := h.formatHadithDisplay(result.Hadith, result.Collection, result.Book)
	h.sendMessage(m.Chat.ID, hadithText)
	h.log.LogResponse(m.From.ID, "random", true)
}

func (h *Handler) handleSearch(m *tgbotapi.Message) {
	args := m.CommandArguments()
	if args == "" {
		h.sendMessage(m.Chat.ID, "Please provide a search keyword.\nUsage: /search <keyword>")
		return
	}

	h.log.LogRequest(int64(m.MessageID), m.From.ID, "/search "+args)
	results := h.hadithService.SearchHadiths(args, 1, 10)

	if len(results.Hadiths) == 0 {
		h.sendMessage(m.Chat.ID, fmt.Sprintf("No results found for: *%s*", args))
		return
	}

	h.sendSearchResults(m.Chat.ID, args, results)
	h.log.LogResponse(m.From.ID, "search", true)
}

func (h *Handler) handleCollections(m *tgbotapi.Message) {
	h.log.LogRequest(int64(m.MessageID), m.From.ID, "/collections")
	collections := h.hadithService.GetCollections()
	h.sendCollectionsMenu(m.Chat.ID, collections, 1)
	h.log.LogResponse(m.From.ID, "collections", true)
}

// --- CALLBACK HANDLER ---

func (h *Handler) handleCallback(c *tgbotapi.CallbackQuery) {
	if !h.rateLimiter.Allow(c.From.ID) {
		callbackCfg := tgbotapi.NewCallback(c.ID, "Please wait a moment...")
		callbackCfg.ShowAlert = true
		h.bot.Request(callbackCfg)
		return
	}

	h.log.LogCallback(c.ID, c.From.ID, c.Data)
	parts := strings.Split(c.Data, ":")
	if len(parts) < 1 {
		return
	}

	// Navigation cleanup logic (Delete previous message to keep chat clean)
	isNavigation := false
	switch parts[0] {
	case "collections", "books", "hadiths", "search", "search_next", "search_prev", "help":
		isNavigation = true
	}

	if isNavigation && c.Message != nil {
		del := tgbotapi.NewDeleteMessage(c.Message.Chat.ID, c.Message.MessageID)
		h.bot.Request(del)
	}

	switch parts[0] {
	case "collections":
		page, _ := strconv.Atoi(parts[1])
		h.sendCollectionsMenu(c.Message.Chat.ID, h.hadithService.GetCollections(), page)
	case "books":
		collectionName := parts[1]
		page, _ := strconv.Atoi(parts[2])
		books := h.hadithService.GetBooks(collectionName)
		h.sendBooksMenu(c.Message.Chat.ID, collectionName, books, page)
	case "hadiths":
		collectionName := parts[1]
		bookNum, _ := strconv.Atoi(parts[2])
		page, _ := strconv.Atoi(parts[3])
		res := h.hadithService.GetHadiths(collectionName, bookNum, page, 10)
		h.sendHadithsMenu(c.Message.Chat.ID, collectionName, bookNum, res)
	case "hadith_detail":
		h.handleHadithDetailCallback(c, parts)
	case "hadith_search":
		h.handleHadithSearchCallback(c, parts)
	case "random":
		h.handleRandomCallback(c)
	case "search":
		h.sendMessage(c.Message.Chat.ID, "Please use /search <keyword> command to search.")
	case "search_next", "search_prev":
		query := parts[1]
		page, _ := strconv.Atoi(parts[2])
		res := h.hadithService.SearchHadiths(query, page, 5)
		h.sendSearchResults(c.Message.Chat.ID, query, res)
	case "help":
		h.handleHelpCallback(c)
	}

	h.bot.Request(tgbotapi.NewCallback(c.ID, ""))
}

// --- INLINE QUERY HANDLER ---

func (h *Handler) handleInlineQuery(q *tgbotapi.InlineQuery) {
	query := strings.TrimSpace(q.Query)
	var results []interface{}

	if query == "random" {
		res := h.hadithService.GetRandomHadith()
		if res.Hadith != nil {
			txt := h.formatHadithForInlineHTML(res.Hadith, res.Collection, res.Book)
			article := tgbotapi.NewInlineQueryResultArticleHTML(q.ID, "🎲 Random Hadith", txt)
			article.Description = fmt.Sprintf("Hadith #%d from %s", res.Hadith.HadithNumber, getCollectionTitle(res.Collection))
			results = append(results, article)
		}
	} else if strings.HasPrefix(query, "search ") {
		keyword := strings.TrimSpace(strings.TrimPrefix(query, "search "))
		if keyword != "" {
			searchRes := h.hadithService.SearchHadiths(keyword, 1, 5)
			for i, hadith := range searchRes.Hadiths {
				colName := h.findCollectionForHadith(hadith)
				col := h.hadithService.GetCollection(colName)
				txt := h.formatHadithForInlineHTML(&hadith, col, nil)

				id := fmt.Sprintf("inline_%d_%d", hadith.HadithNumber, i)
				article := tgbotapi.NewInlineQueryResultArticleHTML(id, fmt.Sprintf("🵿 Hadith #%d", hadith.HadithNumber), txt)
				article.Description = truncate(hadith.English, 50)
				results = append(results, article)
			}
		}
	}

	h.bot.Request(tgbotapi.InlineConfig{
		InlineQueryID: q.ID,
		Results:       results,
		CacheTime:     10,
	})
}

// --- UI HELPERS ---

func (h *Handler) sendCollectionsMenu(chatID int64, collections []models.Collection, page int) {
	const perPage = 6
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(collections) {
		end = len(collections)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range collections[start:end] {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(c.Title, fmt.Sprintf("books:%s:1", c.Name)),
		))
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("⬅️ Previous", fmt.Sprintf("collections:%d", page-1)))
	}
	if end < len(collections) {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("collections:%d", page+1)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	h.sendMessageWithKeyboard(chatID, "📚 *Hadith Collections*\n\nSelect a collection to browse:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendBooksMenu(chatID int64, collection string, books []models.Book, page int) {
	const perPage = 10
	start, end := (page-1)*perPage, page*perPage
	if end > len(books) {
		end = len(books)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, b := range books[start:end] {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("📖 %s", truncate(b.Title, 40)), fmt.Sprintf("hadiths:%s:%d:1", collection, b.BookNumber)),
		))
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("books:%s:%d", collection, page-1)))
	}
	if end < len(books) {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("books:%s:%d", collection, page+1)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Collections", "collections:1")))

	text := fmt.Sprintf("📚 *%s*\n\nSelect a book:", services.GetCollectionDisplayName(collection))
	h.sendMessageWithKeyboard(chatID, text, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendHadithsMenu(chatID int64, collection string, bookNumber int, result models.HadithResponse) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, hadith := range result.Hadiths {
		grade := "Sahih"
		if hadith.Grade != "" {
			grade = hadith.Grade
		}
		data := fmt.Sprintf("hadith_detail:%s:%d:%d:%d", collection, bookNumber, result.Page, i)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🵿 Hadith #%d [%s]", hadith.HadithNumber, grade), data),
		))
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if result.Page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page-1)))
	}
	if result.Page < result.TotalPages {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("hadiths:%s:%d:%d", collection, bookNumber, result.Page+1)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Books", fmt.Sprintf("books:%s:1", collection))))

	text := fmt.Sprintf("📖 Page %d/%d - Showing %d hadiths:", result.Page, result.TotalPages, len(result.Hadiths))
	h.sendMessageWithKeyboard(chatID, text, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendSearchResults(chatID int64, query string, results models.SearchResult) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, hadith := range results.Hadiths {
		col := h.findCollectionForHadith(hadith)
		data := fmt.Sprintf("hadith_search:%s:%d", col, hadith.HadithNumber)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Hadith #%d", hadith.HadithNumber), data),
		))
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if results.Page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("search_prev:%s:%d", query, results.Page-1)))
	}
	if results.Page < results.TotalPages {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("search_next:%s:%d", query, results.Page+1)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	text := fmt.Sprintf("🔍 *Search Results for:* %s\nPage %d/%d", query, results.Page, results.TotalPages)
	h.sendMessageWithKeyboard(chatID, text, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

// --- CORE LOGIC HANDLERS ---

func (h *Handler) handleHadithDetailCallback(c *tgbotapi.CallbackQuery, parts []string) {
	collectionName, bookNum, page, index := parts[1], 0, 0, 0
	fmt.Sscanf(parts[2], "%d", &bookNum)
	fmt.Sscanf(parts[3], "%d", &page)
	fmt.Sscanf(parts[4], "%d", &index)

	result := h.hadithService.GetHadiths(collectionName, bookNum, page, 10)
	if index >= len(result.Hadiths) {
		return
	}

	hadith := result.Hadiths[index]
	txt := h.formatHadithDisplay(&hadith, h.hadithService.GetCollection(collectionName), h.hadithService.GetBook(collectionName, bookNum))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Random", "random"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "search"),
		),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", fmt.Sprintf("hadiths:%s:%d:%d", collectionName, bookNum, page))),
	)
	h.sendMessageWithKeyboard(c.Message.Chat.ID, txt, kb)
}

func (h *Handler) handleHadithSearchCallback(c *tgbotapi.CallbackQuery, parts []string) {
	colName := parts[1]
	hadithNum, _ := strconv.Atoi(parts[2])

	books := h.hadithService.GetBooks(colName)
	for _, b := range books {
		res := h.hadithService.GetHadiths(colName, b.BookNumber, 1, 1000)
		for _, hadith := range res.Hadiths {
			if hadith.HadithNumber == hadithNum {
				txt := h.formatHadithDisplay(&hadith, h.hadithService.GetCollection(colName), &b)
				kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🎲 Random", "random"),
					tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "search"),
				))
				h.sendMessageWithKeyboard(c.Message.Chat.ID, txt, kb)
				return
			}
		}
	}
}

func (h *Handler) handleRandomCallback(c *tgbotapi.CallbackQuery) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith == nil {
		return
	}
	txt := h.formatHadithDisplay(res.Hadith, res.Collection, res.Book)
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random"),
		tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "search"),
	))
	h.sendMessageWithKeyboard(c.Message.Chat.ID, txt, kb)
}

func (h *Handler) handleHelpCallback(c *tgbotapi.CallbackQuery) {
	h.handleHelp(c.Message)
}

// --- FORMATTING & UTILS ---

func (h *Handler) formatHadithDisplay(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	colTitle := "Unknown Collection"
	if collection != nil {
		colTitle = collection.Title
	}
	bookNum := 0
	if book != nil {
		bookNum = book.BookNumber
	}
	grade := "Sahih"
	if hadith.Grade != "" {
		grade = hadith.Grade
	}

	return fmt.Sprintf("🵿 *Hadith*\n\n*Arabic:*\n%s\n\n*English:*\n%s\n\n*Reference:* %s, Book %d, Hadith #%d\n*Grade:* %s",
		hadith.Arabic, hadith.English, colTitle, bookNum, hadith.HadithNumber, grade)
}

func (h *Handler) formatHadithForInlineHTML(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	colTitle := "Unknown Collection"
	if collection != nil {
		colTitle = collection.Title
	}
	bookNum := 0
	if book != nil {
		bookNum = book.BookNumber
	}
	grade := "Sahih"
	if hadith.Grade != "" {
		grade = hadith.Grade
	}

	return fmt.Sprintf("🵿 <b>Hadith</b>\n\n<b>Arabic:</b>\n%s\n\n<b>English:</b>\n%s\n\n<b>Reference:</b> %s, Book %d, Hadith #%d\n<b>Grade:</b> %s",
		html.EscapeString(hadith.Arabic), html.EscapeString(hadith.English), html.EscapeString(colTitle), bookNum, hadith.HadithNumber, html.EscapeString(grade))
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

func (h *Handler) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *Handler) findCollectionForHadith(hadith models.Hadith) string {
	collections := h.hadithService.GetCollections()
	for _, col := range collections {
		books := h.hadithService.GetBooks(col.Name)
		for _, b := range books {
			if b.BookNumber == hadith.ChapterID {
				return col.Name
			}
		}
	}
	return "bukhari"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getCollectionTitle(c *models.Collection) string {
	if c == nil {
		return "Unknown"
	}
	return c.Title
}
