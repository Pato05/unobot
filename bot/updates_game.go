package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Pato05/unobot/cards"
	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: order deck cards by color and number (can be implemented by comparing card.GetGlobalIndex())
// TODO: make sure last card lookup works when implemented card ordering
// TODO: fix reverse when it is thrown by first player (at that point, the turn goes again to the first player) (fixed??)
// TODO: fix PreviousPlayer lookup (fixed??)
// TODO: fix asking for color when it is the last card (fixed??)
// TODO: fix that players can see each other's deck
// TODO: fix that if player throws a Wild card and leaves, others can't play a card other than Wild ones

type Game *uno.Game[*UnoPlayer]

func (self *BotHandler) handleInlineQuery(inlineQuery *tgbotapi.InlineQuery) error {
	userID := inlineQuery.From.ID
	game_, found := self.gameManager.GetPlayerGame(userID)
	game, player := game_.Game, game_.UnoPlayer
	if !found {
		_, err := self.bot.Send(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
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
			InlineQueryID: inlineQuery.ID,
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

	isPlayersTurn := userID == game.CurrentPlayer().GetUID()

	if isPlayersTurn && player.ShouldChooseColor() {
		return self.handleChooseColorInlineQuery(inlineQuery, player, game)
	}

	extra := 0
	if isPlayersTurn {
		extra = 2
		if game.CanCallBluff {
			extra += 1
		}
	}

	Cards := player.Deck().Cards
	cardsLength := len(Cards)
	results := make([]interface{}, cardsLength+extra)

	if isPlayersTurn {
		if game.DidJustDraw {
			results[0] = PassAction()
		} else {
			results[0] = DrawAction(game.DrawCounter)
		}
		results[1] = GameInfoAction(game)
		if game.CanCallBluff {
			results[2] = results[1] // shift result
			results[1] = BluffAction()
		}
	}

	gameInfoStr := GetGameInfoStr(game)

	// Cards := make([]IndexedStruct[*cards.Card], cardsLength)
	// for index, card := range Cards {
	// 	Cards[index] = IndexedStruct[*cards.Card]{
	// 		Value: &card,
	// 		Index: index,
	// 	}
	// }
	// sort.Slice(Cards, func(i, j int) bool {
	// 	return Cards[i].Value.GetGlobalIndex() < Cards[j].Value.GetGlobalIndex()
	// })

	for index, card := range Cards {
		i := index + extra

		canPlayCard := isPlayersTurn && game.CanCurrentPlayerPlayCard(&card)
		if canPlayCard && game.DidJustDraw {
			canPlayCard = index == cardsLength-1
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
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (self *BotHandler) handleChooseColorInlineQuery(inlineQuery *tgbotapi.InlineQuery, player *UnoPlayer, game *uno.Game[*UnoPlayer]) error {
	colorsLength := len(choosableColorsString)
	results := make([]interface{}, colorsLength+1)
	for i, val := range choosableColorsString {
		// starts from blue, index 1
		results[int(i)-1] = tgbotapi.NewInlineQueryResultArticle("choosecolor_"+strconv.Itoa(int(i)), val, val)
	}

	var cardsStr strings.Builder
	for i, card := range player.Deck().Cards {
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
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (self *BotHandler) handleInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	id := chosenInlineResult.ResultID

	if strings.HasPrefix(id, "play_") {
		return self.handlePlayInlineResult(chosenInlineResult)
	}

	if strings.HasPrefix(id, "choosecolor_") {
		return self.handleChooseColorInlineResult(chosenInlineResult)
	}

	switch id {
	case "draw_card":
		return self.handleDrawInlineResult(chosenInlineResult)
	case "pass_turn":
		return self.handlePassInlineResult(chosenInlineResult)
	case "call_bluff":
		return self.handleBluffInlineResult(chosenInlineResult)
	}

	self.logger.Print("WARN: unknown inline result id? ", id)

	return nil
}

func (self *BotHandler) handlePassInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	// sort deck because user took card and passed the turn
	playerGame.UnoPlayer.Deck().Sort()
	game.NextPlayer()
	return self.nextPlayer(playerGame.GameChatId, game)
}

func (self *BotHandler) handleDrawInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	err := game.CurrentPlayerDraw()
	if err != nil {
		self.logDebug(err)
		_, err := self.bot.Send(tgbotapi.NewMessage(playerGame.GameChatId, err.Error()))
		return err
	}

	return self.nextPlayer(playerGame.GameChatId, game)
}

func (self *BotHandler) handleBluffInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := self.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
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
			Text:      fmt.Sprintf("%s didn't bluff, giving 6 cards to %s", game.PreviousPlayer().EscapedName(), game.CurrentPlayer().HTML()),
			ParseMode: tgbotapi.ModeHTML,
		})
	}

	game.NextPlayer()
	return self.nextPlayer(playerGame.GameChatId, game)
}

func (self *BotHandler) handleChooseColorInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	id := chosenInlineResult.ResultID
	user := chosenInlineResult.From
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

	if newColorInt < 1 || newColorInt > 4 {
		log.Println("invalid color, should range 1-4")
		return nil
	}

	game_, found := self.gameManager.GetPlayerGame(user.ID)
	if !found {
		return nil
	}

	game := game_.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	game.ChooseColor(newColor)
	if err := self.checkPlayerWon(game_); err != nil {

	}
	game.NextPlayer()

	return self.nextPlayer(game_.GameChatId, game)
}

func (self *BotHandler) handlePlayInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	id := chosenInlineResult.ResultID
	user := chosenInlineResult.From
	split := strings.Split(id, "_")
	if len(split) != 3 {
		self.logDebug("Wrong ResultID:", id)
		return nil
	}

	cardIndex, err := strconv.Atoi(split[1])
	if err != nil {
		// fail silently
		self.logDebug(err)
		return nil
	}

	digest, err := strconv.Atoi(split[2])
	if err != nil {
		self.logDebug(err)
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

	if game.CurrentPlayer().GetUID() != user.ID {
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

	if err := self.checkPlayerWon(playerGame); err != nil {
		return err
	}

	game.NextPlayer()

	return self.nextPlayer(playerGame.GameChatId, playerGame.Game)
}

func (self *BotHandler) playerWonMsg(chatId int64, player *UnoPlayer) error {
	return self.SendMessageHTML(chatId, player.HTML()+" won!")
}

func (self *BotHandler) playerWon(chatId int64, game *uno.Game[*UnoPlayer], player *UnoPlayer) error {
	self.playerWonMsg(chatId, player)
	delete(self.gameManager.players, player.GetUID())
	err := game.CurrentPlayerWon()
	if err != nil {
		self.SendMessage(chatId, "Game ended!")
		self.gameManager.DeleteGame(chatId)
		return err
	}

	return nil
}

func (self *BotHandler) checkPlayerWon(playerGame PlayerGame) error {
	if playerGame.UnoPlayer.DidWin() {
		err := self.playerWon(playerGame.GameChatId, playerGame.Game, playerGame.UnoPlayer)
		if err != nil {
			return err
		}

		return self.nextPlayer(playerGame.GameChatId, playerGame.Game)
	}

	return nil
}

func (self *BotHandler) nextPlayer(chatId int64, game *uno.Game[*UnoPlayer]) error {
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
