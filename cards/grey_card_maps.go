package cards

// generated by card_gen.py, based on cards/grey_card_constants.go

var GreyCardFileIDs = []string{
	// Green cards
	fileId_Grey_Green_Zero, fileId_Grey_Green_One, fileId_Grey_Green_Two, fileId_Grey_Green_Three, fileId_Grey_Green_Four, fileId_Grey_Green_Five, fileId_Grey_Green_Six, fileId_Grey_Green_Seven, fileId_Grey_Green_Eight, fileId_Grey_Green_Nine, fileId_Grey_Green_PlusTwo, fileId_Grey_Green_Reverse, fileId_Grey_Green_Skip,

	// Red cards
	fileId_Grey_Red_Zero, fileId_Grey_Red_One, fileId_Grey_Red_Two, fileId_Grey_Red_Three, fileId_Grey_Red_Four, fileId_Grey_Red_Five, fileId_Grey_Red_Six, fileId_Grey_Red_Seven, fileId_Grey_Red_Eight, fileId_Grey_Red_Nine, fileId_Grey_Red_PlusTwo, fileId_Grey_Red_Reverse, fileId_Grey_Red_Skip,

	// Yellow cards
	fileId_Grey_Yellow_Zero, fileId_Grey_Yellow_One, fileId_Grey_Yellow_Two, fileId_Grey_Yellow_Three, fileId_Grey_Yellow_Four, fileId_Grey_Yellow_Five, fileId_Grey_Yellow_Six, fileId_Grey_Yellow_Seven, fileId_Grey_Yellow_Eight, fileId_Grey_Yellow_Nine, fileId_Grey_Yellow_PlusTwo, fileId_Grey_Yellow_Reverse, fileId_Grey_Yellow_Skip,

	// Blue cards
	fileId_Grey_Blue_Zero, fileId_Grey_Blue_One, fileId_Grey_Blue_Two, fileId_Grey_Blue_Three, fileId_Grey_Blue_Four, fileId_Grey_Blue_Five, fileId_Grey_Blue_Six, fileId_Grey_Blue_Seven, fileId_Grey_Blue_Eight, fileId_Grey_Blue_Nine, fileId_Grey_Blue_PlusTwo, fileId_Grey_Blue_Reverse, fileId_Grey_Blue_Skip,

	// Black cards
	fileId_Grey_Colorchooser, fileId_Grey_PlusFour,
}

// color -> cards file ids lookup table
var GreyCardFileIDsByColor = map[CardColor][]string{
	Green:  GreyCardFileIDs[0:13],
	Red:    GreyCardFileIDs[13:26],
	Yellow: GreyCardFileIDs[26:39],
	Blue:   GreyCardFileIDs[39:52],
	Wild:   GreyCardFileIDs[52:54],
}
