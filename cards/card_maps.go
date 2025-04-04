package cards

// generated by gen/card_gen.py, based on cards/card_constants.go

// all the uno cards
var Cards = []Card{
	// Wild cards
	{Color: Wild, Special: Special_Colorchooser, Number: -1, CardIndex: 0},
	{Color: Wild, Special: Special_PlusFour, Number: -1, CardIndex: 1},

	// Blue cards
	{Color: Blue, Special: Special_None, Number: 0, CardIndex: 0},
	{Color: Blue, Special: Special_None, Number: 1, CardIndex: 1},
	{Color: Blue, Special: Special_None, Number: 2, CardIndex: 2},
	{Color: Blue, Special: Special_None, Number: 3, CardIndex: 3},
	{Color: Blue, Special: Special_None, Number: 4, CardIndex: 4},
	{Color: Blue, Special: Special_None, Number: 5, CardIndex: 5},
	{Color: Blue, Special: Special_None, Number: 6, CardIndex: 6},
	{Color: Blue, Special: Special_None, Number: 7, CardIndex: 7},
	{Color: Blue, Special: Special_None, Number: 8, CardIndex: 8},
	{Color: Blue, Special: Special_None, Number: 9, CardIndex: 9},
	{Color: Blue, Special: Special_PlusTwo, Number: -1, CardIndex: uint8(Special_PlusTwo)},
	{Color: Blue, Special: Special_Reverse, Number: -1, CardIndex: uint8(Special_Reverse)},
	{Color: Blue, Special: Special_Skip, Number: -1, CardIndex: uint8(Special_Skip)},

	// Yellow cards
	{Color: Yellow, Special: Special_None, Number: 0, CardIndex: 0},
	{Color: Yellow, Special: Special_None, Number: 1, CardIndex: 1},
	{Color: Yellow, Special: Special_None, Number: 2, CardIndex: 2},
	{Color: Yellow, Special: Special_None, Number: 3, CardIndex: 3},
	{Color: Yellow, Special: Special_None, Number: 4, CardIndex: 4},
	{Color: Yellow, Special: Special_None, Number: 5, CardIndex: 5},
	{Color: Yellow, Special: Special_None, Number: 6, CardIndex: 6},
	{Color: Yellow, Special: Special_None, Number: 7, CardIndex: 7},
	{Color: Yellow, Special: Special_None, Number: 8, CardIndex: 8},
	{Color: Yellow, Special: Special_None, Number: 9, CardIndex: 9},
	{Color: Yellow, Special: Special_PlusTwo, Number: -1, CardIndex: uint8(Special_PlusTwo)},
	{Color: Yellow, Special: Special_Reverse, Number: -1, CardIndex: uint8(Special_Reverse)},
	{Color: Yellow, Special: Special_Skip, Number: -1, CardIndex: uint8(Special_Skip)},

	// Green cards
	{Color: Green, Special: Special_None, Number: 0, CardIndex: 0},
	{Color: Green, Special: Special_None, Number: 1, CardIndex: 1},
	{Color: Green, Special: Special_None, Number: 2, CardIndex: 2},
	{Color: Green, Special: Special_None, Number: 3, CardIndex: 3},
	{Color: Green, Special: Special_None, Number: 4, CardIndex: 4},
	{Color: Green, Special: Special_None, Number: 5, CardIndex: 5},
	{Color: Green, Special: Special_None, Number: 6, CardIndex: 6},
	{Color: Green, Special: Special_None, Number: 7, CardIndex: 7},
	{Color: Green, Special: Special_None, Number: 8, CardIndex: 8},
	{Color: Green, Special: Special_None, Number: 9, CardIndex: 9},
	{Color: Green, Special: Special_PlusTwo, Number: -1, CardIndex: uint8(Special_PlusTwo)},
	{Color: Green, Special: Special_Reverse, Number: -1, CardIndex: uint8(Special_Reverse)},
	{Color: Green, Special: Special_Skip, Number: -1, CardIndex: uint8(Special_Skip)},

	// Red cards
	{Color: Red, Special: Special_None, Number: 0, CardIndex: 0},
	{Color: Red, Special: Special_None, Number: 1, CardIndex: 1},
	{Color: Red, Special: Special_None, Number: 2, CardIndex: 2},
	{Color: Red, Special: Special_None, Number: 3, CardIndex: 3},
	{Color: Red, Special: Special_None, Number: 4, CardIndex: 4},
	{Color: Red, Special: Special_None, Number: 5, CardIndex: 5},
	{Color: Red, Special: Special_None, Number: 6, CardIndex: 6},
	{Color: Red, Special: Special_None, Number: 7, CardIndex: 7},
	{Color: Red, Special: Special_None, Number: 8, CardIndex: 8},
	{Color: Red, Special: Special_None, Number: 9, CardIndex: 9},
	{Color: Red, Special: Special_PlusTwo, Number: -1, CardIndex: uint8(Special_PlusTwo)},
	{Color: Red, Special: Special_Reverse, Number: -1, CardIndex: uint8(Special_Reverse)},
	{Color: Red, Special: Special_Skip, Number: -1, CardIndex: uint8(Special_Skip)},
}

// all cards' file id
var CardFileIDs = []CardIDs{
	// Wild cards
	{Normal: fileId_Colorchooser, Gray: fileId_Gray_Colorchooser},
	{Normal: fileId_PlusFour, Gray: fileId_Gray_PlusFour},

	// Blue cards
	{Normal: fileId_Blue_Zero, Gray: fileId_Gray_Blue_Zero},
	{Normal: fileId_Blue_One, Gray: fileId_Gray_Blue_One},
	{Normal: fileId_Blue_Two, Gray: fileId_Gray_Blue_Two},
	{Normal: fileId_Blue_Three, Gray: fileId_Gray_Blue_Three},
	{Normal: fileId_Blue_Four, Gray: fileId_Gray_Blue_Four},
	{Normal: fileId_Blue_Five, Gray: fileId_Gray_Blue_Five},
	{Normal: fileId_Blue_Six, Gray: fileId_Gray_Blue_Six},
	{Normal: fileId_Blue_Seven, Gray: fileId_Gray_Blue_Seven},
	{Normal: fileId_Blue_Eight, Gray: fileId_Gray_Blue_Eight},
	{Normal: fileId_Blue_Nine, Gray: fileId_Gray_Blue_Nine},
	{Normal: fileId_Blue_PlusTwo, Gray: fileId_Gray_Blue_PlusTwo},
	{Normal: fileId_Blue_Reverse, Gray: fileId_Gray_Blue_Reverse},
	{Normal: fileId_Blue_Skip, Gray: fileId_Gray_Blue_Skip},

	// Yellow cards
	{Normal: fileId_Yellow_Zero, Gray: fileId_Gray_Yellow_Zero},
	{Normal: fileId_Yellow_One, Gray: fileId_Gray_Yellow_One},
	{Normal: fileId_Yellow_Two, Gray: fileId_Gray_Yellow_Two},
	{Normal: fileId_Yellow_Three, Gray: fileId_Gray_Yellow_Three},
	{Normal: fileId_Yellow_Four, Gray: fileId_Gray_Yellow_Four},
	{Normal: fileId_Yellow_Five, Gray: fileId_Gray_Yellow_Five},
	{Normal: fileId_Yellow_Six, Gray: fileId_Gray_Yellow_Six},
	{Normal: fileId_Yellow_Seven, Gray: fileId_Gray_Yellow_Seven},
	{Normal: fileId_Yellow_Eight, Gray: fileId_Gray_Yellow_Eight},
	{Normal: fileId_Yellow_Nine, Gray: fileId_Gray_Yellow_Nine},
	{Normal: fileId_Yellow_PlusTwo, Gray: fileId_Gray_Yellow_PlusTwo},
	{Normal: fileId_Yellow_Reverse, Gray: fileId_Gray_Yellow_Reverse},
	{Normal: fileId_Yellow_Skip, Gray: fileId_Gray_Yellow_Skip},

	// Green cards
	{Normal: fileId_Green_Zero, Gray: fileId_Gray_Green_Zero},
	{Normal: fileId_Green_One, Gray: fileId_Gray_Green_One},
	{Normal: fileId_Green_Two, Gray: fileId_Gray_Green_Two},
	{Normal: fileId_Green_Three, Gray: fileId_Gray_Green_Three},
	{Normal: fileId_Green_Four, Gray: fileId_Gray_Green_Four},
	{Normal: fileId_Green_Five, Gray: fileId_Gray_Green_Five},
	{Normal: fileId_Green_Six, Gray: fileId_Gray_Green_Six},
	{Normal: fileId_Green_Seven, Gray: fileId_Gray_Green_Seven},
	{Normal: fileId_Green_Eight, Gray: fileId_Gray_Green_Eight},
	{Normal: fileId_Green_Nine, Gray: fileId_Gray_Green_Nine},
	{Normal: fileId_Green_PlusTwo, Gray: fileId_Gray_Green_PlusTwo},
	{Normal: fileId_Green_Reverse, Gray: fileId_Gray_Green_Reverse},
	{Normal: fileId_Green_Skip, Gray: fileId_Gray_Green_Skip},

	// Red cards
	{Normal: fileId_Red_Zero, Gray: fileId_Gray_Red_Zero},
	{Normal: fileId_Red_One, Gray: fileId_Gray_Red_One},
	{Normal: fileId_Red_Two, Gray: fileId_Gray_Red_Two},
	{Normal: fileId_Red_Three, Gray: fileId_Gray_Red_Three},
	{Normal: fileId_Red_Four, Gray: fileId_Gray_Red_Four},
	{Normal: fileId_Red_Five, Gray: fileId_Gray_Red_Five},
	{Normal: fileId_Red_Six, Gray: fileId_Gray_Red_Six},
	{Normal: fileId_Red_Seven, Gray: fileId_Gray_Red_Seven},
	{Normal: fileId_Red_Eight, Gray: fileId_Gray_Red_Eight},
	{Normal: fileId_Red_Nine, Gray: fileId_Gray_Red_Nine},
	{Normal: fileId_Red_PlusTwo, Gray: fileId_Gray_Red_PlusTwo},
	{Normal: fileId_Red_Reverse, Gray: fileId_Gray_Red_Reverse},
	{Normal: fileId_Red_Skip, Gray: fileId_Gray_Red_Skip},
}

// color -> cards file ids lookup table
var CardFileIDsByColor = map[CardColor][]CardIDs{
	Wild:   CardFileIDs[0:2],
	Blue:   CardFileIDs[2:15],
	Yellow: CardFileIDs[15:28],
	Green:  CardFileIDs[28:41],
	Red:    CardFileIDs[41:54],
}
