package stores

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	RunTests(t, NewMemoryStore())
}
