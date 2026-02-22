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
		h.sendMessage(m.Chat.ID, "⚠️ Please wait a moment before sending another command.")
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

Access authentic hadith collections from the six major books. Use the keyboard below to navigate:`

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

/collections - Browse major hadith books
/search <keyword> - Search for specific topics
/random - Get an inspired daily hadith

_Tip: You can also use this bot in any chat by typing @YourBotName followed by a keyword._`
	h.sendMessage(m.Chat.ID, helpText)
}

func (h *Handler) handleRandom(m *tgbotapi.Message) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith == nil {
		h.sendMessage(m.Chat.ID, "Could not fetch hadith. Try again.")
		return
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random"),
		),
	)
	h.sendHadithResponse(m.Chat.ID, res.Hadith, res.Collection, res.Book, &kb)
}

func (h *Handler) handleSearch(m *tgbotapi.Message) {
	args := m.CommandArguments()
	if args == "" {
		h.sendMessage(m.Chat.ID, "Please provide a keyword: `/search fasting`")
		return
	}
	// Limit query length for safety
	if len(args) > 100 {
		args = args[:100]
	}
	results := h.hadithService.SearchHadiths(args, 1, 5)
	if len(results.Hadiths) == 0 {
		h.sendMessage(m.Chat.ID, "No results found for that keyword.")
		return
	}
	h.sendSearchResults(m.Chat.ID, 0, "", args, results)
}

func (h *Handler) handleCollections(m *tgbotapi.Message) {
	h.sendCollectionsMenu(m.Chat.ID, 0, "", h.hadithService.GetCollections(), 1)
}

// --- CALLBACK HANDLER ---

func (h *Handler) handleCallback(c *tgbotapi.CallbackQuery) {
	if !h.rateLimiter.Allow(c.From.ID) {
		h.bot.Request(tgbotapi.NewCallback(c.ID, "Rate limit reached."))
		return
	}

	parts := strings.Split(c.Data, ":")
	var chatID int64
	var msgID int
	if c.Message != nil {
		chatID = c.Message.Chat.ID
		msgID = c.Message.MessageID
	}
	iMID := c.InlineMessageID

	switch parts[0] {
	case "collections":
		page, _ := strconv.Atoi(parts[1])
		h.sendCollectionsMenu(chatID, msgID, iMID, h.hadithService.GetCollections(), page)
	case "books":
		colName := parts[1]
		page, _ := strconv.Atoi(parts[2])
		h.sendBooksMenu(chatID, msgID, iMID, colName, h.hadithService.GetBooks(colName), page)
	case "hadiths":
		colName := parts[1]
		bookNum, _ := strconv.Atoi(parts[2])
		page, _ := strconv.Atoi(parts[3])
		res := h.hadithService.GetHadiths(colName, bookNum, page, 10)
		h.sendHadithsMenu(chatID, msgID, iMID, colName, bookNum, res)
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
		h.sendSearchResults(chatID, msgID, iMID, query, res)
	case "help":
		h.handleHelp(c.Message)
	}

	h.bot.Request(tgbotapi.NewCallback(c.ID, ""))
}

// --- CORE ROBUST SENDING LOGIC ---

func (h *Handler) sendHadithResponse(chatID int64, hdt *models.Hadith, col *models.Collection, b *models.Book, kb *tgbotapi.InlineKeyboardMarkup) {
	blocks := h.formatHadithBlocks(hdt, col, b)

	for i, block := range blocks {
		// Safety split for unusually long blocks (> 4000 chars)
		subChunks := splitText(block, 4000)

		for j, chunk := range subChunks {
			msg := tgbotapi.NewMessage(chatID, chunk)
			msg.ParseMode = tgbotapi.ModeMarkdown

			// Attach the keyboard ONLY to the very last chunk of the very last block
			if i == len(blocks)-1 && j == len(subChunks)-1 && kb != nil {
				msg.ReplyMarkup = kb
			}

			h.bot.Send(msg)
		}
	}
}

func (h *Handler) formatHadithBlocks(hdt *models.Hadith, col *models.Collection, b *models.Book) []string {
	var blocks []string

	// Block 1: Arabic
	blocks = append(blocks, fmt.Sprintf("🕌 *Arabic*\n\n%s", hdt.Arabic))

	// Block 2: Narrator & English
	narrator := ""
	if hdt.Narrator != "" {
		narrator = fmt.Sprintf("*Narrated by:* %s\n\n", hdt.Narrator)
	}
	blocks = append(blocks, fmt.Sprintf("📖 *Translation*\n\n%s%s", narrator, hdt.English))

	// Block 3: Metadata
	colTitle := "Unknown Collection"
	if col != nil {
		colTitle = col.Title
	}
	bookTitle := "General"
	if b != nil {
		bookTitle = b.Title
	}
	grade := "Sahih"
	if hdt.Grade != "" {
		grade = hdt.Grade
	}
	meta := fmt.Sprintf("📍 *Reference*\n*Collection:* %s\n*Book:* %s\n*Hadith:* #%d\n*Grade:* %s",
		colTitle, bookTitle, hdt.HadithNumber, grade)
	blocks = append(blocks, meta)

	return blocks
}

func splitText(text string, limit int) []string {
	if len(text) <= limit {
		return []string{text}
	}
	var chunks []string
	for len(text) > limit {
		splitIdx := strings.LastIndex(text[:limit], "\n")
		if splitIdx == -1 {
			splitIdx = strings.LastIndex(text[:limit], " ")
		}
		if splitIdx == -1 {
			splitIdx = limit
		}
		chunks = append(chunks, strings.TrimSpace(text[:splitIdx]))
		text = text[splitIdx:]
	}
	if len(text) > 0 {
		chunks = append(chunks, strings.TrimSpace(text))
	}
	return chunks
}

// --- NAVIGATION & CALLBACK HELPERS ---

func (h *Handler) handleHadithDetailCallback(c *tgbotapi.CallbackQuery, parts []string) {
	col, bookNum, page, index := parts[1], 0, 0, 0
	fmt.Sscanf(parts[2], "%d", &bookNum)
	fmt.Sscanf(parts[3], "%d", &page)
	fmt.Sscanf(parts[4], "%d", &index)

	res := h.hadithService.GetHadiths(col, bookNum, page, 10)
	if index < len(res.Hadiths) {
		// Delete old menu to keep chat clean
		if c.Message != nil {
			h.bot.Send(tgbotapi.NewDeleteMessage(c.Message.Chat.ID, c.Message.MessageID))
		}

		hadith := res.Hadiths[index]
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to List", fmt.Sprintf("hadiths:%s:%d:%d", col, bookNum, page)),
		))
		h.sendHadithResponse(c.Message.Chat.ID, &hadith, h.hadithService.GetCollection(col), h.hadithService.GetBook(col, bookNum), &kb)
	}
}

func (h *Handler) handleHadithSearchCallback(c *tgbotapi.CallbackQuery, parts []string) {
	colName := parts[1]
	hadithNum, _ := strconv.Atoi(parts[2])

	// Assuming your service has a way to fetch a single hadith + its book context
	hadith, book := h.hadithService.FindHadithByNumber(colName, hadithNum)

	if hadith != nil {
		if c.Message != nil {
			h.bot.Send(tgbotapi.NewDeleteMessage(c.Message.Chat.ID, c.Message.MessageID))
		}

		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 New Search", "search"),
			tgbotapi.NewInlineKeyboardButtonData("🎲 Random", "random"),
		))
		h.sendHadithResponse(c.Message.Chat.ID, hadith, h.hadithService.GetCollection(colName), book, &kb)
	}
}

func (h *Handler) handleRandomCallback(c *tgbotapi.CallbackQuery) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith != nil {
		if c.Message != nil {
			h.bot.Send(tgbotapi.NewDeleteMessage(c.Message.Chat.ID, c.Message.MessageID))
		}
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random"),
		))
		h.sendHadithResponse(c.Message.Chat.ID, res.Hadith, res.Collection, res.Book, &kb)
	}
}

func (h *Handler) sendSearchResults(chatID int64, msgID int, inlineMsgID string, query string, res models.SearchResult) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, hd := range res.Hadiths {
		col := h.findCollectionForHadith(hd)
		// Preview with snippet
		preview := truncate(fmt.Sprintf("%s #%d: %s", strings.Title(col), hd.HadithNumber, hd.English), 45)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(preview, fmt.Sprintf("hadith_search:%s:%d", col, hd.HadithNumber)),
		))
	}

	var nav []tgbotapi.InlineKeyboardButton
	if res.Page > 1 {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("search_prev:%s:%d", query, res.Page-1)))
	}
	nav = append(nav, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Page %d/%d", res.Page, res.TotalPages), "ignore"))
	if res.Page < res.TotalPages {
		nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next ➡️", fmt.Sprintf("search_next:%s:%d", query, res.Page+1)))
	}
	rows = append(rows, nav)

	text := fmt.Sprintf("🔍 *Results for:* `%s`", query)
	h.editOrSendMessage(chatID, msgID, inlineMsgID, text, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

// --- MENUS ---

func (h *Handler) sendCollectionsMenu(chatID int64, msgID int, inlineMsgID string, collections []models.Collection, page int) {
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

	h.editOrSendMessage(chatID, msgID, inlineMsgID, "📚 *Select a Collection:*", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendBooksMenu(chatID int64, msgID int, inlineMsgID string, col string, books []models.Book, page int) {
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
	h.editOrSendMessage(chatID, msgID, inlineMsgID, fmt.Sprintf("📚 *%s - Books:*", services.GetCollectionDisplayName(col)), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendHadithsMenu(chatID int64, msgID int, inlineMsgID string, col string, bookNum int, result models.HadithResponse) {
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
	h.editOrSendMessage(chatID, msgID, inlineMsgID, fmt.Sprintf("📑 *Page %d/%d*", result.Page, result.TotalPages), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

// --- UTILS ---

func (h *Handler) editOrSendMessage(chatID int64, msgID int, inlineMsgID string, text string, kb tgbotapi.InlineKeyboardMarkup) {
	if inlineMsgID != "" {
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit:        tgbotapi.BaseEdit{InlineMessageID: inlineMsgID, ReplyMarkup: &kb},
			Text:            text,
			ParseMode:       tgbotapi.ModeMarkdown,
		}
		h.bot.Send(edit)
	} else if msgID != 0 {
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
	chunks := splitText(text, 4000)
	for _, chunk := range chunks {
		msg := tgbotapi.NewMessage(chatID, chunk)
		msg.ParseMode = tgbotapi.ModeMarkdown
		h.bot.Send(msg)
	}
}

func (h *Handler) sendMessageWithKeyboard(chatID int64, text string, kb tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = kb
	h.bot.Send(msg)
}

func (h *Handler) handleInlineQuery(q *tgbotapi.InlineQuery) {
	query := strings.TrimSpace(q.Query)
	var results []interface{}

	if query == "" || query == "collections" {
		article := tgbotapi.NewInlineQueryResultArticleMarkdown(q.ID+"_br", "📚 Browse Collections", "Tap below to browse.")
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 Open", "collections:1"),
		))
		article.ReplyMarkup = &kb
		results = append(results, article)
	}

	if strings.HasPrefix(query, "search ") {
		kw := strings.TrimSpace(strings.TrimPrefix(query, "search "))
		if kw != "" {
			searchRes := h.hadithService.SearchHadiths(kw, 1, 5)
			for i, hd := range searchRes.Hadiths {
				txt := h.formatHadithForInlineHTML(&hd, nil, nil)
				if len(txt) > 4000 {
					txt = txt[:3990] + "..."
				}
				article := tgbotapi.NewInlineQueryResultArticleHTML(fmt.Sprintf("%d_%d", hd.HadithNumber, i), fmt.Sprintf("Hadith #%d", hd.HadithNumber), txt)
				article.Description = truncate(hd.English, 60)
				results = append(results, article)
			}
		}
	}

	h.bot.Request(tgbotapi.InlineConfig{InlineQueryID: q.ID, Results: results, CacheTime: 10})
}

func (h *Handler) formatHadithForInlineHTML(hadith *models.Hadith, collection *models.Collection, book *models.Book) string {
	narratorHTML := ""
	if hadith.Narrator != "" {
		narratorHTML = fmt.Sprintf("\n\n<b>Narrator:</b> %s", html.EscapeString(hadith.Narrator))
	}
	return fmt.Sprintf("🕌 <b>Arabic:</b>\n%s%s\n\n📖 <b>English:</b>\n%s",
		html.EscapeString(hadith.Arabic), narratorHTML, html.EscapeString(hadith.English))
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
