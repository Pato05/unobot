package cards

type Card struct {
	Color     CardColor
	Special   CardSpecial
	Number    int8
	CardIndex uint8
}

type CardIDs struct {
	Normal string
	Gray   string
}

func (c *Card) IsSpecial() bool {
	return c.Special != Special_None
}

// unless self.Color is Black and it is the last card available in the deck
func (c *Card) IsPlayable(prevCard *Card) bool {
	if c.Color == Wild {
		return true
	}

	if prevCard.Special == Special_PlusFour || prevCard.Special == Special_Colorchooser {
		return c.Color == prevCard.Color
	}

	return c.Color == prevCard.Color || c.CardIndex == prevCard.CardIndex
}

func (c *Card) GetFileID() CardIDs {
	return CardFileIDsByColor[c.Color][c.CardIndex]
}

// returns index for the "Cards" map
func (c *Card) GetGlobalIndex() uint16 {
	globIndex := uint16(c.CardIndex)
	if c.Color != Wild {
		globIndex = globIndex + (13 * (uint16(c.Color) - 1)) + 2
	}
	return globIndex
}
