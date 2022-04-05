package prepare

import (
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
)

// MintTransactions prepares a transaction group to mint the liquidity pool asset amount in exchange for pool assets.
func MintTransactions(
	validatorAppID,
	asset1ID,
	asset2ID,
	liquidityAssetID,
	asset1Amount,
	asset2Amount,
	liquidityAssetAmount uint64,
	senderAddress string,
	sp types.SuggestedParams,
) (*utils.TransactionGroup, error) {
	var err error
	var tx1 types.Transaction
	var tx2 types.Transaction
	var tx3 types.Transaction
	var tx4 types.Transaction
	var tx5 types.Transaction

	poolAccount, err := contracts.PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID)
	if err != nil {
		return nil, err
	}

	poolAddress, err := poolAccount.Address()
	if err != nil {
		return nil, err
	}

	tx1, err = future.MakePaymentTxn(senderAddress, poolAddress.String(), constants.MintFee, []byte("fee"), "", sp)
	if err != nil {
		return nil, err
	}

	foreignAssets := []uint64{asset1ID, asset2ID, liquidityAssetID}
	if asset2ID == 0 {
		foreignAssets = []uint64{asset1ID, liquidityAssetID}
	}
	tx2, err = future.MakeApplicationNoOpTx(
		validatorAppID,
		[][]byte{[]byte("mint")},
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

	tx3, err = future.MakeAssetTransferTxn(senderAddress, poolAddress.String(), asset1Amount, nil, sp, "", asset1ID)
	if err != nil {
		return nil, err
	}

	if asset2ID > 0 {
		tx4, err = future.MakeAssetTransferTxn(senderAddress, poolAddress.String(), asset2Amount, nil, sp, "", asset2ID)
		if err != nil {
			return nil, err
		}
	} else {
		tx4, err = future.MakePaymentTxn(senderAddress, poolAddress.String(), asset2Amount, nil, "", sp)
		if err != nil {
			return nil, err
		}
	}

	tx5, err = future.MakeAssetTransferTxn(poolAddress.String(), senderAddress, liquidityAssetAmount, nil, sp, "", liquidityAssetID)
	if err != nil {
		return nil, err
	}

	txGroup, err := utils.NewTransactionGroup([]types.Transaction{tx1, tx2, tx3, tx4, tx5})
	if err != nil {
		return nil, err
	}

	if err := txGroup.SignWithLogicSig(poolAccount); err != nil {
		return nil, err
	}

	return txGroup, nil
}
