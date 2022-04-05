package tinyman

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	algoTypes "github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/pools"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// Client represents the Tinyman client
type Client struct {
	assetCache map[uint64]types.Asset
	ac         *algod.Client

	UserAddress    string
	ValidatorAppID uint64
}

// NewClient create a Tinyman client
func NewClient(ac *algod.Client, validatorAppID uint64, userAddress string) *Client {
	return &Client{
		ac:             ac,
		ValidatorAppID: validatorAppID,
		UserAddress:    userAddress,
		assetCache:     make(map[uint64]types.Asset),
	}
}

// NewTestNetClient create a test net Tinyman client
func NewTestNetClient(ac *algod.Client, userAddress string) (*Client, error) {
	if ac == nil {
		a, err := algod.MakeClient(constants.AlgodTestnetHost, "")
		if err != nil {
			return nil, err
		}

		ac = a
	}

	return NewClient(ac, constants.TestnetValidatorAppId, userAddress), nil
}

// NewMainNetClient create a main net Tinyman client
func NewMainNetClient(ac *algod.Client, userAddress string) (*Client, error) {
	if ac == nil {
		a, err := algod.MakeClient(constants.AlgodMainnetHost, "")
		if err != nil {
			return nil, err
		}

		ac = a
	}

	return NewClient(ac, constants.MainnetValidatorAppId, userAddress), nil
}

// FetchPool fetches a pool for given asset1 and asset2
func (c *Client) FetchPool(asset1, asset2 *types.Asset, fetch bool) (*pools.Pool, error) {
	if asset1 == nil || asset2 == nil {
		return nil, fmt.Errorf("asset1 and asset2 are required")
	}

	return pools.NewPool(c.ac, asset1, asset2, nil, c.ValidatorAppID, c.UserAddress, fetch)
}

// FetchAsset fetches an asset for a given asset id
func (c *Client) FetchAsset(ctx context.Context, assetID uint64) (*types.Asset, error) {
	if _, ok := c.assetCache[assetID]; !ok {
		asset := types.Asset{
			ID: assetID,
		}
		if err := asset.Fetch(ctx, c.ac); err != nil {
			return nil, err
		}

		c.assetCache[assetID] = asset
	}

	asset := c.assetCache[assetID]

	return &asset, nil
}

// Submit submits a transaction group to the blockchain
func (c *Client) Submit(ctx context.Context, txGroup *utils.TransactionGroup, wait bool) (string, error) {
	txID, err := txGroup.Submit(ctx, c.ac, wait)
	if err != nil {
		return "", err
	}

	return txID, nil
}

// PrepareAppOptInTransaction prepares an app opt-in transaction and returns a transaction group
func (c *Client) PrepareAppOptInTransaction(ctx context.Context, userAddress string) (*utils.TransactionGroup, error) {
	if len(userAddress) == 0 {
		userAddress = c.UserAddress
	}

	sp, err := c.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	return prepare.AppOptInTransactions(c.ValidatorAppID, userAddress, sp)
}

// PrepareAssetOptInTransactions prepares asset opt-in transaction and returns a transaction group
func (c *Client) PrepareAssetOptInTransactions(ctx context.Context, assetID uint64, userAddress string) (*utils.TransactionGroup, error) {
	if len(userAddress) == 0 {
		userAddress = c.UserAddress
	}

	sp, err := c.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.AssetOptInTransactions(assetID, userAddress, sp)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// FetchExcessAmount fetches user's excess amounts and returns redeem quotes
func (c *Client) FetchExcessAmount(ctx context.Context, userAddr string) ([]types.RedeemQuote, error) {
	var pools []types.RedeemQuote
	if len(userAddr) == 0 {
		userAddr = c.UserAddress
	}

	accountInfo := c.ac.AccountInformation(userAddr)
	account, err := accountInfo.Do(ctx)
	if err != nil {
		return pools, err
	}

	var validatorApp *models.ApplicationLocalState
	for _, as := range account.AppsLocalState {
		if as.Id == c.ValidatorAppID {
			validatorApp = &as
			break
		}
	}

	if validatorApp == nil {
		return pools, nil
	}

	validatorAppState := make(map[string]models.TealValue)
	for _, kv := range validatorApp.KeyValue {
		validatorAppState[kv.Key] = kv.Value
	}

	for k := range validatorAppState {
		key, err := base64.StdEncoding.DecodeString(k)
		if err != nil {
			return nil, err
		}

		if key[len(key)-9] == 'e' {
			value := validatorAppState[k].Uint
			assetID := binary.BigEndian.Uint64(key[len(key)-8:])
			poolAddress, err := algoTypes.EncodeAddress(key[:len(key)-9])
			if err != nil {
				return nil, err
			}

			asset, err := c.FetchAsset(ctx, assetID)
			if err != nil {
				return nil, err
			}

			pools = append(pools, types.RedeemQuote{
				Amount: types.AssetAmount{
					Asset:  asset,
					Amount: value,
				},
				PoolAddress: poolAddress,
			})
		}
	}

	return pools, nil
}

// IsOptedIn checkes whether a user opted in for the application
func (c *Client) IsOptedIn(ctx context.Context, userAddr string) (bool, error) {
	if len(userAddr) == 0 {
		userAddr = c.UserAddress
	}

	accountInfo := c.ac.AccountInformation(userAddr)
	account, err := accountInfo.Do(ctx)
	if err != nil {
		return false, err
	}

	for _, as := range account.AppsLocalState {
		if as.Id == c.ValidatorAppID {
			return true, nil
		}
	}

	return false, nil
}

// IsOptedIn checkes whether a user opted in for asset
func (c *Client) IsAssetOptedIn(ctx context.Context, assetID uint64, userAddr string) (bool, error) {
	if len(userAddr) == 0 {
		userAddr = c.UserAddress
	}

	accountInfo := c.ac.AccountInformation(userAddr)
	account, err := accountInfo.Do(ctx)
	if err != nil {
		return false, err
	}

	for _, asset := range account.Assets {
		if asset.AssetId == assetID {
			return true, nil
		}
	}

	return false, nil
}

// Balance returns an asset balance of a user
func (c *Client) Balance(ctx context.Context, asset types.Asset, userAddress string) (*types.AssetAmount, error) {
	if len(userAddress) == 0 {
		userAddress = c.UserAddress
	}

	accountInfo := c.ac.AccountInformation(userAddress)
	account, err := accountInfo.Do(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range account.Assets {
		if a.AssetId == asset.ID {
			return &types.AssetAmount{
				Asset:  &asset,
				Amount: a.Amount,
			}, nil
		}
	}

	return nil, nil
}
