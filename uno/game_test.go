package uno

import (
	"testing"

	"github.com/Pato05/unobot/cards"
)

func TestPreAutoSkip(t *testing.T) {
	var game Game[*Player]
	player := &Player{
		Id: 1,
	}
	game.Players = append(game.Players, player)
	player.SetShouldChooseColor(true)

	// set prev card to a +4
	game.PreviousCard = &cards.Cards[1]

	if game.PreviousCard.Color != cards.Wild {
		t.Error("card colour isn't Wild.")
	}

	for i := 0; i < 1000; i++ {
		game.PreAutoSkipPlayer()

		t.Log("color: ", game.PreviousCard.Color)

		if game.PreviousCard.Color < cards.Blue || game.PreviousCard.Color > cards.Red {
			t.Error("card colour isn't 1 <= x <= 4, value is ", game.PreviousCard.Color)
		}
	}
}
