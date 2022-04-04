package types

// MintQuote represents a mint quote
type MintQuote struct {
	// AmountsIn is an asset mapping which maps between asset ids and asset amounts
	AmountsIn map[uint64]AssetAmount

	// LiquidityAssetAmount is a liquidity asset amount
	LiquidityAssetAmount AssetAmount

	// Slippage is a slippage
	Slippage float64
}

// LiquidityAssetAmountWithSlippage calculates liquidity asset after applying the slippage
func (m *MintQuote) LiquidityAssetAmountWithSlippage() (*AssetAmount, error) {
	amountWithSlippage, err := m.LiquidityAssetAmount.Mul(nil, &m.Slippage)
	if err != nil {
		return nil, err
	}

	return m.LiquidityAssetAmount.Sub(amountWithSlippage, nil)
}
