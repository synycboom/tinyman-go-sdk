package prepare

import (
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
)

// RedeemTransactions prepares a transaction group to redeem a specified excess asset amount from a pool.
func RedeemTransactions(
	validatorAppID,
	asset1ID,
	asset2ID,
	liquidityAssetID,
	assetID,
	assetAmount uint64,
	senderAddress string,
	sp types.SuggestedParams,
) (*utils.TransactionGroup, error) {
	var err error
	var tx1 types.Transaction
	var tx2 types.Transaction
	var tx3 types.Transaction

	poolAccount, err := contracts.PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID)
	if err != nil {
		return nil, err
	}

	poolAddress, err := poolAccount.Address()
	if err != nil {
		return nil, err
	}

	tx1, err = future.MakePaymentTxn(senderAddress, poolAddress.String(), constants.RedeemFee, []byte("fee"), "", sp)
	if err != nil {
		return nil, err
	}

	foreignAssets := []uint64{asset1ID, asset2ID, liquidityAssetID}
	if asset2ID == 0 {
		foreignAssets = []uint64{asset1ID, liquidityAssetID}
	}

	tx2, err = future.MakeApplicationNoOpTx(
		validatorAppID,
		[][]byte{[]byte("redeem")},
		[]string{senderAddress},
		nil,
		foreignAssets,
		sp,
		poolAddress,
		nil,
		types.Digest{},
		[32]byte{},
		types.Address{},
	)
	if err != nil {
		return nil, err
	}

	if assetID != 0 {
		tx3, err = future.MakeAssetTransferTxn(poolAddress.String(), senderAddress, assetAmount, nil, sp, "", assetID)
		if err != nil {
			return nil, err
		}
	} else {
		tx3, err = future.MakePaymentTxn(poolAddress.String(), senderAddress, assetAmount, nil, "", sp)
		if err != nil {
			return nil, err
		}
	}

	txGroup, err := utils.NewTransactionGroup([]types.Transaction{tx1, tx2, tx3})
	if err != nil {
		return nil, err
	}

	if err := txGroup.SignWithLogicSig(poolAccount); err != nil {
		return nil, err
	}

	return txGroup, nil
}
