package utils

import "math/big"

// BigIntMul multiplies 2 big integers
func BigIntMul(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

// BigIntDiv divides 2 big integers
func BigIntDiv(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Div(x, y)
}

// BigIntAdd adds 2 big integers
func BigIntAdd(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

// BigIntSub subtracts 2 big integers
func BigIntSub(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

// BigIntSqrt returns ⌊√x⌋
func BigIntSqrt(x *big.Int) *big.Int {
	return new(big.Int).Sqrt(x)
}
