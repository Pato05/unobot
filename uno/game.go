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
	LobbyClosed bool
	// the uid of the player that created the game
	GameCreatorUID int64
	// whether the previous player bluffed
	DidBluff bool
	// whether the next player can call bluff
	CanCallBluff bool
}

func (g *Game[T]) Start() (*cards.Card, error) {
	if g.Started {
		return nil, GameAlreadyStartedError{}
	}
	if len(g.Players) < 2 {
		return nil, TooFewPlayersError{}
	}
	g.Started = true
	g.Deck.Fill()
	g.Deck.Shuffle()
	for _, player := range g.Players {
		player.Deck().Fill(&g.Deck)
	}
	var firstCard *cards.Card
	for {
		card, err := g.Deck.DrawOne()
		if err != nil {
			return nil, err
		}
		firstCard = card

		if firstCard.Color != cards.Wild {
			break
		}
		g.Deck.Discard(*firstCard)
	}
	if err := g.PlayCard(firstCard); err != nil {
		return nil, err
	}
	return firstCard, nil
}

func (g *Game[T]) CloseLobby() error {
	if !g.Started {
		return GameNotStartedError{}
	}

	g.LobbyClosed = true
	return nil
}

func (g *Game[T]) OpenLobby() error {
	if !g.Started {
		return GameNotStartedError{}
	}

	g.LobbyClosed = false
	return nil
}

func (g *Game[T]) JoinPlayer(player T) error {
	if g.LobbyClosed {
		return LobbyClosedError{}
	}

	if len(g.Players) > 10 {
		return TooManyPlayersError{}
	}

	for _, player_ := range g.Players {
		if player.GetUID() == player_.GetUID() {
			return PlayerAlreadyExistsError{}
		}
	}
	if g.Started {
		player.Deck().Fill(&g.Deck)
	}

	g.playersCount++
	g.Players = append(g.Players, player)
	return nil
}

// removes a player from the players list, and returns an error if
// the game should be disbanded
func (g *Game[T]) leavePlayer(index uint8) error {
	// TODO: fix
	if g.playersCount != 0 {
		g.playersCount--
		if g.index == g.playersCount {
			g.index = 0
		}
	}

	g.Players = append(g.Players[:index], g.Players[index+1:]...)
	if g.Started {
		if g.playersCount <= 1 {
			return GameDisbandedLastPlayerWon{}
		}
	} else if g.playersCount == 0 {
		return GameDisbandedNoPlayers{}
	}

	return nil
}

// Warn the GameHandler that the current player won,
// returns an error if the game should be disbanded
func (g *Game[T]) CurrentPlayerWon() error {
	return g.LeaveCurrentPlayer()
}

func (g *Game[T]) getPlayerIndex(player T) (uint8, error) {
	var index uint8 = 0
	found := false
	for i, player_ := range g.Players {
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

func (g *Game[T]) LeavePlayer(player T) error {
	index, err := g.getPlayerIndex(player)
	if err != nil {
		return err
	}

	return g.LeavePlayerByIndex(index)

}

func (g *Game[T]) LeavePlayerByIndex(index uint8) error {
	if g.Started {
		player := g.Players[index]
		g.PreAutoSkipPlayer()
		for _, card := range player.Deck().Cards {
			g.Deck.Discard(card)
		}
	}

	return g.leavePlayer(index)
}

func (g *Game[T]) LeaveCurrentPlayer() error {
	return g.LeavePlayerByIndex(g.currentPlayerIndex())
}

func (g *Game[T]) Reverse() {
	g.Reversed = !g.Reversed
}

func (g *Game[T]) previousIndex() uint8 {
	if g.index == 0 {
		return g.playersCount - 1
	}

	return g.index - 1
}

func (g *Game[T]) nextIndex() uint8 {
	if g.index >= g.playersCount-1 {
		return 0
	}

	return g.index + 1
}

func (g *Game[T]) PreviousPlayer() T {
	var currentIndex uint8
	if g.Reversed {
		currentIndex = g.nextIndex()
	} else {
		currentIndex = g.previousIndex()
	}

	return g.Players[currentIndex]
}

func (g *Game[T]) CurrentPlayer() T {
	return g.Players[g.currentPlayerIndex()]
}

func (g *Game[T]) currentPlayerIndex() uint8 {
	return g.index
}

func (g *Game[T]) NextPlayer() T {
	g.DidJustDraw = false

	// if the player played a reverse card and there's only two players, the turn should be of the current player
	if g.PreviousCard != nil && (g.PreviousCard.Special != cards.Special_Reverse || g.playersCount != 2) {
		if g.Reversed {
			g.index = g.previousIndex()
		} else {
			g.index = g.nextIndex()
		}
	}

	return g.Players[g.currentPlayerIndex()]
}

// call this when a player leaves or a player is automatically skipped
func (g *Game[T]) PreAutoSkipPlayer() {
	player := g.CurrentPlayer()
	// if the player should choose the colour (wild card played), let it choose a random colour
	if player.ShouldChooseColor() {
		newColor := cards.CardColor(rand.Intn(int(cards.Red)) + int(cards.Blue))
		g.ChooseColor(newColor)
	}
}

// call this after a player drew (one or more) card(s)
func (g *Game[T]) PlayerDrew() {
	g.DrawCounter = 0
	g.DidJustDraw = true
	g.CanCallBluff = false
}

func (g *Game[T]) CurrentPlayerDraw() error {
	var err error
	if g.DrawCounter == 0 {
		err = g.CurrentPlayer().Deck().Draw(&g.Deck, 1)
		g.PlayerDrew()
	} else {
		err = g.CurrentPlayer().Deck().Draw(&g.Deck, g.DrawCounter)
		g.CurrentPlayer().Deck().Sort()
		g.PlayerDrew()
		g.NextPlayer()
	}

	return err
}

// CallBluff action, NextPlayer() must be called afterwards
//
// returns whether the previous player bluffed or not
func (g *Game[T]) CallBluff() bool {
	defer func() {
		g.PlayerDrew()
	}()

	// check if previous player did bluff.
	previousPlayer := g.PreviousPlayer()
	player := g.CurrentPlayer()
	if g.DidBluff {
		previousPlayer.Deck().Draw(&g.Deck, 4)
		return true
	} else {
		player.Deck().Draw(&g.Deck, 6)
		return false
	}

}

func (g *Game[T]) CanCurrentPlayerPlayCard(card *cards.Card) bool {
	if g.DrawCounter != 0 {
		if g.PreviousCard.Special == cards.Special_PlusTwo {
			return card.Special == cards.Special_PlusTwo
		}
		return false
	}

	return g.IsCardPlayable(card)
}

func (g *Game[T]) IsCardPlayable(card *cards.Card) bool {
	if g.PreviousCard == nil {
		return true
	}

	return card.IsPlayable(g.PreviousCard)
}

// Check if current player is bluffing
//
//	-> whether the current player can play any card that isn't a PlusFour
//
// => Used whenever a +4 is played
func (g *Game[T]) IsCurrentPlayerBluffing() bool {
	for _, card := range g.CurrentPlayer().Deck().Cards {
		if card.Special == cards.Special_PlusFour {
			continue
		}

		if g.IsCardPlayable(&card) {
			return true
		}
	}

	return false
}

func (g *Game[T]) PlayCard(card *cards.Card) error {
	if g.DidBluff {
		// reset bluff state
		g.DidBluff = false
	}

	if !g.IsCardPlayable(card) {
		return CardNotPlayableError{}
	}
	switch card.Special {
	case cards.Special_Colorchooser:
		g.CurrentPlayer().SetShouldChooseColor(true)
	case cards.Special_PlusFour:
		g.DidBluff = g.IsCurrentPlayerBluffing()
		g.CanCallBluff = true // for the next player
		g.DrawCounter += 4
		g.CurrentPlayer().SetShouldChooseColor(true)
	case cards.Special_PlusTwo:
		g.DrawCounter += 2
	case cards.Special_Skip:
		g.NextPlayer() // skip the next player
	case cards.Special_Reverse:
		g.Reverse()
	}
	g.PreviousCard = card
	return nil
}

func (g *Game[T]) ChooseColor(color cards.CardColor) {
	player := g.CurrentPlayer()
	if !player.ShouldChooseColor() {
		return
	}

	player.SetShouldChooseColor(false)
	// Choose color by changing the previous card's color;
	// the previous card being a Wild (Colorchooser or PlusFour)
	g.PreviousCard.Color = color
}
