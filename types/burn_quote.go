package types

// BurnQuote represents a burn quote
type BurnQuote struct {
	// AmountsOut is an asset mapping which maps between asset ids and output asset amounts
	AmountsOut map[uint64]AssetAmount

	// LiquidityAssetAmount is a liquidity asset amount
	LiquidityAssetAmount AssetAmount

	// Slippage is a slippage
	Slippage float64
}

// AmountsOutWithSlippage calculates asset amount out after applying the slippage
func (b *BurnQuote) AmountsOutWithSlippage() (map[uint64]AssetAmount, error) {
	out := make(map[uint64]AssetAmount)
	for k := range b.AmountsOut {
		amountOut := b.AmountsOut[k]
		amountWithSlippage, err := amountOut.Mul(nil, &b.Slippage)
		if err != nil {
			return nil, err
		}

		res, err := amountOut.Sub(amountWithSlippage, nil)
		if err != nil {
			return nil, err
		}

		out[k] = *res
	}

	return out, nil
}
