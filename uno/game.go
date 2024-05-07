package uno

import (
	"math/rand"

	"github.com/Pato05/unobot/cards"
)

type Game[T IPlayer] struct {
	// the game deck
	Deck         Deck
	Players      []T
	playersCount uint8
	// the previous thrown card
	PreviousCard *cards.Card
	// the player index that should play
	index uint8
	// if this is true, the player cannot draw
	// and the player must be only able to play the just-drawn card
	DidJustDraw bool
	// if the turn is reversed (reverse card toggles this state)
	Reversed bool
	// how many cards should the next player draw (if more than one, when the player draws, they're automatically skipped)
	DrawCounter uint8
	// if the game has started
	Started bool
	// if new players are NOT allowed to join
	LobbyClosed    bool
	GameCreatorUID int64
	// whether the previous player bluffed
	DidBluff bool
	// whether the next player can call bluff
	CanCallBluff bool
}

func (self *Game[T]) Start() (*cards.Card, error) {
	if self.Started {
		return nil, GameAlreadyStartedError{}
	}
	if len(self.Players) < 2 {
		return nil, TooFewPlayersError{}
	}
	self.Started = true
	self.Deck.Fill()
	self.Deck.Shuffle()
	for _, player := range self.Players {
		player.Deck().Fill(&self.Deck)
	}
	var firstCard *cards.Card
	for {
		card, err := self.Deck.DrawOne()
		if err != nil {
			return nil, err
		}
		firstCard = card

		if firstCard.Color != cards.Wild {
			break
		}
		self.Deck.Discard(*firstCard)
	}
	if err := self.PlayCard(firstCard); err != nil {
		return nil, err
	}
	return firstCard, nil
}

func (self *Game[T]) CloseLobby() error {
	if !self.Started {
		return GameNotStartedError{}
	}

	self.LobbyClosed = true
	return nil
}

func (self *Game[T]) OpenLobby() error {
	if !self.Started {
		return GameNotStartedError{}
	}

	self.LobbyClosed = false
	return nil
}

func (self *Game[T]) JoinPlayer(player T) error {
	if self.LobbyClosed {
		return LobbyClosedError{}
	}

	if len(self.Players) > 10 {
		return TooManyPlayersError{}
	}

	for _, player_ := range self.Players {
		if player.GetUID() == player_.GetUID() {
			return PlayerAlreadyExistsError{}
		}
	}
	if self.Started {
		player.Deck().Fill(&self.Deck)
	}

	self.playersCount++
	self.Players = append(self.Players, player)
	return nil
}

// removes a player from the players list, and returns an error if
// the game should be disbanded
func (self *Game[T]) leavePlayer(index uint8) error {
	// TODO: fix
	if self.playersCount != 0 {
		self.playersCount--
		if self.index == self.playersCount {
			self.index = 0
		}
	}

	self.Players = append(self.Players[:index], self.Players[index+1:]...)
	if self.Started {
		if self.playersCount <= 1 {
			return GameDisbandedLastPlayerWon{}
		}
	} else if self.playersCount == 0 {
		return GameDisbandedNoPlayers{}
	}

	return nil
}

// Warn the GameHandler that the current player won,
// returns an error if the game should be disbanded
func (self *Game[T]) CurrentPlayerWon() error {
	return self.leavePlayer(self.currentPlayerIndex())
}

func (self *Game[T]) getPlayerIndex(player T) (uint8, error) {
	var index uint8 = 0
	found := false
	for i, player_ := range self.Players {
		if player.GetUID() == player_.GetUID() {
			index = uint8(i)
			found = true
			break
		}
	}
	if !found {
		return 0, PlayerNotInGameError{}
	}

	return index, nil
}

func (self *Game[T]) LeavePlayer(player T) error {
	index, err := self.getPlayerIndex(player)
	if err != nil {
		return err
	}

	return self.LeavePlayerByIndex(index)

}

func (self *Game[T]) LeavePlayerByIndex(index uint8) error {
	if self.Started {
		player := self.Players[index]
		self.PreAutoSkipPlayer()
		for _, card := range player.Deck().Cards {
			self.Deck.Discard(card)
		}
	}

	return self.leavePlayer(index)
}

func (self *Game[T]) Reverse() {
	self.Reversed = !self.Reversed
}

func (self *Game[T]) PreviousPlayer() T {
	currentIndex := self.index
	if currentIndex == 0 {
		currentIndex = self.playersCount - 1
	} else {
		currentIndex--
	}

	return self.Players[self.getActualIndex(currentIndex)]
}

func (self *Game[T]) CurrentPlayer() T {
	return self.Players[self.currentPlayerIndex()]
}

func (self *Game[T]) getActualIndex(index uint8) uint8 {
	if self.Reversed {
		return self.playersCount - index - 1
	}

	return index
}

func (self *Game[T]) currentPlayerIndex() uint8 {
	if self.index >= self.playersCount {
		self.index = 0
	}

	return self.getActualIndex(self.index)
}

func (self *Game[T]) NextPlayer() T {
	self.DidJustDraw = false

	self.index++
	return self.Players[self.currentPlayerIndex()]
}

// call this when a player leaves or a player is automatically skipped
func (self *Game[T]) PreAutoSkipPlayer() {
	player := self.CurrentPlayer()
	// if the player should choose the colour (wild card played), let it choose a random colour
	if player.ShouldChooseColor() {
		self.PreviousCard.Color = cards.CardColor(rand.Intn(int(cards.Red)) + int(cards.Blue))
	}
}

// call this after a player drew (one or more) card(s)
func (self *Game[T]) PlayerDrew() {
	self.DrawCounter = 0
	self.DidJustDraw = true
	self.CanCallBluff = false
}

func (self *Game[T]) CurrentPlayerDraw() error {
	var err error
	if self.DrawCounter == 0 {
		err = self.CurrentPlayer().Deck().Draw(&self.Deck, 1)
		self.PlayerDrew()
	} else {
		err = self.CurrentPlayer().Deck().Draw(&self.Deck, self.DrawCounter)
		self.CurrentPlayer().Deck().Sort()
		self.PlayerDrew()
		self.NextPlayer()
	}

	return err
}

// CallBluff action, NextPlayer() must be called afterwards
//
// returns whether the previous player bluffed or not
func (self *Game[T]) CallBluff() bool {
	defer func() {
		self.PlayerDrew()
	}()

	// check if previous player did bluff.
	previousPlayer := self.PreviousPlayer()
	player := self.CurrentPlayer()
	if self.DidBluff {
		previousPlayer.Deck().Draw(&self.Deck, 4)
		return true
	} else {
		player.Deck().Draw(&self.Deck, 6)
		return false
	}

}

func (self *Game[T]) CanCurrentPlayerPlayCard(card *cards.Card) bool {
	if self.DrawCounter != 0 {
		if self.PreviousCard.Special == cards.Special_PlusTwo {
			return card.Special == cards.Special_PlusTwo
		}
		return false
	}

	return self.IsCardPlayable(card)
}

func (self *Game[T]) IsCardPlayable(card *cards.Card) bool {
	if self.PreviousCard == nil {
		return true
	}

	return card.IsPlayable(self.PreviousCard)
}

// Check if current player is bluffing
//
//	-> whether the current player can play any card that isn't a PlusFour
//
// => Used whenever a +4 is played
func (self *Game[T]) IsCurrentPlayerBluffing() bool {
	for _, card := range self.CurrentPlayer().Deck().Cards {
		if card.Special == cards.Special_PlusFour {
			continue
		}

		if self.IsCardPlayable(&card) {
			return true
		}
	}

	return false
}

func (self *Game[T]) PlayCard(card *cards.Card) error {
	if self.DidBluff {
		// reset bluff state
		self.DidBluff = false
	}

	if !self.IsCardPlayable(card) {
		return CardNotPlayableError{}
	}
	switch card.Special {
	case cards.Special_Colorchooser:
		self.CurrentPlayer().SetShouldChooseColor(true)
	case cards.Special_PlusFour:
		self.DidBluff = self.IsCurrentPlayerBluffing()
		self.CanCallBluff = true // for the next player
		self.DrawCounter += 4
		self.CurrentPlayer().SetShouldChooseColor(true)
	case cards.Special_PlusTwo:
		self.DrawCounter += 2
	case cards.Special_Skip:
		self.NextPlayer() // skip the next player
	case cards.Special_Reverse:
		self.Reverse()
	}
	self.PreviousCard = card
	return nil
}

func (self *Game[T]) ChooseColor(color cards.CardColor) {
	if !self.CurrentPlayer().ShouldChooseColor() {
		return
	}

	self.CurrentPlayer().SetShouldChooseColor(false)
	// Choose color by changing the previous card's color;
	// the previous card being a Wild (Colorchooser or PlusFour)
	self.PreviousCard.Color = color
}
