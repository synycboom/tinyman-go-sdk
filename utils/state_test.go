package utils_test

import (
	"testing"

	"github.com/synycboom/tinyman-go-sdk/utils"
)

func TestState(t *testing.T) {
	var testConverter = func(val uint64, expected []byte) {
		bb, err := utils.IntToBytes(val)
		if err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			return
		}
		if len(bb) != len(expected) {
			t.Errorf("Converter returns wrong value length %v", len(bb))
			return
		}
		for idx, b := range bb {
			if b != expected[idx] {
				t.Errorf("Converter returns wrong value %v at index %d", b, idx)
			}
		}
	}

	testConverter(0, []byte{0, 0, 0, 0, 0, 0, 0, 0})
	testConverter(127, []byte{0, 0, 0, 0, 0, 0, 0, 127})
	testConverter(128, []byte{0, 0, 0, 0, 0, 0, 0, 128})
	testConverter(255, []byte{0, 0, 0, 0, 0, 0, 0, 255})
	testConverter(512, []byte{0, 0, 0, 0, 0, 0, 2, 0})
}
