package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Pato05/unobot/cards"
	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var CardColorEmojis = map[cards.CardColor]string{
	cards.Wild:   "⬛️",
	cards.Blue:   "🟦",
	cards.Green:  "🟩",
	cards.Red:    "🟥",
	cards.Yellow: "🟨",
}

var choosableColorsString = map[cards.CardColor]string{
	cards.Blue:   CardColorEmojis[cards.Blue] + " Blue",
	cards.Green:  CardColorEmojis[cards.Green] + " Green",
	cards.Red:    CardColorEmojis[cards.Red] + " Red",
	cards.Yellow: CardColorEmojis[cards.Yellow] + " Yellow",
}

var CardSpecials = map[cards.CardSpecial]string{
	cards.Special_Colorchooser: "Colorchooser",
	cards.Special_PlusFour:     "+4",
	cards.Special_PlusTwo:      "+2",
	cards.Special_Reverse:      "Reverse",
	cards.Special_Skip:         "Skip",
}

func DrawAction(drawCount uint8) tgbotapi.InlineQueryResultCachedSticker {
	var text string
	if drawCount == 0 {
		text = "Drawing 1 card..."
	} else {
		text = fmt.Sprintf("Drawing %d cards...", drawCount)
	}
	return tgbotapi.InlineQueryResultCachedSticker{
		Type:                "sticker",
		ID:                  "draw_card",
		StickerID:           cards.FileId_DrawIcon,
		InputMessageContent: tgbotapi.InputTextMessageContent{Text: text},
	}
}

func PassAction() tgbotapi.InlineQueryResultCachedSticker {
	return tgbotapi.InlineQueryResultCachedSticker{
		Type:                "sticker",
		ID:                  "pass_turn",
		StickerID:           cards.FileId_PassIcon,
		InputMessageContent: tgbotapi.InputTextMessageContent{Text: "Pass"},
	}
}

func GetGameInfoStr(game *uno.Game[*UnoPlayer]) string {
	currentPlayer := game.CurrentPlayer()
	lastCard := game.PreviousCard
	var cardStr string
	if lastCard.IsSpecial() {
		cardStr = CardSpecials[lastCard.Special]
	} else {
		cardStr = strconv.Itoa(int(lastCard.Number))
	}

	iteratePlayer := func(player *UnoPlayer, currentPlayer *UnoPlayer, isFirst bool) string {
		prefix := " → "
		if isFirst {
			prefix = ""
		}

		content := fmt.Sprintf("%s (%d cards)", player.EscapedName(), player.CardCount())

		if currentPlayer.GetUID() == player.GetUID() {
			return prefix + "<b>" + content + "</b>"
		}
		return prefix + content
	}
	var playersStr strings.Builder

	if game.Reversed {
		length := len(game.Players)
		for i := length - 1; i >= 0; i-- {
			playersStr.WriteString(iteratePlayer(game.Players[i], currentPlayer, i == length-1))
		}
	} else {
		for i, player := range game.Players {
			playersStr.WriteString(iteratePlayer(player, currentPlayer, i == 0))
		}
	}

	return fmt.Sprintf(
		"Current player: <a href=\"tg://user?id=%d\">%s</a>\n"+
			"Last card: %s %s\n"+
			"Players: %s",

		currentPlayer.GetUID(),
		currentPlayer.EscapedName(),

		CardColorEmojis[lastCard.Color],
		cardStr,

		playersStr.String(),
	)

}

func GameInfoAction(game *uno.Game[*UnoPlayer]) tgbotapi.InlineQueryResultCachedSticker {
	return tgbotapi.InlineQueryResultCachedSticker{
		Type:      "sticker",
		ID:        "gameinfo",
		StickerID: cards.FileId_GameInfo,
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text:      GetGameInfoStr(game),
			ParseMode: tgbotapi.ModeHTML,
		},
	}
}

func BluffAction() tgbotapi.InlineQueryResultCachedSticker {
	return tgbotapi.InlineQueryResultCachedSticker{
		Type:                "sticker",
		ID:                  "call_bluff",
		StickerID:           cards.FileId_CallBluff,
		InputMessageContent: tgbotapi.InputTextMessageContent{Text: "I'm calling your bluff!"},
	}
}
