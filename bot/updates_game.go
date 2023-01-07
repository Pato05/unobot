package bot

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/Pato05/unobot/cards"
	"github.com/Pato05/unobot/messages"
	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: order deck cards by color and number (can be implemented by comparing card.GetGlobalIndex())
// TODO: make sure last card lookup works when implemented card ordering
// TODO: fix reverse when it is thrown by first player (at that point, the turn goes again to the first player) (fixed??)
// TODO: fix PreviousPlayer lookup (fixed??)
// TODO: fix asking for color when it is the last card (fixed??)
// TODO: fix that players can see each other's deck

type Game *uno.Game[*UnoPlayer]

func (self *BotHandler) handleInlineQuery(update tgbotapi.Update) error {
	userID := update.InlineQuery.From.ID
	game_, found := self.gameManager.GetPlayerGame(userID)
	game, player := game_.Game, game_.UnoPlayer
	if !found {
		_, err := self.bot.Send(tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results: []interface{}{
				tgbotapi.NewInlineQueryResultArticle(
					"join_game_error",
					"Join a game first!",
					"Join a game first!",
				),
			},
			CacheTime: 1,
		})
		return err
	}

	if !game.Started {
		_, err := self.bot.Send(tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results: []interface{}{
				tgbotapi.NewInlineQueryResultArticle(
					"game_not_started_error",
					"Let the game start first!",
					"Let the game start first!",
				),
			},
			CacheTime: 1,
		})
		return err
	}

	cards_ := player.Deck().Cards
	isPlayersTurn := userID == game.CurrentPlayer().UserId

	if isPlayersTurn && player.ShouldChooseColor() {
		return self.handleChooseColorInlineQuery(update, player, game)
	}

	extra := 0
	if isPlayersTurn {
		extra = 2
		if game.CanCallBluff {
			extra += 1
		}
	}
	cardsLength := len(cards_)
	results := make([]interface{}, cardsLength+extra)

	if isPlayersTurn {

		if game.DidJustDraw {
			results[0] = PassAction()
		} else {
			results[0] = DrawAction(game.DrawCounter)
		}
		results[1] = GameInfoAction(game)
		if game.CanCallBluff {
			results[2] = results[1]
			results[1] = BluffAction()
		}
	}

	gameInfoStr := GetGameInfoStr(game)
	Cards := make([]IndexedStruct[*cards.Card], cardsLength)
	for index, card := range cards_ {
		Cards[index] = IndexedStruct[*cards.Card]{
			Value: &card,
			Index: index,
		}
	}
	sort.Slice(Cards, func(i, j int) bool {
		return Cards[i].Value.GetGlobalIndex() < Cards[j].Value.GetGlobalIndex()
	})

	for index, el := range Cards {
		i := index + extra
		card := el.Value

		canPlayCard := isPlayersTurn && game.CanCurrentPlayerPlayCard(card)
		if canPlayCard && game.DidJustDraw {
			canPlayCard = el.Index == cardsLength-1
		}

		if canPlayCard {
			cardDigest := int(card.Color)*100 + int(card.CardIndex)
			results[i] = tgbotapi.InlineQueryResultCachedSticker{
				Type:      "sticker",
				ID:        fmt.Sprintf("play_%d_%d", index, cardDigest),
				StickerID: card.GetFileID().Normal,
			}
		} else {
			results[i] = tgbotapi.InlineQueryResultCachedSticker{
				Type:      "sticker",
				ID:        fmt.Sprintf("grey_%d", i),
				StickerID: card.GetFileID().Grey,
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:      gameInfoStr,
					ParseMode: tgbotapi.ModeHTML,
				},
			}
		}

	}

	_, err := self.bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (self *BotHandler) handleChooseColorInlineQuery(update tgbotapi.Update, player *UnoPlayer, game *uno.Game[*UnoPlayer]) error {
	colorsLength := len(choosableColorsString)
	results := make([]interface{}, colorsLength+1)
	for i, val := range choosableColorsString {
		// starts from blue, index 1
		results[int(i)-1] = tgbotapi.NewInlineQueryResultArticle("choosecolor_"+strconv.Itoa(int(i)), val, val)
	}

	var cardsStr strings.Builder
	for i, card := range player.deck.Cards {
		if i > 0 {
			cardsStr.WriteString(", ")
		}
		cardsStr.WriteString(CardColorEmojis[card.Color] + " ")
		if card.IsSpecial() {
			cardsStr.WriteString(CardSpecials[card.Special])
		} else {
			cardsStr.WriteString(strconv.Itoa(int(card.Number)))
		}
	}

	results[colorsLength] = tgbotapi.InlineQueryResultArticle{
		Type:        "article",
		ID:          "getinfo",
		Title:       "Cards (tap for game info)",
		Description: cardsStr.String(),
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text:      GetGameInfoStr(game),
			ParseMode: tgbotapi.ModeHTML,
		},
	}

	_, err := self.bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (self *BotHandler) handleInlineResult(update tgbotapi.Update) error {
	id := update.ChosenInlineResult.ResultID

	if strings.HasPrefix(id, "play_") {
		return self.handlePlayInlineResult(update)
	}

	if strings.HasPrefix(id, "choosecolor_") {
		return self.handleChooseColorInlineResult(update)
	}

	switch id {
	case "draw_card":
		return self.handleDrawInlineResult(update)
	case "pass_turn":
		return self.handlePassInlineResult(update)
	case "call_bluff":
		return self.handleBluffInlineResult(update)
	}

	return nil
}

func (self *BotHandler) handlePassInlineResult(update tgbotapi.Update) error {
	user := update.ChosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().UserId != user.ID {
		// fail silently
		return nil
	}

	game.NextPlayer()
	return self.nextPlayer(update, playerGame.GameChatId, game)
}

func (self *BotHandler) handleDrawInlineResult(update tgbotapi.Update) error {
	user := update.ChosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().UserId != user.ID {
		// fail silently
		return nil
	}

	err := game.CurrentPlayerDraw()
	if err != nil {
		self.logDebug(err)
		_, err := self.bot.Send(tgbotapi.NewMessage(playerGame.GameChatId, err.Error()))
		return err
	}

	return self.nextPlayer(update, playerGame.GameChatId, game)
}

func (self *BotHandler) handleBluffInlineResult(update tgbotapi.Update) error {
	user := update.ChosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().UserId != user.ID {
		// fail silently
		return nil
	}

	didBluff := game.CallBluff()
	if didBluff {
		self.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: playerGame.GameChatId,
			},
			Text:      fmt.Sprintf("%s bluffed, giving them 4 cards.", game.PreviousPlayer().HTML()),
			ParseMode: tgbotapi.ModeHTML,
		})
	} else {
		self.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: playerGame.GameChatId,
			},
			Text:      fmt.Sprintf("%s didn't bluff, giving 6 cards to %s", game.PreviousPlayer().Name, game.CurrentPlayer().HTML()),
			ParseMode: tgbotapi.ModeHTML,
		})
	}

	game.NextPlayer()
	return self.nextPlayer(update, playerGame.GameChatId, game)
}

func (self *BotHandler) handleChooseColorInlineResult(update tgbotapi.Update) error {
	id := update.ChosenInlineResult.ResultID
	user := update.ChosenInlineResult.From
	split := strings.Split(id, "_")
	if len(split) != 2 {
		log.Println("Wrong ResultID:", id)
		return nil
	}
	newColorInt, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal(err)
		return nil
	}

	newColor := cards.CardColor(newColorInt)

	if newColor > cards.Yellow {
		// invalid color, should range 0-4
		return nil
	}

	game_, found := self.gameManager.GetPlayerGame(user.ID)
	if !found {
		return nil
	}

	game := game_.Game

	if game.CurrentPlayer().UserId != user.ID {
		// fail silently
		return nil
	}

	game.ChooseColor(newColor)
	game.NextPlayer()

	return self.nextPlayer(update, game_.GameChatId, game)
}

func (self *BotHandler) handlePlayInlineResult(update tgbotapi.Update) error {
	id := update.ChosenInlineResult.ResultID
	user := update.ChosenInlineResult.From
	split := strings.Split(id, "_")
	if len(split) != 3 {
		log.Print("Wrong ResultID:", id)
		return nil
	}

	cardIndex, err := strconv.Atoi(split[1])
	if err != nil {
		// fail silently
		log.Fatal(err)
		return nil
	}

	digest, err := strconv.Atoi(split[2])
	if err != nil {
		log.Fatal(err)
		return nil
	}

	playerGame := self.gameManager.players[user.ID]
	game, player := playerGame.Game, playerGame.UnoPlayer

	if len(player.Deck().Cards) < cardIndex {
		self.logDebug("Wrong index: got ", cardIndex, ", array is smaller than given index.")
		return nil
	}
	card := player.Deck().Cards[cardIndex]

	d := int(card.Color)*100 + int(card.CardIndex)

	if d != digest {
		// current card's index does not match previous card
		self.logDebug("Wrong card color and index: expecting ", d, " got ", digest)
		return nil
	}

	if game.CurrentPlayer().UserId != user.ID {
		// fail silently
		return nil
	}

	err = game.PlayCard(&card)
	if err != nil {
		_, err := self.bot.Send(tgbotapi.NewMessage(playerGame.GameChatId, err.Error()))
		return err
	}

	player.RemoveCard(uint(cardIndex))
	game.Deck.Discard(card)

	if player.ShouldShoutUNO() {
		self.bot.Send(tgbotapi.NewMessage(playerGame.GameChatId, "UNO!"))
	}

	if player.ShouldChooseColor() {
		msg := tgbotapi.NewMessage(playerGame.GameChatId, "Please choose a color")
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
				self.NewInlineKeyboardButtonSwitchCurrentChat("Choose a color", ""),
			}},
		}
		_, err := self.bot.Send(msg)
		return err
	}

	if player.DidWin() {
		err := self.playerWon(update, playerGame.GameChatId, game, player)
		if err != nil {
			return nil
		}

		return self.nextPlayer(update, playerGame.GameChatId, game)
	}

	game.NextPlayer()

	return self.nextPlayer(update, playerGame.GameChatId, playerGame.Game)
}

func (self *BotHandler) handleNewGame(update tgbotapi.Update) error {
	err := self.gameManager.NewGame(update.FromChat().ID, update.SentFrom().ID)
	if err != nil {
		_, err := self.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
		return err
	}

	_, err = self.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, messages.GAME_CREATED_SUCCESS))
	return err
}

func (self *BotHandler) handleJoinGame(update tgbotapi.Update) error {
	err := self.gameManager.PlayerJoin(update.FromChat().ID, update.SentFrom())
	if err != nil {
		_, err := self.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, err.Error()))
		return err
	}

	_, err = self.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           update.FromChat().ID,
			ReplyToMessageID: update.Message.MessageID,
		},
		Text: messages.GAME_JOINED_SUCCESS,
	})
	return err
}

func (self *BotHandler) handleLeaveGame(update tgbotapi.Update) error {
	lastPlayer, _ := self.gameManager.PlayerLeave(update.FromChat().ID, update.SentFrom().ID)

	err := self.ReplyMessage(update.Message.MessageID, update.FromChat().ID, messages.GAME_LEAVE_SUCCESS)

	if err != nil {
		self.logDebug(err)
	}

	if lastPlayer != nil {
		return self.playerWonMsg(update, update.FromChat().ID, lastPlayer)
	}

	return nil
}

func (self *BotHandler) handleGetPlayers(update tgbotapi.Update) error {
	players, err := self.gameManager.GetPlayersInGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}
	if len(players) == 0 {
		return self.SendMessage(update.FromChat().ID, messages.NO_PLAYERS_IN_GAME_ERROR)
	}
	var message strings.Builder
	for _, player := range players {
		message.WriteString("- " + player.HTML() + "\n")
	}
	return self.SendMessageHTML(update.FromChat().ID, message.String())
}

func (self *BotHandler) handleGameInfo(update tgbotapi.Update) error {
	game, err := self.gameManager.GetGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	return self.SendMessage(update.FromChat().ID, fmt.Sprintf(
		"Game Info:\n"+
			"  Deck:\n"+
			"    Available cards: %d\n"+
			"    Discarded cards: %d",
		len(game.Deck.Cards),
		len(game.Deck.Discarded),
	))
}

func (self *BotHandler) playerWon(update tgbotapi.Update, chatId int64, game *uno.Game[*UnoPlayer], player *UnoPlayer) error {
	self.playerWonMsg(update, chatId, player)
	delete(self.gameManager.players, player.UserId)
	err := game.CurrentPlayerWon()
	if err != nil {
		self.SendMessage(chatId, "Game ended!")
		self.gameManager.DeleteGame(chatId)
		return err
	}

	return nil
}

func (self *BotHandler) playerWonMsg(update tgbotapi.Update, chatId int64, player *UnoPlayer) error {
	return self.SendMessageHTML(chatId, player.HTML()+" won!")
}

func (self *BotHandler) nextPlayer(update tgbotapi.Update, chatId int64, game *uno.Game[*UnoPlayer]) error {
	nextPlayer := game.CurrentPlayer()
	_, err := self.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           chatId,
			ReplyToMessageID: 0,
			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
					self.NewInlineKeyboardButtonSwitchCurrentChat("Play a card", ""),
				}},
			},
		},
		Text:      "Next player: " + nextPlayer.HTML(),
		ParseMode: tgbotapi.ModeHTML,
	})

	return err
}

func (self *BotHandler) handleGameStart(update tgbotapi.Update) error {
	game, err := self.gameManager.GetGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}
	firstCard, err := game.Start()
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}
	_, err = self.bot.Send(tgbotapi.NewSticker(update.FromChat().ID, tgbotapi.FileID(firstCard.GetFileID().Normal)))
	if err != nil {
		return err
	}
	game.NextPlayer()
	return self.nextPlayer(update, update.FromChat().ID, game)
}

func (self *BotHandler) handleCloseLobby(update tgbotapi.Update) error {
	game, err := self.gameManager.GetGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	if game.GameCreatorUID != update.SentFrom().ID {
		return self.ReplyMessage(update.Message.MessageID, update.FromChat().ID, "You're not allowed to do this!")
	}

	err = game.CloseLobby()
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	return self.SendMessage(update.FromChat().ID, messages.LOBBY_CLOSED_SUCCESS)
}

func (self *BotHandler) handleOpenLobby(update tgbotapi.Update) error {
	game, err := self.gameManager.GetGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	if game.GameCreatorUID != update.SentFrom().ID {
		return self.ReplyMessage(update.Message.MessageID, update.FromChat().ID, "You're not allowed to do this!")
	}

	err = game.OpenLobby()
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	return self.SendMessage(update.FromChat().ID, messages.LOBBY_OPEN_SUCCESS)
}

func (self *BotHandler) handleKillGame(update tgbotapi.Update) error {
	game, err := self.gameManager.GetGame(update.FromChat().ID)
	if err != nil {
		return self.SendMessage(update.FromChat().ID, err.Error())
	}

	if game.GameCreatorUID != update.SentFrom().ID {
		return self.ReplyMessage(update.Message.MessageID, update.FromChat().ID, "You're not allowed to do this!")
	}

	self.gameManager.DeleteGame(update.FromChat().ID)
	return self.SendMessage(update.FromChat().ID, "Game killed successfully.")
}
