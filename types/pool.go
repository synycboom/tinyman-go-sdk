package types

// PoolInfo represents pool information
type PoolInfo struct {
	// Address is a pool address
	Address string

	// Asset1ID is an asset1 id
	Asset1ID uint64

	// Asset2ID is an asset2 id
	Asset2ID uint64

	// Asset1UnitName is an asset1 unit name
	Asset1UnitName string

	// Asset2UnitName is an asset2 unit name
	Asset2UnitName string

	// LiquidityAssetID is an asset id for the liquidity
	LiquidityAssetID uint64

	// LiquidityAssetID is an asset name for the liquidity
	LiquidityAssetName string

	// Asset1Reserves is an asset1's reserves value
	Asset1Reserves uint64

	// Asset2Reserves is an asset2's reserves value
	Asset2Reserves uint64

	// IssuedLiquidity is the total issued liquidity
	IssuedLiquidity uint64

	// UnclaimedProtocolFee is an unclaimed protocol fee
	UnclaimedProtocolFee uint64

	// OutstandingAsset1Amount is an outstanding asset1 amount
	OutstandingAsset1Amount uint64

	// OutstandingAsset2Amount is an outstanding asset2 amount
	OutstandingAsset2Amount uint64

	// OutstandingLiquidityAssetAmount is an outstanding liquidity asset amount
	OutstandingLiquidityAssetAmount uint64

	// ValidatorAppID is the validator app id
	ValidatorAppID uint64

	// AlgoBalance is a balance of the pool
	AlgoBalance uint64

	// Round is the latest fetch round
	Round uint64
}

// PoolPosition represents a user position in the pool
type PoolPosition struct {
	// Asset1 is an asset1
	Asset1 AssetAmount

	// Asset2 is an asset2
	Asset2 AssetAmount

	// LiquidityAsset is a liquidity asset
	LiquidityAsset AssetAmount

	// Share is a share of user which can be calculated as a percentage by (share * 100)
	Share float64
}
