package types

import (
	"fmt"
	"math"
)

// AssetAmount represents an asset amount
type AssetAmount struct {
	// Asset is the underlying asset
	Asset *Asset

	// Amount is the amount of the asset
	Amount uint64
}

// NewAssetAmount creates a new asset amount
func NewAssetAmount(asset *Asset, amount uint64) (*AssetAmount, error) {
	if asset == nil {
		return nil, fmt.Errorf("an asset cannot be nil")
	}

	return &AssetAmount{Asset: asset, Amount: amount}, nil
}

// Mul multiplies the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Mul(other *AssetAmount, numOther *uint64) (*AssetAmount, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return nil, err
	}
	if other == nil && numOther != nil {
		return &AssetAmount{
			Asset:  a.Asset,
			Amount: a.Amount * *numOther,
		}, nil
	}

	if !a.Asset.Equal(other.Asset) {
		return nil, fmt.Errorf("the other asset is mismatched")
	}

	return &AssetAmount{
		Asset:  a.Asset,
		Amount: a.Amount * other.Amount,
	}, nil
}

// Add adds the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Add(other *AssetAmount, numOther *uint64) (*AssetAmount, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return nil, err
	}
	if other == nil && numOther != nil {
		return &AssetAmount{
			Asset:  a.Asset,
			Amount: a.Amount + *numOther,
		}, nil
	}

	if !a.Asset.Equal(other.Asset) {
		return nil, fmt.Errorf("the other asset is mismatched")
	}

	return &AssetAmount{
		Asset:  a.Asset,
		Amount: a.Amount + other.Amount,
	}, nil
}

// Sub subtracts the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Sub(other *AssetAmount, numOther *uint64) (*AssetAmount, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return nil, err
	}
	if other == nil && numOther != nil {
		return &AssetAmount{
			Asset:  a.Asset,
			Amount: a.Amount - *numOther,
		}, nil
	}

	if !a.Asset.Equal(other.Asset) {
		return nil, fmt.Errorf("the other asset is mismatched")
	}

	return &AssetAmount{
		Asset:  a.Asset,
		Amount: a.Amount - other.Amount,
	}, nil
}

// Div divides the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Div(other *AssetAmount, numOther *uint64) (*AssetAmount, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return nil, err
	}
	if other == nil && numOther != nil {
		if *numOther == 0 {
			return nil, fmt.Errorf("the other value is zero")
		}

		return &AssetAmount{
			Asset:  a.Asset,
			Amount: a.Amount / *numOther,
		}, nil
	}

	if !a.Asset.Equal(other.Asset) {
		return nil, fmt.Errorf("the other asset is mismatched")
	}

	if other.Amount == 0 {
		return nil, fmt.Errorf("the other value is zero")
	}

	return &AssetAmount{
		Asset:  a.Asset,
		Amount: a.Amount / other.Amount,
	}, nil
}

// Gt compares (>) the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Gt(other *AssetAmount, numOther *uint64) (bool, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return false, err
	}
	if other == nil && numOther != nil {
		return a.Amount > *numOther, nil
	}
	if !a.Asset.Equal(other.Asset) {
		return false, fmt.Errorf("the other asset is mismatched")
	}

	return a.Amount > other.Amount, nil
}

// Lt compares (<) the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Lt(other *AssetAmount, numOther *uint64) (bool, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return false, err
	}
	if other == nil && numOther != nil {
		return a.Amount < *numOther, nil
	}
	if !a.Asset.Equal(other.Asset) {
		return false, fmt.Errorf("the other asset is mismatched")
	}

	return a.Amount < other.Amount, nil
}

// Eq compares (=) the asset with the other (if it is not nil) or otherwise numOther, and return a new one
func (a *AssetAmount) Eq(other *AssetAmount, numOther *uint64) (bool, error) {
	if err := validateAssetAmountParams(other, numOther); err != nil {
		return false, err
	}
	if other == nil && numOther != nil {
		return a.Amount == *numOther, nil
	}
	if !a.Asset.Equal(other.Asset) {
		return false, fmt.Errorf("the other asset is mismatched")
	}

	return a.Amount == other.Amount, nil
}

// String returns a string representing an asest amount
func (a *AssetAmount) String() string {
	amount := float64(a.Amount) / (math.Pow(10, float64(a.Asset.Decimals)))

	return fmt.Sprintf("%s('%f')", a.Asset.UnitName, amount)
}

func validateAssetAmountParams(other *AssetAmount, numOther *uint64) error {
	if other == nil && numOther == nil {
		return fmt.Errorf("requires one parameter")
	}

	return nil
}
