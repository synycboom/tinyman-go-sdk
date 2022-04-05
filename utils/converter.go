package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

// IntToBytes convert int into 8-bit bytes
func IntToBytes(n uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, n); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToBigUint converts an unsigned 64-bit integer to big integer
func ToBigUint(v uint64) *big.Int {
	return new(big.Int).SetUint64(v)
}

// ToBigFloat converts an unsigned 64-bit integer to big float
func ToBigFloat(v uint64) *big.Float {
	return new(big.Float).SetUint64(v)
}
