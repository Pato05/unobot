package cards

type Card struct {
	Color     CardColor
	Special   CardSpecial
	Number    int8
	CardIndex uint8
}

type CardIDs struct {
	Normal string
	Grey   string
}

func (self *Card) IsSpecial() bool {
	return self.Special != Special_None
}

// unless self.Color is Black and it is the last card available in the deck
func (self *Card) IsPlayable(prevCard *Card) bool {
	if self.Color == Wild {
		return true
	}

	if prevCard.Special == Special_PlusFour || prevCard.Special == Special_Colorchooser {
		return self.Color == prevCard.Color
	}

	return self.Color == prevCard.Color || self.CardIndex == prevCard.CardIndex
}

func (self *Card) GetFileID() CardIDs {
	return CardFileIDsByColor[self.Color][self.CardIndex]
}

// returns index for the "Cards" map
func (self *Card) GetGlobalIndex() uint16 {
	globIndex := uint16(self.CardIndex)
	if self.Color != Wild {
		globIndex = globIndex + (13 * (uint16(self.Color) - 1)) + 2
	}
	return globIndex
}
