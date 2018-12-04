package storage

import "testing"

func TestMemoryStore(t *testing.T) {
	testStore(t, NewMemoryStore())
}
