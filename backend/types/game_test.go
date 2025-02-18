package types

import (
	"testing"
)

func TestDrawCard(t *testing.T) {

	t.Run("DrawOneCard", func(t *testing.T) {
		cards := *NewCards()
		target := &Cards{}


		cardsLength := len(cards)
		
		cards.DrawCard(target)

		if len(cards) != cardsLength - 1 {
			t.Errorf("Expected %d cards in cards, got %d", cardsLength - 1, len(cards))
		}
		if len(*target) != 1 {
			t.Errorf("Expected 2 cards in target, got %d", len(*target))
		}
	})

	t.Run("DrawAllCards", func(t *testing.T) {
		cards := *NewCards()
		target := &Cards{}


		cardsLength := len(cards)
		for i := 0; i < cardsLength; i++ {
			cards.DrawCard(target)
		}
		if len(cards) != 0 {
			t.Errorf("Expected 0 cards in cards, got %d", len(cards))
		}
		if len(*target) != cardsLength {
			t.Errorf("Expected %d cards in target, got %d", cardsLength, len(*target))
		}
	})

	t.Run("DrawNoCards", func(t *testing.T) {
		cards := &Cards{}
		target := &Cards{}

		if cards.DrawCard(target) {
			t.Errorf("Expected false, got true")
		}
	})
}
