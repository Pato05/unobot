package uno

type Player interface {
	GetUID() int64
	Deck() *PlayerDeck
	ShouldChooseColor() bool
	SetShouldChooseColor(bool)
	DidWin() bool
}
