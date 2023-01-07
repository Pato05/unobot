package bot

import (
	"strconv"

	"github.com/Pato05/unobot/uno"
)

type UnoPlayer struct {
	UserId            int64
	Name              string
	deck              uno.PlayerDeck
	shouldChooseColor bool
}

func (self *UnoPlayer) CardCount() int {
	return len(self.deck.Cards)
}

func (self *UnoPlayer) ShouldShoutUNO() bool {
	return self.CardCount() == 1
}

func (self *UnoPlayer) GetUID() int64 {
	return self.UserId
}
func (self *UnoPlayer) Deck() *uno.PlayerDeck {
	return &self.deck
}

func (self *UnoPlayer) DidWin() bool {
	return self.CardCount() == 0 && !self.shouldChooseColor
}

func (self *UnoPlayer) RemoveCard(cardIndex uint) {
	self.deck.Cards = append(self.deck.Cards[:cardIndex], self.deck.Cards[cardIndex+1:]...)
}

func (self *UnoPlayer) HTML() string {
	return "<a href=\"tg://user?id=" + strconv.Itoa(int(self.UserId)) + "\">" + self.Name + "</a>"
}

func (self *UnoPlayer) ShouldChooseColor() bool {
	return self.shouldChooseColor
}

func (self *UnoPlayer) SetShouldChooseColor(shouldChooseColor bool) {
	self.shouldChooseColor = shouldChooseColor
}
