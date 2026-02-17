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

type Handler struct {
	bot           *tgbotapi.BotAPI
	hadithService *services.HadithService
	log           *logger.Logger
	rateLimiter   *RateLimiter
}

func NewHandler(bot *tgbotapi.BotAPI, hadithService *services.HadithService, log *logger.Logger, rateLimitRequests int, rateLimitWindow time.Duration) *Handler {
	return &Handler{
		bot:           bot,
		hadithService: hadithService,
		log:           log,
		rateLimiter:   NewRateLimiter(rateLimitRequests, rateLimitWindow),
	}
}

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
	text := `*Welcome to Hadith Portal Bot* 🕌

Access authentic hadith collections from the six major books.

Use the keyboard below to navigate:`

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Browse Collections", "collections:1"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search Hadith", "search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Random Hadith", "random"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
		),
	)
	h.sendMessageWithKeyboard(m.Chat.ID, text, kb)
}

func (h *Handler) handleHelp(m *tgbotapi.Message) {
	helpText := `*Hadith Portal Bot - Help* ❓

/collections - Browse hadith collections
/search <keyword> - Search for hadiths
/random - Get a random hadith`
	h.sendMessage(m.Chat.ID, helpText)
}

func (h *Handler) handleRandom(m *tgbotapi.Message) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith == nil {
		h.sendMessage(m.Chat.ID, "Could not fetch hadith. Try again.")
		return
	}
	txt := h.formatHadithDisplay(res.Hadith, res.Collection, res.Book)
	h.sendMessage(m.Chat.ID, txt)
}

func (h *Handler) handleSearch(m *tgbotapi.Message) {
	args := m.CommandArguments()
	if args == "" {
		h.sendMessage(m.Chat.ID, "Please provide a keyword: /search prayer")
		return
	}
	results := h.hadithService.SearchHadiths(args, 1, 10)
	if len(results.Hadiths) == 0 {
		h.sendMessage(m.Chat.ID, "No results found.")
		return
	}
	h.sendSearchResults(m.Chat.ID, 0, args, results)
}

func (h *Handler) handleCollections(m *tgbotapi.Message) {
	h.sendCollectionsMenu(m.Chat.ID, 0, h.hadithService.GetCollections(), 1)
}

// --- CALLBACK HANDLER ---

func (h *Handler) handleCallback(c *tgbotapi.CallbackQuery) {
	if !h.rateLimiter.Allow(c.From.ID) {
		h.bot.Request(tgbotapi.NewCallback(c.ID, "Rate limit reached."))
		return
	}

	parts := strings.Split(c.Data, ":")
	chatID := c.Message.Chat.ID
	msgID := c.Message.MessageID

	switch parts[0] {
	case "collections":
		page, _ := strconv.Atoi(parts[1])
		h.sendCollectionsMenu(chatID, msgID, h.hadithService.GetCollections(), page)
	case "books":
		colName := parts[1]
		page, _ := strconv.Atoi(parts[2])
		h.sendBooksMenu(chatID, msgID, colName, h.hadithService.GetBooks(colName), page)
	case "hadiths":
		colName := parts[1]
		bookNum, _ := strconv.Atoi(parts[2])
		page, _ := strconv.Atoi(parts[3])
		res := h.hadithService.GetHadiths(colName, bookNum, page, 10)
		h.sendHadithsMenu(chatID, msgID, colName, bookNum, res)
	case "hadith_detail":
		h.handleHadithDetailCallback(c, parts)
	case "hadith_search":
		h.handleHadithSearchCallback(c, parts)
	case "random":
		h.handleRandomCallback(c)
	case "search_next", "search_prev":
		query := parts[1]
		page, _ := strconv.Atoi(parts[2])
		res := h.hadithService.SearchHadiths(query, page, 5)
		h.sendSearchResults(chatID, msgID, query, res)
	case "help":
		h.sendMessage(chatID, "Use /help for commands.")
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

// --- NAVIGATION MENUS ---

func (h *Handler) sendCollectionsMenu(chatID int64, msgID int, collections []models.Collection, page int) {
	const perPage = 6
	start, end := (page-1)*perPage, page*perPage
	if end > len(collections) {
		end = len(collections)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range collections[start:end] {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(c.Title, fmt.Sprintf("books:%s:1", c.Name))))
	}

	var nav []tgbotapi.InlineKeyboardButton
	if page > 1 {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("collections:%d", page-1)))
	}
	if end < len(collections) {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("collections:%d", page+1)))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	h.editOrSendMessage(chatID, msgID, "📚 *Select a Collection:*", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendBooksMenu(chatID int64, msgID int, col string, books []models.Book, page int) {
	const perPage = 10
	start, end := (page-1)*perPage, page*perPage
	if end > len(books) {
		end = len(books)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, b := range books[start:end] {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(truncate(b.Title, 35), fmt.Sprintf("hadiths:%s:%d:1", col, b.BookNumber))))
	}

	var nav []tgbotapi.InlineKeyboardButton
	if page > 1 {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Page", fmt.Sprintf("books:%s:%d", col, page-1)))
	}
	if end < len(books) {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Page ➡️", fmt.Sprintf("books:%s:%d", col, page+1)))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Collections", "collections:1")))
	h.editOrSendMessage(chatID, msgID, fmt.Sprintf("📚 *%s - Books:*", services.GetCollectionDisplayName(col)), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendHadithsMenu(chatID int64, msgID int, col string, bookNum int, result models.HadithResponse) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, hadith := range result.Hadiths {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🵿 Hadith #%d", hadith.HadithNumber), fmt.Sprintf("hadith_detail:%s:%d:%d:%d", col, bookNum, result.Page, i))))
	}

	var nav []tgbotapi.InlineKeyboardButton
	if result.Page > 1 {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("hadiths:%s:%d:%d", col, bookNum, result.Page-1)))
	}
	if result.Page < result.TotalPages {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("hadiths:%s:%d:%d", col, bookNum, result.Page+1)))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Books", fmt.Sprintf("books:%s:1", col))))
	h.editOrSendMessage(chatID, msgID, fmt.Sprintf("📑 *Page %d/%d*", result.Page, result.TotalPages), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) handleHadithDetailCallback(c *tgbotapi.CallbackQuery, parts []string) {
	col, bookNum, page, index := parts[1], 0, 0, 0
	fmt.Sscanf(parts[2], "%d", &bookNum)
	fmt.Sscanf(parts[3], "%d", &page)
	fmt.Sscanf(parts[4], "%d", &index)

	res := h.hadithService.GetHadiths(col, bookNum, page, 10)
	if index < len(res.Hadiths) {
		hadith := res.Hadiths[index]
		txt := h.formatHadithDisplay(&hadith, h.hadithService.GetCollection(col), h.hadithService.GetBook(col, bookNum))
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", fmt.Sprintf("hadiths:%s:%d:%d", col, bookNum, page))))
		h.editOrSendMessage(c.Message.Chat.ID, c.Message.MessageID, txt, kb)
	}
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
					tgbotapi.NewInlineKeyboardButtonData("🔍 New Search", "search"),
					tgbotapi.NewInlineKeyboardButtonData("🎲 Random", "random"),
				))
				h.editOrSendMessage(c.Message.Chat.ID, c.Message.MessageID, txt, kb)
				return
			}
		}
	}
}

func (h *Handler) handleRandomCallback(c *tgbotapi.CallbackQuery) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith != nil {
		txt := h.formatHadithDisplay(res.Hadith, res.Collection, res.Book)
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random")))
		h.editOrSendMessage(c.Message.Chat.ID, c.Message.MessageID, txt, kb)
	}
}

// --- FORMATTING & UTILS ---

func (h *Handler) editOrSendMessage(chatID int64, msgID int, text string, kb tgbotapi.InlineKeyboardMarkup) {
	if msgID != 0 {
		edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
		edit.ParseMode = tgbotapi.ModeMarkdown
		edit.ReplyMarkup = &kb
		h.bot.Send(edit)
	} else {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = kb
		h.bot.Send(msg)
	}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

func (h *Handler) sendMessageWithKeyboard(chatID int64, text string, kb tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = kb
	h.bot.Send(msg)
}

func (h *Handler) sendSearchResults(chatID int64, msgID int, query string, res models.SearchResult) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, hd := range res.Hadiths {
		col := h.findCollectionForHadith(hd)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Hadith #%d (%s)", hd.HadithNumber, col), fmt.Sprintf("hadith_search:%s:%d", col, hd.HadithNumber))))
	}

	if res.TotalPages > 1 {
		var nav []tgbotapi.InlineKeyboardButton
		if res.Page > 1 {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("search_prev:%s:%d", query, res.Page-1)))
		}
		if res.Page < res.TotalPages {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("search_next:%s:%d", query, res.Page+1)))
		}
		rows = append(rows, nav)
	}

	h.editOrSendMessage(chatID, msgID, fmt.Sprintf("🔍 *Results for:* %s", query), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) formatHadithDisplay(hdt *models.Hadith, col *models.Collection, b *models.Book) string {
	colTitle := "Unknown"
	if col != nil {
		colTitle = col.Title
	}
	bookNum := 0
	if b != nil {
		bookNum = b.BookNumber
	}
	grade := "Sahih"
	if hdt.Grade != "" {
		grade = hdt.Grade
	}

	return fmt.Sprintf("🵿 *Hadith*\n\n%s\n\n%s\n\n*Reference:* %s, Book %d, #%d\n*Grade:* %s",
		hdt.Arabic, hdt.English, colTitle, bookNum, hdt.HadithNumber, grade)
}

func (h *Handler) formatHadithForInlineHTML(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	colTitle := "Unknown"
	if collection != nil {
		colTitle = collection.Title
	}
	grade := "Sahih"
	if hadith.Grade != "" {
		grade = hadith.Grade
	}

	return fmt.Sprintf("🵿 <b>Hadith</b>\n\n%s\n\n%s\n\n<b>Ref:</b> %s, #%d\n<b>Grade:</b> %s",
		html.EscapeString(hadith.Arabic), html.EscapeString(hadith.English), html.EscapeString(colTitle), hadith.HadithNumber, html.EscapeString(grade))
}

func (h *Handler) findCollectionForHadith(hadith models.Hadith) string {
	for _, col := range h.hadithService.GetCollections() {
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
