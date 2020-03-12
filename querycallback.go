package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func showQuotes(update *tgbotapi.Update) {
	mp  := make(map[string]string)
	err := json.Unmarshal([]byte(update.CallbackQuery.Data), &mp)
	if err != nil {
		Log.Error(err)
		return
	}
	if mp["arg"] == "" {
		return
	}
	Log.Debug("recv callback:", update.CallbackQuery.Data, "chatID:", update.CallbackQuery.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
	msg.Text = "@" + update.CallbackQuery.From.UserName + " - " + update.CallbackQuery.From.FirstName + "\n"
	msg.ParseMode = "Markdown"

	msg.Text += _quotes(mp["arg"])

	send(&msg, 5)
}
