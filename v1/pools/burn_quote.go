package pools

import (
	"context"
	"fmt"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
)

// FetchBurnQuote returns a burn quote
func (p *Pool) FetchBurnQuote(ctx context.Context, liquidityAsset types.AssetAmount, slippage float64) (*types.BurnQuote, error) {
	if slippage == 0 {
		slippage = 0.05
	}

	if !liquidityAsset.Asset.Equal(p.LiquidityAsset) {
		return nil, fmt.Errorf("the liquidity asset is not the same as one in a pool")
	}

	if err := p.Refresh(ctx, nil); err != nil {
		return nil, err
	}

	asset1Amount := utils.BigIntDiv(
		utils.BigIntMul(
			utils.ToBigUint(liquidityAsset.Amount),
			utils.ToBigUint(p.Asset1Reserves),
		),
		utils.ToBigUint(p.IssuedLiquidity),
	)
	asset2Amount := utils.BigIntDiv(
		utils.BigIntMul(
			utils.ToBigUint(liquidityAsset.Amount),
			utils.ToBigUint(p.Asset2Reserves),
		),
		utils.ToBigUint(p.IssuedLiquidity),
	)

	return &types.BurnQuote{
		AmountsOut: map[uint64]types.AssetAmount{
			p.Asset1.ID: {
				Asset:  p.Asset1,
				Amount: asset1Amount.Uint64(),
			},
			p.Asset2.ID: {
				Asset:  p.Asset2,
				Amount: asset2Amount.Uint64(),
			},
		},
		LiquidityAssetAmount: liquidityAsset,
		Slippage:             slippage,
	}, nil
}
