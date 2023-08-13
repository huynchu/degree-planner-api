package utils

import "testing"

func TestSliceUtil__Move(t *testing.T) {
	original := []string{"a", "b", "c", "d", "e"}
	slice := []string{"a", "b", "c", "d", "e"}
	sliceMoved := Move(slice, 1, 3)

	expected := []string{"a", "c", "d", "b", "e"}
	if !equals(original, slice) {
		t.Errorf("Original changed to %v", slice)
	}

	if !equals(sliceMoved, expected) {
		t.Errorf("Expected %v, got %v", expected, slice)
	}
}

func equals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
