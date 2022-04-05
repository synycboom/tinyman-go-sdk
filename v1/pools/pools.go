package pools

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/contracts"
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

// LogicSig returns a logic signature account
func (p *Pool) LogicSig() (*crypto.LogicSigAccount, error) {
	return contracts.PoolLogicSigAccount(p.ValidatorAppID, p.Asset1.ID, p.Asset2.ID)
}

// Address returns a logic signature address
func (p *Pool) Address() (string, error) {
	poolAccount, err := p.LogicSig()
	if err != nil {
		return "", err
	}

	addr, err := poolAccount.Address()
	if err != nil {
		return "", err
	}

	return addr.String(), nil
}

// Asset1Price returns asset1 price
func (p *Pool) Asset1Price() float64 {
	return float64(p.Asset2Reserves) / float64(p.Asset1Reserves)
}

// Asset2Price returns asset2 price
func (p *Pool) Asset2Price() float64 {
	return float64(p.Asset1Reserves) / float64(p.Asset2Reserves)
}

// Info returns pool information
func (p *Pool) Info() (*types.PoolInfo, error) {
	address, err := p.Address()
	if err != nil {
		return nil, err
	}

	return &types.PoolInfo{
		Address:                         address,
		Asset1ID:                        p.Asset1.ID,
		Asset2ID:                        p.Asset2.ID,
		Asset1UnitName:                  p.Asset1.UnitName,
		Asset2UnitName:                  p.Asset2.UnitName,
		LiquidityAssetID:                p.LiquidityAsset.ID,
		LiquidityAssetName:              p.LiquidityAsset.Name,
		Asset1Reserves:                  p.Asset1Reserves,
		Asset2Reserves:                  p.Asset2Reserves,
		IssuedLiquidity:                 p.IssuedLiquidity,
		UnclaimedProtocolFee:            p.UnclaimedProtocolFee,
		OutstandingAsset1Amount:         p.OutstandingAsset1Amount,
		OutstandingAsset2Amount:         p.OutstandingAsset2Amount,
		OutstandingLiquidityAssetAmount: p.OutstandingLiquidityAssetAmount,
		ValidatorAppID:                  p.ValidatorAppID,
		AlgoBalance:                     p.AlgoBalance,
		Round:                           p.LastRefreshedRound,
	}, nil
}

// Convert converts one asset to another
func (p *Pool) Convert(amount *types.AssetAmount) (*types.AssetAmount, error) {
	if amount == nil {
		return nil, fmt.Errorf("amount is required")
	}

	if amount.Asset.Equal(p.Asset1) {
		price := p.Asset1Price()
		newAmount, err := amount.Mul(nil, &price)
		if err != nil {
			return nil, err
		}

		return &types.AssetAmount{
			Asset:  p.Asset2,
			Amount: newAmount.Amount,
		}, nil
	} else if amount.Asset.Equal(p.Asset2) {
		price := p.Asset2Price()
		newAmount, err := amount.Mul(nil, &price)
		if err != nil {
			return nil, err
		}

		return &types.AssetAmount{
			Asset:  p.Asset1,
			Amount: newAmount.Amount,
		}, nil
	}

	return nil, fmt.Errorf("mismatch asset")
}

// FetchStateInt returns an application state int value of the pool by a given key
func (p *Pool) FetchStateInt(ctx context.Context, key string) (uint64, error) {
	poolAddress, err := p.Address()
	if err != nil {
		return 0, err
	}

	accountInfo, err := p.ac.AccountInformation(poolAddress).Do(ctx)
	if err != nil {
		return 0, err
	}

	if len(accountInfo.AppsLocalState) == 0 {
		return 0, fmt.Errorf("application has no local state")
	}

	validatorAppState := make(map[string]models.TealValue)
	for _, kv := range accountInfo.AppsLocalState[0].KeyValue {
		validatorAppState[kv.Key] = kv.Value
	}

	return utils.StateInt(validatorAppState, key), nil
}

// FetchStateBytes returns an application state bytes value of the pool by a given key
func (p *Pool) FetchStateBytes(ctx context.Context, key string) ([]byte, error) {
	poolAddress, err := p.Address()
	if err != nil {
		return nil, err
	}

	accountInfo, err := p.ac.AccountInformation(poolAddress).Do(ctx)
	if err != nil {
		return nil, err
	}

	if len(accountInfo.AppsLocalState) == 0 {
		return nil, fmt.Errorf("application has no local state")
	}

	validatorAppState := make(map[string]models.TealValue)
	for _, kv := range accountInfo.AppsLocalState[0].KeyValue {
		validatorAppState[kv.Key] = kv.Value
	}

	return utils.StateBytes(validatorAppState, key), nil
}
