package bot

import (
	"log"
	"os"
	"strings"

	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	gameManager GameManager
	bot         *tgbotapi.BotAPI
	logger      *log.Logger
}

func NewBotHandler(bot *tgbotapi.BotAPI, verbose bool) *BotHandler {
	var logger *log.Logger

	if verbose {
		logger = log.New(os.Stderr, "BotHandler", log.LstdFlags)
	}

	return &BotHandler{
		gameManager: GameManager{
			games:   make(map[int64]*uno.Game[*UnoPlayer]),
			players: make(map[int64]PlayerGame),
		},
		bot:    bot,
		logger: logger,
	}
}

func (self *BotHandler) ParseCommand(message string) string {
	// there is actually a parser for this with `update.Message.Command()`, but since
	// telegram sends UTF-16 entities, and that function treats it as UTF-8, it is not trustable.
	// moreover, that function doesn't check if the username after the "@" symbol matches
	// the bot's, so it could potentially make the bot answer messages it's not supposed to.
	if !strings.HasPrefix(message, "/") {
		return ""
	}

	message = strings.ToLower(message)
	idx := strings.Index(message, "@"+strings.ToLower(self.bot.Self.UserName))
	if idx == -1 {
		idx = strings.Index(message, " ")
	}

	if idx == -1 {
		return message[1:]
	}

	return message[1:idx]
}

// dispatches the update according to its fields,
// creating a copy of the structs so that the variables
// are not replaced by a subsequent update.
func (self *BotHandler) ProcessUpdate(update tgbotapi.Update) error {
	if update.Message != nil {
		// create a copy of the struct
		msg := *update.Message
		if update.Message.Chat.ID > 0 {
			return self.handlePrivateMessage(&msg)
		}

		return self.handleGroupMessage(&msg)
	}

	if update.InlineQuery != nil {
		inlineQuery := *update.InlineQuery
		return self.handleInlineQuery(&inlineQuery)
	}

	if update.ChosenInlineResult != nil {
		inlineResult := *update.ChosenInlineResult
		return self.handleInlineResult(&inlineResult)
	}

	return nil
}

/// Helper methods

func (self *BotHandler) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)

	_, err := self.bot.Send(msg)
	return err
}

func (self *BotHandler) SendMessageHTML(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := self.bot.Send(msg)
	return err
}

func (self *BotHandler) ReplyMessage(replyToMessageId int, chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.BaseChat.ReplyToMessageID = replyToMessageId

	_, err := self.bot.Send(msg)
	return err
}

func (self *BotHandler) ReplyMessageHTML(replyToMessageId int, chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.BaseChat.ReplyToMessageID = replyToMessageId

	_, err := self.bot.Send(msg)
	return err
}

func (self *BotHandler) NewInlineKeyboardButtonSwitchCurrentChat(text string, switchInlineQuery string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:                         text,
		SwitchInlineQueryCurrentChat: &switchInlineQuery,
	}
}

func (self *BotHandler) logDebug(v ...any) {
	if self.logger != nil {
		self.logger.Print(v...)
	}
}

func (self *BotHandler) logDebugf(format string, v ...any) {
	if self.logger != nil {
		self.logger.Printf(format, v...)
	}
}
