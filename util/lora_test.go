package util

import (
	"testing"
)

func TestCalculatePacketTime(t *testing.T) {
	if time, err := CalculatePacketTime(23, "SF7BW125"); time != 56.576 {
		t.Errorf("expected 56.576, got %f with error %v", time, err)
	}
	// TODO: one example isn't enough
}
