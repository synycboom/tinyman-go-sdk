package utils_test

import (
	"testing"

	tUtils "github.com/synycboom/tinyman-go-sdk/utils"
)

func TestEncodeValue(t *testing.T) {
	_, err := tUtils.EncodeValue(1, "unknown")
	if err == nil {
		t.Errorf("EncodeValue should support only int type")
	}

	_, err = tUtils.EncodeValue(1, "int")
	if err == nil {
		t.Errorf("EncodeValue should not support untyped value")
	}

	_, err = tUtils.EncodeValue(uint64(1), "int")
	if err != nil {
		t.Errorf("EncodeValue should support uint64 value")
	}
}

func TestEncodeVarInt(t *testing.T) {
	var testEncodeVarInt = func(val uint64, expected []byte) {
		bb := tUtils.EncodeVarInt(val)
		if len(bb) != len(expected) {
			t.Errorf("EncodeValue returns wrong value length %v", len(bb))

			return
		}
		for idx, b := range bb {
			if b != expected[idx] {
				t.Errorf("EncodeValue returns wrong value %v at index %d", b, idx)
			}
		}
	}

	testEncodeVarInt(0, []byte{0x00})
	testEncodeVarInt(127, []byte{0x7F})
	testEncodeVarInt(128, []byte{0x80, 0x01})
	testEncodeVarInt(255, []byte{0xFF, 0x01})
	testEncodeVarInt(256, []byte{0x80, 0x02})
}
