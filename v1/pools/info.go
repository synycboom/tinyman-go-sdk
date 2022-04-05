package pools

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
)

// PoolInfo returns pool information for the given asset1 and asset2
func PoolInfo(ctx context.Context, ac *algod.Client, validatorAppID, asset1ID, asset2ID uint64) (*types.PoolInfo, error) {
	poolAccount, err := contracts.PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID)
	if err != nil {
		return nil, err
	}

	poolAddress, err := poolAccount.Address()
	if err != nil {
		return nil, err
	}

	accountInfo, err := ac.AccountInformation(poolAddress.String()).Do(ctx)
	if err != nil {
		return nil, err
	}

	return poolInfoFromAccountInfo(accountInfo)
}

func poolInfoFromAccountInfo(accountInfo models.Account) (*types.PoolInfo, error) {
	if len(accountInfo.AppsLocalState) == 0 {
		return nil, nil
	}

	validatorAppID := accountInfo.AppsLocalState[0].Id
	validatorAppState := make(map[string]models.TealValue)
	for _, kv := range accountInfo.AppsLocalState[0].KeyValue {
		validatorAppState[kv.Key] = kv.Value
	}

	asset1ID := utils.StateInt(validatorAppState, "a1")
	asset2ID := utils.StateInt(validatorAppState, "a2")
	poolAccount, err := contracts.PoolLogicSigAccount(validatorAppID, asset1ID, asset2ID)
	if err != nil {
		return nil, err
	}

	poolAddress, err := poolAccount.Address()
	if err != nil {
		return nil, err
	}

	if accountInfo.Address != poolAddress.String() {
		return nil, fmt.Errorf(
			"pool address '%s' is not matched an account address '%s'",
			poolAddress.String(),
			accountInfo.Address,
		)
	}

	asset1Reserves := utils.StateInt(validatorAppState, "s1")
	asset2Reserves := utils.StateInt(validatorAppState, "s2")
	issuedLiquidity := utils.StateInt(validatorAppState, "ilt")
	unclaimedProtocolFees := utils.StateInt(validatorAppState, "p")
	if len(accountInfo.CreatedAssets) == 0 {
		return nil, fmt.Errorf("account does not have created assets")
	}

	liquidityAsset := accountInfo.CreatedAssets[0]
	liquidityAssetID := liquidityAsset.Index
	outstandingAsset1Key, err := utils.OutstandingAssetStateKey(asset1ID)
	if err != nil {
		return nil, err
	}
	outstandingAsset2Key, err := utils.OutstandingAssetStateKey(asset2ID)
	if err != nil {
		return nil, err
	}
	outstandingLiquidityAssetKey, err := utils.OutstandingAssetStateKey(liquidityAssetID)
	if err != nil {
		return nil, err
	}

	outstandingAsset1Amount := utils.StateInt(validatorAppState, string(outstandingAsset1Key))
	outstandingAsset2Amount := utils.StateInt(validatorAppState, string(outstandingAsset2Key))
	outstandingLiquidityAssetAmount := utils.StateInt(validatorAppState, string(outstandingLiquidityAssetKey))

	return &types.PoolInfo{
		Address:                         poolAddress.String(),
		Asset1ID:                        asset1ID,
		Asset2ID:                        asset2ID,
		LiquidityAssetID:                liquidityAssetID,
		LiquidityAssetName:              liquidityAsset.Params.Name,
		Asset1Reserves:                  asset1Reserves,
		Asset2Reserves:                  asset2Reserves,
		IssuedLiquidity:                 issuedLiquidity,
		UnclaimedProtocolFee:            unclaimedProtocolFees,
		OutstandingAsset1Amount:         outstandingAsset1Amount,
		OutstandingAsset2Amount:         outstandingAsset2Amount,
		OutstandingLiquidityAssetAmount: outstandingLiquidityAssetAmount,
		ValidatorAppID:                  validatorAppID,
		AlgoBalance:                     accountInfo.Amount,
		Round:                           accountInfo.Round,
	}, nil
}
