package cards

type CardColor uint8
type CardSpecial uint8

// colors
const (
	Wild CardColor = iota
	Blue
	Yellow
	Green
	Red
)

// specials
const (
	Special_None CardSpecial = iota
	Special_Colorchooser
	Special_PlusFour

	// these are also used as card indexes
	Special_PlusTwo CardSpecial = 10
	Special_Reverse CardSpecial = 11
	Special_Skip    CardSpecial = 12
)
