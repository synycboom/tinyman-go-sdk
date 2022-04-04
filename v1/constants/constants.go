package constants

const (
	// TinyManURL is a tiny man web url
	TinyManURL = "https://tinyman.org"

	// AlgodTestnetHost is the algorand test net url
	AlgodTestnetHost = "https://testnet-api.algonode.cloud"

	// AlgodMainnetHost is the algorand main net url
	AlgodMainnetHost = "https://mainnet-api.algonode.cloud"

	// TestnetValidatorAppIdV1_1 is the Tinyman test net validator app id version 1.1
	TestnetValidatorAppIdV1_1 uint64 = 62368684

	// MainnetValidatorAppIdV1_1 is the Tinyman main net validator app id version 1.1
	MainnetValidatorAppIdV1_1 uint64 = 552635992

	// TestnetValidatorAppId is an alias for the current Tinyman test net validator app id
	TestnetValidatorAppId = TestnetValidatorAppIdV1_1

	// MainnetValidatorAppId is an alias for the current Tinyman main net validator app id
	MainnetValidatorAppId = MainnetValidatorAppIdV1_1

	// BurnFee is a burn transaction fee
	BurnFee uint64 = 4000

	// MintFee is a mint transaction fee
	MintFee uint64 = 2000

	// RedeemFee is a redeem transaction fee
	RedeemFee uint64 = 2000

	// SwapFee is a swap transaction fee
	SwapFee uint64 = 2000
)
