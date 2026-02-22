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

const telegramMessageMaxRunes = 3800

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
		h.sendMessage(m.Chat.ID, "⏳ Please wait a moment before sending another command.")
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
	text := `🕌 *Welcome to Hadith Portal Bot*

Explore authentic hadith from major collections with a clean, simple menu.

✨ *Quick actions:* use the buttons below to browse, search, or get a random hadith.`

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
	helpText := `❓ *Hadith Portal Bot — Help*

*Commands*
• */start* — Open the main menu
• */collections* — Browse hadith collections
• */search <keyword>* — Search hadith text
• */random* — Get a random hadith
• */help* — Show this help message

💡 *Examples*
• */search prayer*
• */search patience*`
	h.sendMessage(m.Chat.ID, helpText)
}

func (h *Handler) handleRandom(m *tgbotapi.Message) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith == nil || res.Collection == nil {
		h.sendMessage(m.Chat.ID, "⚠️ Could not fetch a hadith right now. Please try again.")
		return
	}
	h.sendRandomHadithPaged(m.Chat.ID, 0, "", res.Collection.Name, res.Hadith.HadithNumber, 0)
}

func (h *Handler) handleSearch(m *tgbotapi.Message) {
	args := m.CommandArguments()
	if args == "" {
		h.sendMessage(m.Chat.ID, "🔎 Please provide a keyword. Example: */search prayer*")
		return
	}
	results := h.hadithService.SearchHadiths(args, 1, 10)
	if len(results.Hadiths) == 0 {
		h.sendMessage(m.Chat.ID, "No results found for your search. Try a different keyword.")
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
		h.bot.Request(tgbotapi.NewCallback(c.ID, "⏳ Slow down a little."))
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
	case "hadith_page":
		h.handleHadithPageCallback(c, parts)
	case "random":
		h.handleRandomCallback(c)
	case "search_next", "search_prev":
		query := parts[1]
		page, _ := strconv.Atoi(parts[2])
		res := h.hadithService.SearchHadiths(query, page, 5)
		h.sendSearchResults(chatID, msgID, iMID, query, res)
	case "help":
		h.sendMessage(chatID, "Use */help* to view all commands and examples.")
	}

	h.bot.Request(tgbotapi.NewCallback(c.ID, ""))
}

// --- INLINE QUERY HANDLER ---

func (h *Handler) handleInlineQuery(q *tgbotapi.InlineQuery) {
	query := strings.TrimSpace(q.Query)
	var results []interface{}

	if query == "" || query == "collections" {
		text := "📚 *Browse Hadith Collections*\nSelect a collection to view its books and hadiths."
		article := tgbotapi.NewInlineQueryResultArticleMarkdown(q.ID+"_browse", "📚 Browse Collections", text)
		article.Description = "View Bukhari, Muslim, and other major collections"
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📚 Open Collections", "collections:1"),
				tgbotapi.NewInlineKeyboardButtonData("🎲 Random Hadith", "random"),
			),
		)
		article.ReplyMarkup = &kb
		results = append(results, article)
	}

	if query == "random" {
		res := h.hadithService.GetRandomHadith()
		if res.Hadith != nil && res.Collection != nil {
			txt := h.formatHadithDisplay(res.Hadith, res.Collection, res.Book)
			pages := splitTelegramMessage(txt, telegramMessageMaxRunes)
			if len(pages) == 0 {
				pages = []string{txt}
			}

			display := pages[0]
			if len(pages) > 1 {
				display = fmt.Sprintf("*Page 1/%d*\n\n%s", len(pages), display)
			}

			article := tgbotapi.NewInlineQueryResultArticleHTML(q.ID, "🎲 Random Hadith", txt)
			article.Description = fmt.Sprintf("Hadith #%d from %s", res.Hadith.HadithNumber, getCollectionTitle(res.Collection))
			article.InputMessageContent = tgbotapi.InputTextMessageContent{
				Text:      display,
				ParseMode: tgbotapi.ModeMarkdown,
			}

			var rows [][]tgbotapi.InlineKeyboardButton
			if len(pages) > 1 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Next Part ➡️", fmt.Sprintf("hadith_page:r:%s:%d:%d", res.Collection.Name, res.Hadith.HadithNumber, 1)),
				))
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random"),
			))
			kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
			article.ReplyMarkup = &kb
			results = append(results, article)
		}
	} else if strings.HasPrefix(query, "search ") {
		keyword := strings.TrimSpace(strings.TrimPrefix(query, "search "))
		if keyword != "" {
			searchRes := h.hadithService.SearchHadiths(keyword, 1, 5)
			for i, hadith := range searchRes.Hadiths {
				colName := h.findCollectionForHadith(hadith)
				col := h.hadithService.GetCollection(colName)
				txt := h.formatHadithDisplay(&hadith, col, nil)
				pages := splitTelegramMessage(txt, telegramMessageMaxRunes)
				if len(pages) == 0 {
					pages = []string{txt}
				}

				display := pages[0]
				if len(pages) > 1 {
					display = fmt.Sprintf("*Page 1/%d*\n\n%s", len(pages), display)
				}

				id := fmt.Sprintf("inline_%d_%d", hadith.HadithNumber, i)
				article := tgbotapi.NewInlineQueryResultArticleHTML(id, fmt.Sprintf("🵿 Hadith #%d", hadith.HadithNumber), txt)
				article.Description = truncate(hadith.English, 50)
				article.InputMessageContent = tgbotapi.InputTextMessageContent{
					Text:      display,
					ParseMode: tgbotapi.ModeMarkdown,
				}

				var rows [][]tgbotapi.InlineKeyboardButton
				if len(pages) > 1 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Next Part ➡️", fmt.Sprintf("hadith_page:s:%s:%d:%d", colName, hadith.HadithNumber, 1)),
					))
				}
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Open in Bot", fmt.Sprintf("hadith_search:%s:%d", colName, hadith.HadithNumber)),
				))
				kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
				article.ReplyMarkup = &kb
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
	h.editOrSendMessage(chatID, msgID, inlineMsgID, fmt.Sprintf("📚 *%s — Books*", services.GetCollectionDisplayName(col)), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) sendHadithsMenu(chatID int64, msgID int, inlineMsgID string, col string, bookNum int, result models.HadithResponse) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, hadith := range result.Hadiths {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("📜 Hadith #%d", hadith.HadithNumber), fmt.Sprintf("hadith_detail:%s:%d:%d:%d", col, bookNum, result.Page, i))))
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
	h.editOrSendMessage(chatID, msgID, inlineMsgID, fmt.Sprintf("📑 *Hadith List — Page %d/%d*", result.Page, result.TotalPages), tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) handleHadithDetailCallback(c *tgbotapi.CallbackQuery, parts []string) {
	col, bookNum, page, index := parts[1], 0, 0, 0
	fmt.Sscanf(parts[2], "%d", &bookNum)
	fmt.Sscanf(parts[3], "%d", &page)
	fmt.Sscanf(parts[4], "%d", &index)

	chatID := int64(0)
	msgID := 0
	if c.Message != nil {
		chatID = c.Message.Chat.ID
		msgID = c.Message.MessageID
	}
	h.sendHadithDetailPaged(chatID, msgID, c.InlineMessageID, col, bookNum, page, index, 0)
}

func (h *Handler) handleHadithPageCallback(c *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	chatID := int64(0)
	msgID := 0
	if c.Message != nil {
		chatID = c.Message.Chat.ID
		msgID = c.Message.MessageID
	}

	switch parts[1] {
	case "d":
		if len(parts) < 7 {
			return
		}
		col := parts[2]
		bookNum, _ := strconv.Atoi(parts[3])
		listPage, _ := strconv.Atoi(parts[4])
		index, _ := strconv.Atoi(parts[5])
		textPage, _ := strconv.Atoi(parts[6])
		h.sendHadithDetailPaged(chatID, msgID, c.InlineMessageID, col, bookNum, listPage, index, textPage)
	case "s":
		if len(parts) < 5 {
			return
		}
		col := parts[2]
		hadithNum, _ := strconv.Atoi(parts[3])
		textPage, _ := strconv.Atoi(parts[4])
		h.sendSearchHadithPaged(chatID, msgID, c.InlineMessageID, col, hadithNum, textPage)
	case "r":
		if len(parts) < 5 {
			return
		}
		col := parts[2]
		hadithNum, _ := strconv.Atoi(parts[3])
		textPage, _ := strconv.Atoi(parts[4])
		h.sendRandomHadithPaged(chatID, msgID, c.InlineMessageID, col, hadithNum, textPage)
	}
}

func (h *Handler) sendHadithDetailPaged(chatID int64, msgID int, inlineMsgID, col string, bookNum, page, index, textPage int) {

	res := h.hadithService.GetHadiths(col, bookNum, page, 10)
	if index < len(res.Hadiths) {
		hadith := res.Hadiths[index]
		txt := h.formatHadithDisplay(&hadith, h.hadithService.GetCollection(col), h.hadithService.GetBook(col, bookNum))
		pages := splitTelegramMessage(txt, telegramMessageMaxRunes)
		if len(pages) == 0 {
			pages = []string{txt}
		}

		if textPage < 0 {
			textPage = 0
		}
		if textPage >= len(pages) {
			textPage = len(pages) - 1
		}

		display := pages[textPage]
		if len(pages) > 1 {
			display = fmt.Sprintf("*Page %d/%d*\n\n%s", textPage+1, len(pages), display)
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		if len(pages) > 1 {
			var nav []tgbotapi.InlineKeyboardButton
			if textPage > 0 {
				nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev Part", fmt.Sprintf("hadith_page:d:%s:%d:%d:%d:%d", col, bookNum, page, index, textPage-1)))
			}
			if textPage < len(pages)-1 {
				nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next Part ➡️", fmt.Sprintf("hadith_page:d:%s:%d:%d:%d:%d", col, bookNum, page, index, textPage+1)))
			}
			if len(nav) > 0 {
				rows = append(rows, nav)
			}
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", fmt.Sprintf("hadiths:%s:%d:%d", col, bookNum, page))))
		h.editOrSendMessage(chatID, msgID, inlineMsgID, display, tgbotapi.NewInlineKeyboardMarkup(rows...))
	}
}

func (h *Handler) handleHadithSearchCallback(c *tgbotapi.CallbackQuery, parts []string) {
	colName := parts[1]
	hadithNum, _ := strconv.Atoi(parts[2])

	chatID := int64(0)
	msgID := 0
	if c.Message != nil {
		chatID = c.Message.Chat.ID
		msgID = c.Message.MessageID
	}
	h.sendSearchHadithPaged(chatID, msgID, c.InlineMessageID, colName, hadithNum, 0)
}

func (h *Handler) sendSearchHadithPaged(chatID int64, msgID int, inlineMsgID, colName string, hadithNum, textPage int) {
	hadith, book := h.hadithService.FindHadithByNumber(colName, hadithNum)
	if hadith == nil {
		var rows [][]tgbotapi.InlineKeyboardButton
		if inlineMsgID == "" {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔍 New Search", "search"),
			))
		}

		h.editOrSendMessage(chatID, msgID, inlineMsgID, "⚠️ Could not find this hadith. Please try searching again.", tgbotapi.NewInlineKeyboardMarkup(rows...))
		return
	}

	txt := h.formatHadithDisplay(hadith, h.hadithService.GetCollection(colName), book)
	pages := splitTelegramMessage(txt, telegramMessageMaxRunes)
	if len(pages) == 0 {
		pages = []string{txt}
	}

	if textPage < 0 {
		textPage = 0
	}
	if textPage >= len(pages) {
		textPage = len(pages) - 1
	}

	display := pages[textPage]
	if len(pages) > 1 {
		display = fmt.Sprintf("*Page %d/%d*\n\n%s", textPage+1, len(pages), display)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	if len(pages) > 1 {
		var nav []tgbotapi.InlineKeyboardButton
		if textPage > 0 {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev Part", fmt.Sprintf("hadith_page:s:%s:%d:%d", colName, hadithNum, textPage-1)))
		}
		if textPage < len(pages)-1 {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next Part ➡️", fmt.Sprintf("hadith_page:s:%s:%d:%d", colName, hadithNum, textPage+1)))
		}
		if len(nav) > 0 {
			rows = append(rows, nav)
		}
	}

	if inlineMsgID == "" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 New Search", "search"),
			tgbotapi.NewInlineKeyboardButtonData("🎲 Random", "random"),
		))
	}

	h.editOrSendMessage(chatID, msgID, inlineMsgID, display, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (h *Handler) handleRandomCallback(c *tgbotapi.CallbackQuery) {
	res := h.hadithService.GetRandomHadith()
	if res.Hadith != nil && res.Collection != nil {
		chatID := int64(0)
		msgID := 0
		if c.Message != nil {
			chatID = c.Message.Chat.ID
			msgID = c.Message.MessageID
		}
		h.sendRandomHadithPaged(chatID, msgID, c.InlineMessageID, res.Collection.Name, res.Hadith.HadithNumber, 0)
	}
}

func (h *Handler) sendRandomHadithPaged(chatID int64, msgID int, inlineMsgID, colName string, hadithNum, textPage int) {
	hadith, book := h.hadithService.FindHadithByNumber(colName, hadithNum)
	if hadith == nil {
		h.sendMessage(chatID, "⚠️ Could not fetch a hadith right now. Please try again.")
		return
	}

	txt := h.formatHadithDisplay(hadith, h.hadithService.GetCollection(colName), book)
	pages := splitTelegramMessage(txt, telegramMessageMaxRunes)
	if len(pages) == 0 {
		pages = []string{txt}
	}

	if textPage < 0 {
		textPage = 0
	}
	if textPage >= len(pages) {
		textPage = len(pages) - 1
	}

	display := pages[textPage]
	if len(pages) > 1 {
		display = fmt.Sprintf("*Page %d/%d*\n\n%s", textPage+1, len(pages), display)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	if len(pages) > 1 {
		var nav []tgbotapi.InlineKeyboardButton
		if textPage > 0 {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev Part", fmt.Sprintf("hadith_page:r:%s:%d:%d", colName, hadithNum, textPage-1)))
		}
		if textPage < len(pages)-1 {
			nav = append(nav, tgbotapi.NewInlineKeyboardButtonData("Next Part ➡️", fmt.Sprintf("hadith_page:r:%s:%d:%d", colName, hadithNum, textPage+1)))
		}
		if len(nav) > 0 {
			rows = append(rows, nav)
		}
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🎲 Another Random", "random")))
	h.editOrSendMessage(chatID, msgID, inlineMsgID, display, tgbotapi.NewInlineKeyboardMarkup(rows...))
}

// --- FORMATTING & UTILS ---

func (h *Handler) editOrSendMessage(chatID int64, msgID int, inlineMsgID string, text string, kb tgbotapi.InlineKeyboardMarkup) {
	if inlineMsgID != "" {
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				InlineMessageID: inlineMsgID,
				ReplyMarkup:     &kb,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeMarkdown,
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

func (h *Handler) sendSearchResults(chatID int64, msgID int, inlineMsgID string, query string, res models.SearchResult) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, hd := range res.Hadiths {
		col := h.findCollectionForHadith(hd)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("📜 Hadith #%d (%s)", hd.HadithNumber, services.GetCollectionDisplayName(col)), fmt.Sprintf("hadith_search:%s:%d", col, hd.HadithNumber))))
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

	h.editOrSendMessage(chatID, msgID, inlineMsgID, fmt.Sprintf("🔍 *Results for:* %s", escapeMarkdown(query)), tgbotapi.NewInlineKeyboardMarkup(rows...))
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

	narrator := ""
	if hdt.Narrator != "" {
		narrator = fmt.Sprintf("\n*:* %s\n", escapeMarkdown(hdt.Narrator))
	}

	return fmt.Sprintf("📜 *Hadith*\n\n%s%s\n\n%s\n\n*Reference:* %s, Book %d, #%d\n*Grade:* %s",
		escapeMarkdown(hdt.Arabic), narrator, escapeMarkdown(hdt.English), escapeMarkdown(colTitle), bookNum, hdt.HadithNumber, escapeMarkdown(grade))
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

	narratorHTML := ""
	if hadith.Narrator != "" {
		narratorHTML = fmt.Sprintf("\n\n<b>:</b> %s", html.EscapeString(hadith.Narrator))
	}

	return fmt.Sprintf("📜 <b>Hadith</b>\n\n%s%s\n\n%s\n\n<b>Ref:</b> %s, #%d\n<b>Grade:</b> %s",
		html.EscapeString(hadith.Arabic), narratorHTML, html.EscapeString(hadith.English), html.EscapeString(colTitle), hadith.HadithNumber, html.EscapeString(grade))
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

func splitTelegramMessage(text string, maxRunes int) []string {
	if maxRunes <= 0 {
		return []string{text}
	}

	runes := []rune(text)
	if len(runes) <= maxRunes {
		return []string{text}
	}

	chunks := make([]string, 0)
	start := 0

	for start < len(runes) {
		end := start + maxRunes
		if end >= len(runes) {
			chunks = append(chunks, strings.TrimSpace(string(runes[start:])))
			break
		}

		splitAt := -1
		searchFrom := start + (maxRunes * 2 / 3)
		for i := end; i >= searchFrom && i > start; i-- {
			if runes[i-1] == '\n' {
				splitAt = i
				break
			}
		}

		if splitAt == -1 {
			splitAt = end
		}

		chunk := strings.TrimSpace(string(runes[start:splitAt]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
		start = splitAt
	}

	return chunks
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"`", "\\`",
	)
	return replacer.Replace(text)
}
