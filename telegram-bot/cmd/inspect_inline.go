package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tucnak/telebot"
)

func main() {
	hadithText := "🵿 *Hadith*\n\n*Arabic:*\n...arabic...\n\n*English:*\nThis is English text.\n\n*Reference:* Sunan an-Nasa'i, Book 34, Hadith #2517\n*Grade:* Sahih\n"

	var content telebot.InputMessageContent = &telebot.InputTextMessageContent{
		Text:      hadithText,
		ParseMode: "MarkdownV2",
	}

	article := &telebot.ArticleResult{
		ResultBase: telebot.ResultBase{
			ID:      fmt.Sprintf("test_%d", time.Now().UnixNano()),
			Content: &content,
		},
		Title:       "Test",
		Description: "desc",
		Text:        "preview",
	}

	results := telebot.Results{article}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
