package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Pato05/unobot/cards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: order deck cards by color and number (can be implemented by comparing card.GetGlobalIndex())
// TODO: make sure last card lookup works when implemented card ordering
// TODO: fix reverse when it is thrown by first player (at that point, the turn goes again to the first player) (fixed??)
// TODO: fix PreviousPlayer lookup (fixed??)
// TODO: fix asking for color when it is the last card (fixed??)
// TODO: fix that players can see each other's deck
// TODO: fix that if player throws a Wild card and leaves, others can't play a card other than Wild ones

type Game *UnoGame

func (bh *BotHandler) handleInlineQuery(inlineQuery *tgbotapi.InlineQuery) error {
	userID := inlineQuery.From.ID
	game_, found := bh.gameManager.GetPlayerGame(userID)
	game, player := game_.Game, game_.UnoPlayer
	if !found {
		_, err := bh.bot.Send(tgbotapi.InlineConfig{
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
		_, err := bh.bot.Send(tgbotapi.InlineConfig{
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
		return bh.handleChooseColorInlineQuery(inlineQuery, player, game)
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
				ID:        fmt.Sprintf("gray_%d", i),
				StickerID: card.GetFileID().Gray,
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:      gameInfoStr,
					ParseMode: tgbotapi.ModeHTML,
				},
			}
		}

	}

	_, err := bh.bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (bh *BotHandler) handleChooseColorInlineQuery(inlineQuery *tgbotapi.InlineQuery, player *UnoPlayer, game *UnoGame) error {
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

	_, err := bh.bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		CacheTime:     1,
	})

	return err
}

func (bh *BotHandler) handleInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	id := chosenInlineResult.ResultID

	if strings.HasPrefix(id, "play_") {
		return bh.handlePlayInlineResult(chosenInlineResult)
	}

	if strings.HasPrefix(id, "choosecolor_") {
		return bh.handleChooseColorInlineResult(chosenInlineResult)
	}

	if strings.HasPrefix(id, "gray_") {
		return nil
	}

	switch id {
	case "draw_card":
		return bh.handleDrawInlineResult(chosenInlineResult)
	case "pass_turn":
		return bh.handlePassInlineResult(chosenInlineResult)
	case "call_bluff":
		return bh.handleBluffInlineResult(chosenInlineResult)
	case "gameinfo":
		return nil
	}

	bh.logger.Print("WARN: unknown inline result id? ", id)

	return nil
}

func (bh *BotHandler) handlePassInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := bh.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	// sort deck because user took card and passed the turn
	playerGame.UnoPlayer.Deck().Sort()
	return bh.nextPlayer(game)
}

func (bh *BotHandler) handleDrawInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := bh.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	err := game.CurrentPlayerDraw()
	if err != nil {
		bh.logDebug(err)
		_, err := bh.bot.Send(tgbotapi.NewMessage(game.ChatId, err.Error()))
		return err
	}

	playerGame.Game.SetTimer(bh)
	return bh.announceNextPlayer(game)
}

func (bh *BotHandler) handleBluffInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	user := chosenInlineResult.From
	playerGame := bh.gameManager.players[user.ID]
	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	didBluff := game.CallBluff()
	if didBluff {
		bh.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: game.ChatId,
			},
			Text:      fmt.Sprintf("%s bluffed, giving them 4 cards.", game.PreviousPlayer().HTML()),
			ParseMode: tgbotapi.ModeHTML,
		})
	} else {
		bh.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: game.ChatId,
			},
			Text:      fmt.Sprintf("%s didn't bluff, giving 6 cards to %s", game.PreviousPlayer().EscapedName(), game.CurrentPlayer().HTML()),
			ParseMode: tgbotapi.ModeHTML,
		})
	}

	return bh.nextPlayer(game)
}

func (bh *BotHandler) handleChooseColorInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
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

	if newColorInt < 1 || newColorInt > 4 {
		log.Println("invalid color, should range 1-4")
		return nil
	}

	newColor := cards.CardColor(newColorInt)

	playerGame, found := bh.gameManager.GetPlayerGame(user.ID)
	if !found {
		return nil
	}

	game := playerGame.Game

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	game.ChooseColor(newColor)

	return bh.nextPlayer(game)
}

func (bh *BotHandler) handlePlayInlineResult(chosenInlineResult *tgbotapi.ChosenInlineResult) error {
	id := chosenInlineResult.ResultID
	user := chosenInlineResult.From
	split := strings.Split(id, "_")
	if len(split) != 3 {
		bh.logDebug("Wrong ResultID:", id)
		return nil
	}

	cardIndex, err := strconv.Atoi(split[1])
	if err != nil {
		// fail silently
		bh.logDebug(err)
		return nil
	}

	digest, err := strconv.Atoi(split[2])
	if err != nil {
		bh.logDebug(err)
		return nil
	}

	playerGame, found := bh.gameManager.GetPlayerGame(user.ID)
	if !found {
		// fail silently
		return nil
	}
	game, player := playerGame.Game, playerGame.UnoPlayer

	if len(player.Deck().Cards) < cardIndex {
		bh.logDebug("Wrong index: got ", cardIndex, ", array is smaller than given index.")
		return nil
	}
	card := player.Deck().Cards[cardIndex]

	d := int(card.Color)*100 + int(card.CardIndex)

	if d != digest {
		// current card's index does not match previous card
		bh.logDebug("Wrong card color and index: expecting ", d, " got ", digest)
		return nil
	}

	if game.CurrentPlayer().GetUID() != user.ID {
		// fail silently
		return nil
	}

	// player has played a card, so reset their autoskipcount
	player.ResetAutoSkipCount()

	err = game.PlayCard(&card)
	if err != nil {
		_, err := bh.bot.Send(tgbotapi.NewMessage(game.ChatId, err.Error()))
		return err
	}

	player.RemoveCard(uint(cardIndex))
	game.Deck.Discard(card)

	if player.ShouldShoutUNO() {
		bh.bot.Send(tgbotapi.NewMessage(game.ChatId, "UNO!"))
	}

	if player.ShouldChooseColor() {
		msg := tgbotapi.NewMessage(game.ChatId, "Please choose a color")
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
				bh.NewInlineKeyboardButtonSwitchCurrentChat("Choose a color", ""),
			}},
		}
		game.SetTimer(bh)
		_, err := bh.bot.Send(msg)
		return err
	}

	return bh.nextPlayer(game)
}

func (bh *BotHandler) playerWonMsg(chatId int64, player *UnoPlayer) error {
	return bh.SendMessageHTML(chatId, player.HTML()+" won!")
}

func (bh *BotHandler) playerWon(chatId int64, game *UnoGame) error {
	player := game.CurrentPlayer()
	bh.playerWonMsg(chatId, player)
	delete(bh.gameManager.players, player.GetUID())
	err := game.CurrentPlayerWon()
	if err != nil {
		bh.SendMessage(chatId, "Game ended!")
		bh.gameManager.DeleteGame(chatId)
		return err
	}

	return nil
}

// checks if a player won, and announces the next player
func (bh *BotHandler) nextPlayer(game *UnoGame) error {
	// disable the current timer
	game.StopCurrentTimer()
	bh.logDebug("nextPlayer()")

	if game.CurrentPlayer().DidWin() {
		err := bh.playerWon(game.ChatId, game)
		if err != nil {
			// stop the timer as the game is likely not to continue
			game.StopCurrentTimer()
			return err
		}
	} else {
		game.NextPlayer()
	}

	shouldAnnounceNextPlayer, _ := game.SetTimer(bh)
	if shouldAnnounceNextPlayer {
		return bh.announceNextPlayer(game)
	}

	return nil
}

func (bh *BotHandler) announceNextPlayer(game *UnoGame) error {
	player := game.CurrentPlayer()

	_, err := bh.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           game.ChatId,
			ReplyToMessageID: 0,
			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
					bh.NewInlineKeyboardButtonSwitchCurrentChat("Play a card", ""),
				}},
			},
		},
		Text:      "Next player: " + player.HTML(),
		ParseMode: tgbotapi.ModeHTML,
	})

	return err
}

func (bh *BotHandler) announceKickedAFKPlayer(chatId int64, player *UnoPlayer) error {
	return bh.SendMessageHTML(chatId, fmt.Sprintf("%s has been kicked for inactivity.", player.HTML()))
}

func (bh *BotHandler) announcePlayerSkipped(chatId int64, player *UnoPlayer) error {
	return bh.SendMessageHTML(chatId, fmt.Sprintf("%s has been skipped. The skip timer has been reduced to %d seconds.", player.HTML(), uint(player.SkipTimer().Seconds())))
}
