package cards

import "testing"

// Test if the math's correct
func TestCardGlobalIndex(t *testing.T) {
	for index, card := range Cards {
		if int(card.GetGlobalIndex()) != index {
			t.Error("Wrong global index for ", card, ": is ", card.GetGlobalIndex(), " should be ", index, "\n")
		}
	}
}
