package types

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// Asset is an Algorand token
type Asset struct {
	isFetched bool

	// ID is an asset id
	ID uint64

	// Decimals is an asset decimals
	Decimals uint64

	// Name is an asset name
	Name string

	// UnitName is an asset unit name
	UnitName string
}

// NewAsset creates an asset
func NewAsset(id, decimals uint64, name, unitName string) *Asset {
	return &Asset{
		ID:       id,
		Decimals: decimals,
		Name:     name,
		UnitName: unitName,
	}
}

// String returns a string representing the asset
func (a *Asset) String() string {
	return fmt.Sprintf("Asset(%s - %d)", a.UnitName, a.ID)
}

// Fetch fetches and updates the asset information
func (a *Asset) Fetch(ctx context.Context, ac *algod.Client) error {
	if a.ID > 0 {
		asset, err := ac.GetAssetByID(a.ID).Do(ctx)
		if err != nil {
			return err
		}

		a.Name = asset.Params.Name
		a.UnitName = asset.Params.UnitName
		a.Decimals = asset.Params.Decimals
	} else {
		a.Name = constants.AlgoTokenName
		a.UnitName = constants.AlgoTokenUnitName
		a.Decimals = constants.AlgoTokenDecimals
	}

	a.isFetched = true

	return nil
}

// Equal checks the equality of the asset with an other
func (a *Asset) Equal(other *Asset) bool {
	if other == nil {
		return false
	}

	return a.ID == other.ID
}

// IsFetchingRequired checks whether the asset needs information fetching
func (a *Asset) IsFetchingRequired() bool {
	if a.isFetched {
		return false
	}

	if len(a.Name) == 0 || len(a.UnitName) == 0 || a.Decimals == 0 {
		return true
	}

	return false
}
