package uno

import (
	"reflect"
	"testing"

	"github.com/Pato05/unobot/cards"
)

func BenchmarkDeckShuffle(b *testing.B) {
	var deck Deck
	deck.Fill()
	for i := 0; i < b.N; i++ {
		deck.Shuffle()
	}
}

func BenchmarkDeckDraw(b *testing.B) {
	var deck Deck
	deck.Fill()
	for i := 0; i < b.N; i++ {
		deck.DrawOne()
	}
}

func BenchmarkDeckFill(b *testing.B) {
	var deck Deck
	for i := 0; i < b.N; i++ {
		deck.Fill()
	}
}

func BenchmarkPlayerDeckSortOnce(b *testing.B) {
	var deck Deck
	deck.Fill()
	var pDeck PlayerDeck
	pDeck.Fill(&deck)

	for i := 0; i < b.N; i++ {
		pDeck.Sort()
	}
}

func BenchmarkPlayerDeckSortEvery(b *testing.B) {
	var deck Deck
	deck.Fill()
	Cards := deck.Cards
	var pDeck PlayerDeck
	pDeck.Fill(&deck)

	for i := 0; i < b.N; i++ {
		pDeck.Cards = []cards.Card{}
		deck.Cards = Cards
		pDeck.Fill(&deck)
		pDeck.Sort()
	}
}

func TestDeckCards(t *testing.T) {
	var deck Deck
	deck.Fill()

	if len(deck.Cards) != 108 {
		t.Errorf("Deck cards should've been 108, but they're %d", len(deck.Cards))
	}

	specialsCount := map[cards.CardSpecial]uint32{}
	// 4 Wild and 4 Wild +4; then two special cards for each color
	shouldSpecialsBe := map[cards.CardSpecial]uint32{
		cards.Special_Colorchooser: 4,
		cards.Special_PlusFour:     4,
		cards.Special_PlusTwo:      8,
		cards.Special_Reverse:      8,
		cards.Special_Skip:         8,
	}

	numbersCount := map[uint8]uint32{}
	// 2 for each color, except 0 which is present only once
	shouldNumbersCountBe := map[uint8]uint32{
		0: 4,
		1: 8,
		2: 8,
		3: 8,
		4: 8,
		5: 8,
		6: 8,
		7: 8,
		8: 8,
		9: 8,
	}

	for _, card := range deck.Cards {
		if card.IsSpecial() {
			specialsCount[card.Special]++
		} else {
			numbersCount[uint8(card.Number)]++
		}
	}

	if !reflect.DeepEqual(specialsCount, shouldSpecialsBe) {
		t.Error("Special cards are not distributed in the correct way\n", specialsCount, "\nShould be:\n", shouldSpecialsBe)
	}

	if !reflect.DeepEqual(numbersCount, shouldNumbersCountBe) {
		t.Error("Numeral cards are not distributed in the correct way\n", numbersCount, "\nShould be:\n", shouldNumbersCountBe)
	}
}
