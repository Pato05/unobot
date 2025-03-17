package uno

import (
	"time"

	"github.com/Pato05/unobot/constants"
)

type IPlayer interface {
	ShouldChooseColor() bool
	CardCount() int
	ShouldShoutUNO() bool
	GetUID() int64
	DidWin() bool
	Deck() *PlayerDeck
	SetShouldChooseColor(shouldChooseColor bool)
}

type Player struct {
	Id                int64
	deck              PlayerDeck
	shouldChooseColor bool
	// How many times this player has been automatically skipped
	AutoSkipCount int
}

func (p *Player) ShouldChooseColor() bool {
	return p.shouldChooseColor
}

func (p *Player) SetShouldChooseColor(shouldChooseColor bool) {
	p.shouldChooseColor = shouldChooseColor
}

func (p *Player) CardCount() int {
	return len(p.deck.Cards)
}

func (p *Player) ShouldShoutUNO() bool {
	return p.CardCount() == 1
}

func (p *Player) GetUID() int64 {
	return p.Id
}

func (p *Player) DidWin() bool {
	return p.CardCount() == 0 && !p.shouldChooseColor
}

func (p *Player) Deck() *PlayerDeck {
	return &p.deck
}

func (p *Player) RemoveCard(cardIndex uint) {
	p.deck.Cards = append(p.deck.Cards[:cardIndex], p.deck.Cards[cardIndex+1:]...)
}

func (p *Player) SkipTimer() time.Duration {
	return time.Duration(constants.DEFAULT_SKIP_TIMER-(p.AutoSkipCount*15)) * time.Second
}

func (p *Player) ResetAutoSkipCount() {
	p.AutoSkipCount = 0
}

func (p *Player) IncreaseAutoSkipCount() {
	p.AutoSkipCount++
}
