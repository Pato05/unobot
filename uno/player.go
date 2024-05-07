package uno

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
	Name              string
	deck              PlayerDeck
	shouldChooseColor bool
}

func (self *Player) ShouldChooseColor() bool {
	return self.shouldChooseColor
}

func (self *Player) SetShouldChooseColor(shouldChooseColor bool) {
	self.shouldChooseColor = shouldChooseColor
}

func (self *Player) CardCount() int {
	return len(self.deck.Cards)
}

func (self *Player) ShouldShoutUNO() bool {
	return self.CardCount() == 1
}

func (self *Player) GetUID() int64 {
	return self.Id
}

func (self *Player) DidWin() bool {
	return self.CardCount() == 0 && !self.shouldChooseColor
}

func (self *Player) Deck() *PlayerDeck {
	return &self.deck
}

func (self *Player) RemoveCard(cardIndex uint) {
	self.deck.Cards = append(self.deck.Cards[:cardIndex], self.deck.Cards[cardIndex+1:]...)
}
