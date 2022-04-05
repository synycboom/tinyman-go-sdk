package utils

import "math/big"

// BigFloatDiv divides 2 big floats
func BigFloatDiv(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).Quo(x, y)
}
