package utility

import (
	"testing"
)

func TestHack_value(t *testing.T) {
	var hack Hack = ""
	t.Log(hack.value())
}
