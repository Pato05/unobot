package bot

import (
	"log"
	"os"
	"strings"

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
			games:   make(map[int64]*UnoGame),
			players: make(map[int64]PlayerGame),
		},
		bot:    bot,
		logger: logger,
	}
}

func (bh *BotHandler) ParseCommand(message string) string {
	// there is actually a parser for this with `update.Message.Command()`, but since
	// telegram sends UTF-16 entities, and that function treats it as UTF-8, it is not trustable.
	// moreover, that function doesn't check if the username after the "@" symbol matches
	// the bot's, so it could potentially make the bot answer messages it's not supposed to.
	if !strings.HasPrefix(message, "/") {
		return ""
	}

	message = strings.ToLower(message)
	idx := strings.Index(message, "@"+strings.ToLower(bh.bot.Self.UserName))
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
func (bh *BotHandler) ProcessUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		// create a copy of the struct
		msg := *update.Message
		if update.Message.Chat.Type == "private" {
			go bh.handlePrivateMessage(&msg)
			return
		}

		if update.Message.Chat.Type == "channel" {
			bh.logDebug("Ignoring message from chat type: channel")
			return
		}

		go bh.handleGroupMessage(&msg)
		return
	}

	if update.InlineQuery != nil {
		inlineQuery := *update.InlineQuery
		go bh.handleInlineQuery(&inlineQuery)
		return
	}

	if update.ChosenInlineResult != nil {
		inlineResult := *update.ChosenInlineResult
		go bh.handleInlineResult(&inlineResult)
		return
	}
}

/// Helper methods

func (bh *BotHandler) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)

	_, err := bh.bot.Send(msg)
	return err
}

func (bh *BotHandler) SendMessageHTML(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := bh.bot.Send(msg)
	return err
}

func (bh *BotHandler) ReplyMessage(replyToMessageId int, chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.BaseChat.ReplyToMessageID = replyToMessageId

	_, err := bh.bot.Send(msg)
	return err
}

func (bh *BotHandler) ReplyMessageHTML(replyToMessageId int, chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.BaseChat.ReplyToMessageID = replyToMessageId

	_, err := bh.bot.Send(msg)
	return err
}

func (bh *BotHandler) NewInlineKeyboardButtonSwitchCurrentChat(text string, switchInlineQuery string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:                         text,
		SwitchInlineQueryCurrentChat: &switchInlineQuery,
	}
}

func (bh *BotHandler) logDebug(v ...any) {
	if bh.logger != nil {
		bh.logger.Print(v...)
	}
}

func (bh *BotHandler) logDebugf(format string, v ...any) {
	if bh.logger != nil {
		bh.logger.Printf(format, v...)
	}
}
