package uno

import (
	"math/rand"
	"slices"

	"github.com/Pato05/unobot/cards"
)

const INITIAL_CARDS_COUNT = 7

// thrown if the deck is empty and there are no (more) cards in the discard pile
type DeckEmptyError struct{}

func (e DeckEmptyError) Error() string {
	return "No more cards in deck."
}

type PlayerDeck struct {
	Cards []cards.Card
}

func (pd *PlayerDeck) Sort() {
	slices.SortStableFunc(pd.Cards, func(a, b cards.Card) int {
		return int(b.GetGlobalIndex()) - int(a.GetGlobalIndex())
	})
}

func (pd *PlayerDeck) Fill(deck *Deck) {
	pd.Draw(deck, INITIAL_CARDS_COUNT)
	pd.Sort()
}

func (pd *PlayerDeck) Draw(deck *Deck, n uint8) error {
	newCards, err := deck.DrawMulti(n)
	if err != nil {
		return err
	}

	// put the new cards on the "top" of the deck
	newSlice := make([]cards.Card, int(n)+len(pd.Cards))
	copy(newSlice[:n], newCards)
	copy(newSlice[n:], pd.Cards)
	pd.Cards = newSlice

	return nil
}

type Deck struct {
	Cards     []cards.Card
	Discarded []cards.Card
}

func (d *Deck) Fill() {
	// 108 cards total
	d.Cards = make([]cards.Card, 0)

	// 4 cards of each Wild
	d.Cards = append(
		d.Cards,
		cards.Cards[0],
		cards.Cards[0],
		cards.Cards[0],
		cards.Cards[0],
		cards.Cards[1],
		cards.Cards[1],
		cards.Cards[1],
		cards.Cards[1],
	)

	// two of each color, except the 0, which is only present once
	// Blue
	d.Cards = append(d.Cards, cards.Cards[2:15]...)
	d.Cards = append(d.Cards, cards.Cards[3:15]...)
	// Yellow
	d.Cards = append(d.Cards, cards.Cards[15:28]...)
	d.Cards = append(d.Cards, cards.Cards[16:28]...)
	// Green
	d.Cards = append(d.Cards, cards.Cards[28:41]...)
	d.Cards = append(d.Cards, cards.Cards[29:41]...)
	// Red
	d.Cards = append(d.Cards, cards.Cards[41:54]...)
	d.Cards = append(d.Cards, cards.Cards[42:54]...)
}

func (d *Deck) DrawOne() (*cards.Card, error) {
	cards_, err := d.DrawMulti(1)
	if err != nil {
		return nil, err
	}

	card := cards_[0]
	return &card, nil
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) FillFromDiscarded() error {
	if len(d.Discarded) == 0 {
		return DeckEmptyError{}
	}
	d.Cards = append(d.Cards, d.Discarded...)
	d.Discarded = []cards.Card{}
	d.Shuffle()
	return nil
}

func (d *Deck) DrawMulti(n uint8) ([]cards.Card, error) {
	if int(n) > len(d.Cards) {
		n -= uint8(len(d.Cards))
		if err := d.FillFromDiscarded(); err != nil {
			return nil, err
		}
		if n > 0 {
			return d.DrawMulti(n)
		}
	}
	cards := d.Cards[0:n]
	d.Cards = d.Cards[n:]
	return cards, nil
}

func (d *Deck) Discard(cards ...cards.Card) {
	d.Discarded = append(d.Discarded, cards...)
}
