package main

import "testing"

func TestAddReturnsCorrectSum(t *testing.T) {
	if result := add(2, 3); result != 5 {
		t.Errorf("expected 5 got %d", result)
	}
}
