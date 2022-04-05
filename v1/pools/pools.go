package pools

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// Pool represents a liquidity pool
type Pool struct {
	ac     *algod.Client
	exists bool

	ValidatorAppID                  uint64
	LiquidityAsset                  *types.Asset
	Asset1                          *types.Asset
	Asset2                          *types.Asset
	AlgoBalance                     uint64
	Asset1Reserves                  uint64
	Asset2Reserves                  uint64
	IssuedLiquidity                 uint64
	MinBalance                      uint64
	UnclaimedProtocolFee            uint64
	OutstandingAsset1Amount         uint64
	OutstandingAsset2Amount         uint64
	OutstandingLiquidityAssetAmount uint64
	LastRefreshedRound              uint64
	UserAddress                     string
}

func NewPool(
	ac *algod.Client,
	assetA,
	assetB *types.Asset,
	info *types.PoolInfo,
	validatorAppID uint64,
	userAddress string,
	fetch bool,
) (*Pool, error) {
	ctx := context.Background()
	p := Pool{}
	p.ac = ac
	p.ValidatorAppID = validatorAppID
	p.UserAddress = userAddress

	if ac == nil {
		return nil, fmt.Errorf("algod client is required")
	}
	if assetA == nil || assetB == nil {
		return nil, fmt.Errorf("both assetA and assetB are required")
	}

	if assetA.IsFetchingRequired() {
		if err := assetA.Fetch(ctx, p.ac); err != nil {
			return nil, err
		}
	}

	if assetB.IsFetchingRequired() {
		if err := assetB.Fetch(ctx, p.ac); err != nil {
			return nil, err
		}
	}

	if assetA.ID > assetB.ID {
		p.Asset1 = assetA
		p.Asset2 = assetB
	} else {
		p.Asset1 = assetB
		p.Asset2 = assetA
	}

	if fetch {
		if err := p.Refresh(ctx, nil); err != nil {
			return nil, err
		}
	} else if info != nil {
		if err := p.UpdateFromInfo(ctx, info); err != nil {
			return nil, err
		}
	}

	return &p, nil
}

// FromAccountInfo create a pool from an account
func FromAccountInfo(account models.Account, ac *algod.Client, userAddress string) (*Pool, error) {
	info, err := poolInfoFromAccountInfo(account)
	if err != nil {
		return nil, err
	}

	assetA := &types.Asset{ID: info.Asset1ID}
	assetB := &types.Asset{ID: info.Asset2ID}

	return NewPool(ac, assetA, assetB, info, info.ValidatorAppID, userAddress, false)
}

// Refresh refreshes pool information
func (p *Pool) Refresh(ctx context.Context, info *types.PoolInfo) error {
	if info == nil {
		i, err := PoolInfo(ctx, p.ac, p.ValidatorAppID, p.Asset1.ID, p.Asset2.ID)
		if err != nil {
			return err
		}

		info = i
	}

	if info == nil {
		return nil
	}

	if err := p.UpdateFromInfo(ctx, info); err != nil {
		return err
	}

	return nil
}

// UpdateFromInfo updates pool information from a given pool info
func (p *Pool) UpdateFromInfo(ctx context.Context, info *types.PoolInfo) error {
	if info.LiquidityAssetID > 0 {
		p.exists = true
	}

	p.LiquidityAsset = &types.Asset{
		ID:       info.LiquidityAssetID,
		Decimals: 6,
		Name:     info.LiquidityAssetName,
		UnitName: constants.LiquidityAssetUnitName,
	}
	p.Asset1Reserves = info.Asset1Reserves
	p.Asset2Reserves = info.Asset2Reserves
	p.IssuedLiquidity = info.IssuedLiquidity
	p.UnclaimedProtocolFee = info.UnclaimedProtocolFee
	p.OutstandingAsset1Amount = info.OutstandingAsset1Amount
	p.OutstandingAsset2Amount = info.OutstandingAsset2Amount
	p.OutstandingLiquidityAssetAmount = info.OutstandingLiquidityAssetAmount
	p.LastRefreshedRound = info.Round
	p.AlgoBalance = info.AlgoBalance
	p.MinBalance = p.MinimumBalance()
	if p.Asset2.ID == 0 {
		p.Asset2Reserves = (p.AlgoBalance - p.MinBalance) - p.OutstandingAsset2Amount
	}
	if p.IssuedLiquidity > 0 {
		asset, err := p.ac.GetAssetByID(p.LiquidityAsset.ID).Do(ctx)
		if err != nil {
			return err
		}

		p.LiquidityAsset = &types.Asset{
			ID:       asset.Index,
			Decimals: asset.Params.Decimals,
			Name:     asset.Params.Name,
			UnitName: asset.Params.UnitName,
		}
	}

	return nil
}

// MinimumBalance calculates minimum balance
func (p *Pool) MinimumBalance() uint64 {
	numAssets := 3
	if p.Asset2.ID == 0 {
		numAssets = 2
	}

	numCreatedApps := 0
	numLocalApps := 1
	totalUints := 16
	totalByteSlices := 0

	return uint64(constants.MinBalancePerAccount +
		(constants.MinBalancePerAsset * numAssets) +
		(constants.MinBalancePerApp * (numCreatedApps + numLocalApps)) +
		(constants.MinBalancePerAppUint * totalUints) +
		(constants.MinBalancePerAppByteSlice * totalByteSlices))
}
