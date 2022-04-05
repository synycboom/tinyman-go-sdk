package utils

import (
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// StateInt returns unsigned int field from ApplicationLocalState local key-value mapping
func StateInt(state map[string]models.TealValue, key string) uint64 {
	key = base64.StdEncoding.EncodeToString([]byte(key))
	value, ok := state[key]
	if !ok {
		return 0
	}

	return value.Uint
}

// StateBytes returns bytes field from ApplicationLocalState local key-value mapping
func StateBytes(state map[string]models.TealValue, key string) []byte {
	key = base64.StdEncoding.EncodeToString([]byte(key))
	value, ok := state[key]
	if !ok {
		return nil
	}

	return []byte(value.Bytes)
}

// OutstandingAssetStateKey returns an outstanding key of the given asset id
func OutstandingAssetStateKey(assetID uint64) ([]byte, error) {
	assetIDInBytes, err := IntToBytes(assetID)
	if err != nil {
		return nil, err
	}

	var key []byte
	key = append(key, []byte("o")[0])
	key = append(key, assetIDInBytes...)

	return key, nil
}

// ExcessAssetStateKey returns an excess key of the given asset id
func ExcessAssetStateKey(poolAddress string, assetID uint64) ([]byte, error) {
	addr, err := types.DecodeAddress(poolAddress)
	if err != nil {
		return nil, err
	}

	assetIDInBytes, err := IntToBytes(assetID)
	if err != nil {
		return nil, err
	}

	var key []byte
	key = append(key, []byte(addr.String())...)
	key = append(key, []byte("e")[0])
	key = append(key, assetIDInBytes...)

	return key, nil
}
