package types_test

import (
	"testing"

	"github.com/synycboom/tinyman-go-sdk/types"
)

func TestNewAsetAmount(t *testing.T) {
	_, err := types.NewAssetAmount(nil, 0)
	if err == nil {
		t.Error("It should not return an error if an asset is nil")

		return
	}

	asset := types.NewAsset(0, 0, "", "")
	assetAmount, err := types.NewAssetAmount(asset, 10)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount.Amount != 10 {
		t.Error("Amount is not correct")

		return
	}
	if !assetAmount.Asset.Equal(asset) {
		t.Error("Asset is not equal")

		return
	}
}

func prepare(t *testing.T) (*types.AssetAmount, *types.AssetAmount, *types.AssetAmount) {
	asset1 := types.NewAsset(1, 0, "1", "1")
	asset2 := types.NewAsset(2, 0, "2", "2")
	assetAmount1_1, err := types.NewAssetAmount(asset1, 10)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return nil, nil, nil
	}
	assetAmount1_2, err := types.NewAssetAmount(asset1, 20)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return nil, nil, nil
	}
	assetAmount2_1, err := types.NewAssetAmount(asset2, 10)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return nil, nil, nil
	}

	return assetAmount1_1, assetAmount1_2, assetAmount2_1
}

func TestAsetAmountMul(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Mul(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Mul(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}
	out, err := assetAmount1_1.Mul(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount*assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d * %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
	out, err = assetAmount1_1.Mul(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount*assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d * %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
}

func TestAsetAmountAdd(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Add(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Add(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}
	out, err := assetAmount1_1.Add(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount+assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d + %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
	out, err = assetAmount1_1.Add(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount+assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d + %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
}

func TestAsetAmountSub(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Sub(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Sub(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}
	out, err := assetAmount1_1.Sub(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount-assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d - %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
	out, err = assetAmount1_1.Sub(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount-assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d - %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
}

func TestAsetAmountDiv(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Div(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Div(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}
	out, err := assetAmount1_1.Div(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount/assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d / %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}
	out, err = assetAmount1_1.Div(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if assetAmount1_1.Amount/assetAmount1_2.Amount != out.Amount {
		t.Errorf("Wrong calculation for %d / %d == %d", assetAmount1_1.Amount, assetAmount1_2.Amount, out.Amount)

		return
	}

	assetAmount1_2.Amount = 0
	if _, err := assetAmount1_1.Div(assetAmount1_2, nil); err == nil {
		t.Error("It should return a zero division error")

		return
	}
	if _, err := assetAmount1_1.Div(nil, &assetAmount1_2.Amount); err == nil {
		t.Error("It should return a zero division error")

		return
	}
}

func TestAsetAmountGt(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Gt(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Gt(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}

	gt, err := assetAmount1_1.Gt(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if gt {
		t.Errorf("Wrong calculation for %d > %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, gt)

		return
	}

	gt, err = assetAmount1_1.Gt(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if gt {
		t.Errorf("Wrong calculation for %d > %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, gt)

		return
	}

	assetAmount1_2.Amount = assetAmount1_1.Amount - 1
	gt, err = assetAmount1_1.Gt(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !gt {
		t.Errorf("Wrong calculation for %d > %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, gt)

		return
	}

	gt, err = assetAmount1_1.Gt(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !gt {
		t.Errorf("Wrong calculation for %d > %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, gt)

		return
	}
}

func TestAsetAmountLt(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Lt(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Lt(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}

	lt, err := assetAmount1_1.Lt(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !lt {
		t.Errorf("Wrong calculation for %d < %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, lt)

		return
	}

	lt, err = assetAmount1_1.Lt(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !lt {
		t.Errorf("Wrong calculation for %d < %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, lt)

		return
	}

	assetAmount1_2.Amount = assetAmount1_1.Amount - 1
	lt, err = assetAmount1_1.Lt(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if lt {
		t.Errorf("Wrong calculation for %d < %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, lt)

		return
	}

	lt, err = assetAmount1_1.Lt(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if lt {
		t.Errorf("Wrong calculation for %d < %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, lt)

		return
	}
}

func TestAsetAmountEq(t *testing.T) {
	assetAmount1_1, assetAmount1_2, assetAmount2_1 := prepare(t)
	if assetAmount1_1 == nil {
		return
	}

	if _, err := assetAmount1_1.Eq(nil, nil); err == nil {
		t.Error("It should return an error when two params are nil")

		return
	}
	if _, err := assetAmount1_1.Eq(assetAmount2_1, nil); err == nil {
		t.Error("It should return a mismatch error")

		return
	}

	eq, err := assetAmount1_1.Eq(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if eq {
		t.Errorf("Wrong calculation for %d == %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, eq)

		return
	}

	eq, err = assetAmount1_1.Eq(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if eq {
		t.Errorf("Wrong calculation for %d == %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, eq)

		return
	}

	assetAmount1_2.Amount = assetAmount1_1.Amount
	eq, err = assetAmount1_1.Eq(assetAmount1_2, nil)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !eq {
		t.Errorf("Wrong calculation for %d == %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, eq)

		return
	}
	eq, err = assetAmount1_1.Eq(nil, &assetAmount1_2.Amount)
	if err != nil {
		t.Errorf("It should not return an error: %s", err.Error())

		return
	}
	if !eq {
		t.Errorf("Wrong calculation for %d == %d == %v", assetAmount1_1.Amount, assetAmount1_2.Amount, eq)

		return
	}
}
