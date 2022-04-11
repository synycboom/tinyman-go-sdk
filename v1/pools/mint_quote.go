package pools

import (
	"context"
	"fmt"
	"math"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
)

// FetchMintQuote returns a mint quote
func (p *Pool) FetchMintQuote(
	ctx context.Context,
	amountA *types.AssetAmount,
	amountB *types.AssetAmount,
	slippage float64,
) (*types.MintQuote, error) {
	if amountA == nil {
		return nil, fmt.Errorf("amountA is required")
	}
	if slippage == 0 {
		slippage = 0.05
	}

	var liquidityAssetAmount uint64
	var amount1 *types.AssetAmount
	var amount2 *types.AssetAmount
	if amountA.Asset.Equal(p.Asset1) {
		amount1 = amountA
		amount2 = amountB
	} else {
		amount1 = amountB
		amount2 = amountA
	}

	if err := p.Refresh(ctx, nil); err != nil {
		return nil, err
	}

	if !p.exists {
		return nil, fmt.Errorf("pool has not been bootstrapped yet")
	}

	if p.IssuedLiquidity > 0 {
		if amount1 == nil {
			amount, err := p.Convert(amount2)
			if err != nil {
				return nil, err
			}

			amount1 = amount
		}
		if amount2 == nil {
			amount, err := p.Convert(amount1)
			if err != nil {
				return nil, err
			}

			amount2 = amount
		}

		liquidityAssetAmount = uint64(
			math.Min(
				float64(utils.BigIntDiv(
					utils.BigIntMul(
						utils.ToBigUint(amount1.Amount),
						utils.ToBigUint(p.IssuedLiquidity),
					),
					utils.ToBigUint(p.Asset1Reserves),
				).Uint64()),
				float64(utils.BigIntDiv(
					utils.BigIntMul(
						utils.ToBigUint(amount2.Amount),
						utils.ToBigUint(p.IssuedLiquidity),
					),
					utils.ToBigUint(p.Asset2Reserves),
				).Uint64()),
			),
		)
	} else {
		if amount1 == nil || amount2 == nil {
			return nil, fmt.Errorf("amounts required for both assets for first mint")
		}

		liquidityAssetAmount = utils.BigIntSub(
			utils.BigIntSqrt(
				utils.BigIntMul(
					utils.ToBigUint(amount1.Amount),
					utils.ToBigUint(amount2.Amount),
				),
			),
			utils.ToBigUint(1000),
		).Uint64()

		slippage = 0
	}

	return &types.MintQuote{
		AmountsIn: map[uint64]types.AssetAmount{
			p.Asset1.ID: *amount1,
			p.Asset2.ID: *amount2,
		},
		LiquidityAssetAmount: types.AssetAmount{
			Asset:  p.LiquidityAsset,
			Amount: liquidityAssetAmount,
		},
		Slippage: slippage,
	}, nil
}
