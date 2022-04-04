package types

// RedeemQuote represents a redeem quote
type RedeemQuote struct {
	// Amount is an asset amount that can be redeemed
	Amount AssetAmount

	// PoolAddress is a pool address
	PoolAddress string
}
