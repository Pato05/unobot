package uno

import (
	"math/rand"
	"time"

	"github.com/Pato05/unobot/cards"
)

const INITIAL_CARDS_COUNT = 7

// thrown if the deck is empty and there are no (more) cards in the discard pile
type DeckEmptyError struct{}

func (self DeckEmptyError) Error() string {
	return "No more cards in deck."
}

type PlayerDeck struct {
	Cards []cards.Card
}

func (self *PlayerDeck) Fill(deck *Deck) {
	self.Draw(deck, INITIAL_CARDS_COUNT)
}

func (self *PlayerDeck) Draw(deck *Deck, n uint8) error {
	cards, err := deck.DrawMulti(n)
	if err != nil {
		return err
	}
	self.Cards = append(self.Cards, cards...)
	return nil
}

type Deck struct {
	Cards     []cards.Card
	Discarded []cards.Card
}

func (self *Deck) Fill() {
	// 108 cards total
	self.Cards = make([]cards.Card, 0)

	// 4 cards of each Wild
	self.Cards = append(
		self.Cards,
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
	self.Cards = append(self.Cards, cards.Cards[2:15]...)
	self.Cards = append(self.Cards, cards.Cards[3:15]...)
	// Yellow
	self.Cards = append(self.Cards, cards.Cards[15:28]...)
	self.Cards = append(self.Cards, cards.Cards[16:28]...)
	// Green
	self.Cards = append(self.Cards, cards.Cards[28:41]...)
	self.Cards = append(self.Cards, cards.Cards[29:41]...)
	// Red
	self.Cards = append(self.Cards, cards.Cards[41:54]...)
	self.Cards = append(self.Cards, cards.Cards[42:54]...)
}

func (self *Deck) DrawOne() (*cards.Card, error) {
	cards_, err := self.DrawMulti(1)
	if err != nil {
		return nil, err
	}

	card := cards_[0]
	return &card, nil
}

func (self *Deck) Shuffle() {
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(self.Cards), func(i, j int) {
		self.Cards[i], self.Cards[j] = self.Cards[j], self.Cards[i]
	})
}

func (self *Deck) FillFromDiscarded() error {
	if len(self.Discarded) == 0 {
		return DeckEmptyError{}
	}
	self.Cards = append(self.Cards, self.Discarded...)
	self.Discarded = []cards.Card{}
	self.Shuffle()
	return nil
}

func (self *Deck) DrawMulti(n uint8) ([]cards.Card, error) {
	if int(n) > len(self.Cards) {
		n -= uint8(len(self.Cards))
		if err := self.FillFromDiscarded(); err != nil {
			return nil, err
		}
		if n > 0 {
			return self.DrawMulti(n)
		}
	}
	cards := self.Cards[0:n]
	self.Cards = self.Cards[n:]
	return cards, nil
}

func (self *Deck) Discard(cards ...cards.Card) {
	self.Discarded = append(self.Discarded, cards...)
}
