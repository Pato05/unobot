package bot

import (
	"fmt"
	"strings"

	"github.com/Pato05/unobot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bh *BotHandler) handleGroupMessage(message *tgbotapi.Message) error {
	command := bh.ParseCommand(message.Text)

	switch command {
	case "gonew":
		return bh.handleNewGame(message)
	case "gojoin":
		return bh.handleJoinGame(message)
	case "goleave":
		return bh.handleLeaveGame(message)
	case "goplayers":
		return bh.handleGetPlayers(message)
	case "gostart":
		return bh.handleGameStart(message)
	case "goopen":
		return bh.handleOpenLobby(message)
	case "goclose":
		return bh.handleCloseLobby(message)
	case "gokill":
		return bh.handleKillGame(message)
	case "goinfo":
		return bh.handleGameInfo(message)
	}

	return nil
}

func (bh *BotHandler) handleNewGame(message *tgbotapi.Message) error {
	err := bh.gameManager.NewGame(message.Chat.ID, message.From.ID)
	if err != nil {
		_, err := bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return err
	}

	_, err = bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, messages.GAME_CREATED_SUCCESS))
	return err
}

func (bh *BotHandler) handleJoinGame(message *tgbotapi.Message) error {
	err := bh.gameManager.PlayerJoin(message.Chat.ID, message.From)
	if err != nil {
		_, err := bh.bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
		return err
	}

	_, err = bh.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.Chat.ID,
			ReplyToMessageID: message.MessageID,
		},
		Text: messages.GAME_JOINED_SUCCESS,
	})
	return err
}

func (bh *BotHandler) handleLeaveGame(message *tgbotapi.Message) error {
	lastPlayer, _ := bh.gameManager.PlayerLeave(message.Chat.ID, message.From.ID)

	err := bh.ReplyMessage(message.MessageID, message.Chat.ID, messages.GAME_LEAVE_SUCCESS)

	if err != nil {
		bh.logDebug(err)
	}

	if lastPlayer != nil {
		return bh.playerWonMsg(message.Chat.ID, lastPlayer)
	}

	return nil
}

func (bh *BotHandler) handleGetPlayers(message *tgbotapi.Message) error {
	players, err := bh.gameManager.GetPlayersInGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}
	if len(players) == 0 {
		return bh.SendMessage(message.Chat.ID, messages.NO_PLAYERS_IN_GAME_ERROR)
	}
	var response strings.Builder
	for _, player := range players {
		response.WriteString("- " + player.HTML() + "\n")
	}
	return bh.SendMessageHTML(message.Chat.ID, response.String())
}

func (bh *BotHandler) handleGameInfo(message *tgbotapi.Message) error {
	game, err := bh.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	return bh.SendMessage(message.Chat.ID, fmt.Sprintf(
		"Game Info:\n"+
			"  Deck:\n"+
			"    Available cards: %d\n"+
			"    Discarded cards: %d",
		len(game.Deck.Cards),
		len(game.Deck.Discarded),
	))
}

func (bh *BotHandler) handleGameStart(message *tgbotapi.Message) error {
	game, err := bh.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}
	firstCard, err := game.Start()
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}
	_, err = bh.bot.Send(tgbotapi.NewSticker(message.Chat.ID, tgbotapi.FileID(firstCard.GetFileID().Normal)))
	if err != nil {
		return err
	}

	game.NextPlayer()
	game.SetTimer(bh)
	return bh.announceNextPlayer(game)
}

func (bh *BotHandler) handleCloseLobby(message *tgbotapi.Message) error {
	game, err := bh.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return bh.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	err = game.CloseLobby()
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	return bh.SendMessage(message.Chat.ID, messages.LOBBY_CLOSED_SUCCESS)
}

func (bh *BotHandler) handleOpenLobby(message *tgbotapi.Message) error {
	game, err := bh.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return bh.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	err = game.OpenLobby()
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	return bh.SendMessage(message.Chat.ID, messages.LOBBY_OPEN_SUCCESS)
}

func (bh *BotHandler) handleKillGame(message *tgbotapi.Message) error {
	game, err := bh.gameManager.GetGame(message.Chat.ID)
	if err != nil {
		return bh.SendMessage(message.Chat.ID, err.Error())
	}

	if game.GameCreatorUID != message.From.ID {
		return bh.ReplyMessage(message.MessageID, message.Chat.ID, "You're not allowed to do this!")
	}

	bh.gameManager.DeleteGame(message.Chat.ID)
	return bh.SendMessage(message.Chat.ID, "Game killed successfully.")
}
