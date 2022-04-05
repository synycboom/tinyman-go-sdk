package prepare

import (
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/utils"
)

// AppOptInTransactions prepares a transaction group to opt-in of Tinyman
func AppOptInTransactions(validatorAppID uint64, senderAddress string, sp types.SuggestedParams) (*utils.TransactionGroup, error) {
	addr, err := types.DecodeAddress(senderAddress)
	if err != nil {
		return nil, err
	}

	tx, err := future.MakeApplicationOptInTx(
		validatorAppID,
		nil,
		nil,
		nil,
		nil,
		sp,
		addr,
		nil,
		types.Digest{},
		[32]byte{},
		types.Address{},
	)
	if err != nil {
		return nil, err
	}

	txGroup, err := utils.NewTransactionGroup([]types.Transaction{tx})
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// AssetOptInTransactions prepares a transaction group to opt-in an asset
func AssetOptInTransactions(assetID uint64, senderAddress string, sp types.SuggestedParams) (*utils.TransactionGroup, error) {
	tx, err := future.MakeAssetAcceptanceTxn(senderAddress, nil, sp, assetID)
	if err != nil {
		return nil, err
	}

	txGroup, err := utils.NewTransactionGroup([]types.Transaction{tx})
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}
