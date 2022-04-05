package pools

import (
	"context"
	"fmt"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// FetchFixedInputSwapQuote returns a fixed input swap quote
func (p *Pool) FetchFixedInputSwapQuote(ctx context.Context, amountIn types.AssetAmount, slippage float64) (*types.SwapQuote, error) {
	if slippage == 0 {
		slippage = 0.05
	}

	assetIn := amountIn.Asset
	assetInAmount := amountIn.Amount
	if err := p.Refresh(ctx, nil); err != nil {
		return nil, err
	}

	var assetOut *types.Asset
	var inputSupply uint64
	var outputSupply uint64
	if assetIn.Equal(p.Asset1) {
		assetOut = p.Asset2
		inputSupply = p.Asset1Reserves
		outputSupply = p.Asset2Reserves
	} else {
		assetOut = p.Asset1
		inputSupply = p.Asset2Reserves
		outputSupply = p.Asset1Reserves
	}

	if inputSupply == 0 || outputSupply == 0 {
		return nil, fmt.Errorf("pool has no liquidity")
	}

	bigInputSupply := utils.ToBigUint(inputSupply)
	bigOutputSupply := utils.ToBigUint(outputSupply)
	bigAssetInAmount := utils.ToBigUint(assetInAmount)
	k := utils.BigIntMul(bigInputSupply, bigOutputSupply)
	bigAssetInAmountMinusFee := utils.BigIntDiv(
		utils.BigIntMul(bigAssetInAmount, utils.ToBigUint(997)),
		utils.ToBigUint(1000),
	)
	swapFees := assetInAmount - bigAssetInAmountMinusFee.Uint64()
	bigAssetOutAmount := utils.BigIntSub(
		bigOutputSupply,
		utils.BigIntDiv(k, utils.BigIntAdd(bigInputSupply, bigAssetInAmountMinusFee)),
	)
	amountOut := types.AssetAmount{
		Asset:  assetOut,
		Amount: bigAssetOutAmount.Uint64(),
	}

	return &types.SwapQuote{
		SwapType:  constants.SwapFixedInput,
		AmountIn:  amountIn,
		AmountOut: amountOut,
		SwapFee: types.AssetAmount{
			Asset:  amountIn.Asset,
			Amount: swapFees,
		},
		Slippage: slippage,
	}, nil
}
