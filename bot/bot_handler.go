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

func (self *BotHandler) handlePrivateMessage(update tgbotapi.Update) error {
	command := self.ParseCommand(update.Message.Text)

	switch command {
	case "start":
		return self.handleStart(update)
	}

	return nil
}

func (self *BotHandler) handleGroupMessage(update tgbotapi.Update) error {
	command := self.ParseCommand(update.Message.Text)

	switch command {
	case "gonew":
		return self.handleNewGame(update)
	case "gojoin":
		return self.handleJoinGame(update)
	case "goleave":
		return self.handleLeaveGame(update)
	case "goplayers":
		return self.handleGetPlayers(update)
	case "gostart":
		return self.handleGameStart(update)
	case "goopen":
		return self.handleOpenLobby(update)
	case "goclose":
		return self.handleCloseLobby(update)
	case "gokill":
		return self.handleKillGame(update)
	case "goinfo":
		return self.handleGameInfo(update)
	}

	return nil
}

func (self *BotHandler) ProcessUpdate(update tgbotapi.Update) error {
	if update.Message != nil {
		if update.Message.Chat.ID > 0 {
			return self.handlePrivateMessage(update)
		}

		return self.handleGroupMessage(update)
	}

	if update.InlineQuery != nil {
		return self.handleInlineQuery(update)
	}

	if update.ChosenInlineResult != nil {
		return self.handleInlineResult(update)
	}

	return nil
}

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
