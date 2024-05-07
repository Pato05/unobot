package bot

import (
	"fmt"
	"strings"

	"github.com/Pato05/unobot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (self *BotHandler) handleGroupMessage(message *tgbotapi.Message) error {
	command := self.ParseCommand(message.Text)

	switch command {
	case "gonew":
		return self.handleNewGame(message)
	case "gojoin":
		return self.handleJoinGame(message)
	case "goleave":
		return self.handleLeaveGame(message)
	case "goplayers":
		return self.handleGetPlayers(message)
	case "gostart":
		return self.handleGameStart(message)
	case "goopen":
		return self.handleOpenLobby(message)
	case "goclose":
		return self.handleCloseLobby(message)
	case "gokill":
		return self.handleKillGame(message)
	case "goinfo":
		return self.handleGameInfo(message)
	}

	return nil
}

func (self *BotHandler) handleNewGame(message *tgbotapi.Message) error {
	err := self.gameManager.NewGame(message.Chat.ID, message.From.ID)
	if err != nil {
		_, err := self.bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return err
	}

	_, err = self.bot.Send(tgbotapi.NewMessage(message.Chat.ID, messages.GAME_CREATED_SUCCESS))
	return err
}

func (self *BotHandler) handleJoinGame(message *tgbotapi.Message) error {
	err := self.gameManager.PlayerJoin(message.Chat.ID, message.From)
	if err != nil {
		_, err := self.bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return err
	}

	_, err = self.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.Chat.ID,
			ReplyToMessageID: message.MessageID,
		},
		Text: messages.GAME_JOINED_SUCCESS,
	})
	return err
}

func (self *BotHandler) handleLeaveGame(message *tgbotapi.Message) error {
	lastPlayer, _ := self.gameManager.PlayerLeave(message.Chat.ID, message.From.ID)

	err := self.ReplyMessage(message.MessageID, message.Chat.ID, messages.GAME_LEAVE_SUCCESS)

	if err != nil {
		self.logDebug(err)
	}

	if lastPlayer != nil {
		return self.playerWonMsg(message.Chat.ID, lastPlayer)
	}

	return nil
}

func (self *BotHandler) handleGetPlayers(message *tgbotapi.Message) error {
	players, err := self.gameManager.GetPlayersInGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}
	if len(players) == 0 {
		return self.SendMessage(message.Chat.ID, messages.NO_PLAYERS_IN_GAME_ERROR)
	}
	var response strings.Builder
	for _, player := range players {
		response.WriteString("- " + player.HTML() + "\n")
	}
	return self.SendMessageHTML(message.Chat.ID, response.String())
}

func (self *BotHandler) handleGameInfo(message *tgbotapi.Message) error {
	game, err := self.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	return self.SendMessage(message.Chat.ID, fmt.Sprintf(
		"Game Info:\n"+
			"  Deck:\n"+
			"    Available cards: %d\n"+
			"    Discarded cards: %d",
		len(game.Deck.Cards),
		len(game.Deck.Discarded),
	))
}

func (self *BotHandler) handleGameStart(message *tgbotapi.Message) error {
	game, err := self.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}
	firstCard, err := game.Start()
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}
	_, err = self.bot.Send(tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(firstCard.GetFileID().Normal)))
	if err != nil {
		return err
	}
	game.NextPlayer()
	return self.nextPlayer(message.Chat.ID, game)
}

func (self *BotHandler) handleCloseLobby(message *tgbotapi.Message) error {
	game, err := self.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return self.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	err = game.CloseLobby()
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	return self.SendMessage(message.Chat.ID, messages.LOBBY_CLOSED_SUCCESS)
}

func (self *BotHandler) handleOpenLobby(message *tgbotapi.Message) error {
	game, err := self.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return self.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	err = game.OpenLobby()
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	return self.SendMessage(message.Chat.ID, messages.LOBBY_OPEN_SUCCESS)
}

func (self *BotHandler) handleKillGame(message *tgbotapi.Message) error {
	game, err := self.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return self.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return self.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	self.gameManager.DeleteGame(message.Chat.ID)
	return self.SendMessage(message.Chat.ID, "Game killed successfully.")
}
