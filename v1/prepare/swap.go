package prepare

import (
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
)

// SwapTransactions Prepare a transaction group to swap assets.
func SwapTransactions(
	validatorAppID,
	asset1ID,
	asset2ID,
	liquidityAssetID,
	assetInID,
	assetInAmount,
	assetOutAmount uint64,
	swapType string,
	senderAddress string,
	sp types.SuggestedParams,
) (*utils.TransactionGroup, error) {
	var err error
	var tx1 types.Transaction
	var tx2 types.Transaction
	var tx3 types.Transaction
	var tx4 types.Transaction

	poolAccount, err := contracts.PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID)
	if err != nil {
		return nil, err
	}

	poolAddress, err := poolAccount.Address()
	if err != nil {
		return nil, err
	}

	assetOutID := asset1ID
	if assetInID == asset1ID {
		assetOutID = asset2ID
	}

	tx1, err = future.MakePaymentTxn(senderAddress, poolAddress.String(), constants.SwapFee, []byte("fee"), "", sp)
	if err != nil {
		return nil, err
	}

	appIdx := validatorAppID
	appArgs := [][]byte{[]byte("swap"), []byte(constants.SwapTypeMapping[swapType])}
	accounts := []string{senderAddress}
	foreignAssets := []uint64{asset1ID, asset2ID, liquidityAssetID}
	if asset2ID == 0 {
		foreignAssets = []uint64{asset1ID, liquidityAssetID}
	}

	tx2, err = future.MakeApplicationNoOpTx(
		appIdx,
		appArgs,
		accounts,
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

	if assetInID != 0 {
		tx3, err = future.MakeAssetTransferTxn(senderAddress, poolAddress.String(), assetInAmount, nil, sp, "", assetInID)
		if err != nil {
			return nil, err
		}
	} else {
		tx3, err = future.MakePaymentTxn(senderAddress, poolAddress.String(), assetInAmount, nil, "", sp)
		if err != nil {
			return nil, err
		}
	}

	if assetOutID != 0 {
		tx4, err = future.MakeAssetTransferTxn(poolAddress.String(), senderAddress, assetOutAmount, nil, sp, "", assetOutID)
		if err != nil {
			return nil, err
		}
	} else {
		tx4, err = future.MakePaymentTxn(poolAddress.String(), senderAddress, assetOutAmount, nil, "", sp)
		if err != nil {
			return nil, err
		}
	}

	txGroup, err := utils.NewTransactionGroup([]types.Transaction{tx1, tx2, tx3, tx4})
	if err != nil {
		return nil, err
	}

	if err := txGroup.SignWithLogicSig(poolAccount); err != nil {
		return nil, err
	}

	return txGroup, nil
}
