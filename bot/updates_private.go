package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (self *BotHandler) handlePrivateMessage(message *tgbotapi.Message) error {
	command := self.ParseCommand(message.Text)

	switch command {
	case "start":
		return self.handleStart(message)
	}

	return nil
}

func (self *BotHandler) handleStart(message *tgbotapi.Message) error {
	uri := "https://t.me/" + self.bot.Self.UserName + "?startgroup=1"
	msg := tgbotapi.NewMessage(message.Chat.ID, "Hello! I am Go UNO bot, an UNO bot written in Go. Add me to a group to play with your friends!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "➕ Add me to a group!", URL: &uri},
		),
	)
	_, err := self.bot.Send(msg)

	return err
}
