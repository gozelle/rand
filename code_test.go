package rand_test

import (
	"github.com/gozelle/rand"
	"testing"
)

func TestCode(t *testing.T) {
	t.Log(rand.Code(6))
}
