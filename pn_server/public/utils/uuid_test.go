package utils

import "testing"

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(UUID())
	}
}
