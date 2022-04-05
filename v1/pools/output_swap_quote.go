package pools

import (
	"context"
	"fmt"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// FetchFixedOutputSwapQuote returns a fixed input swap quote
func (p *Pool) FetchFixedOutputSwapQuote(ctx context.Context, amountOut types.AssetAmount, slippage float64) (*types.SwapQuote, error) {
	if slippage == 0 {
		slippage = 0.05
	}

	assetOut := amountOut.Asset
	assetOutAmount := amountOut.Amount
	if err := p.Refresh(ctx, nil); err != nil {
		return nil, err
	}

	var assetIn *types.Asset
	var inputSupply uint64
	var outputSupply uint64
	if assetOut.Equal(p.Asset1) {
		assetIn = p.Asset2
		inputSupply = p.Asset2Reserves
		outputSupply = p.Asset1Reserves
	} else {
		assetIn = p.Asset1
		inputSupply = p.Asset1Reserves
		outputSupply = p.Asset2Reserves
	}

	if inputSupply == 0 || outputSupply == 0 {
		return nil, fmt.Errorf("pool has no liquidity")
	}

	bigInputSupply := utils.ToBigUint(inputSupply)
	bigOutputSupply := utils.ToBigUint(outputSupply)
	bigAssetOutAmount := utils.ToBigUint(assetOutAmount)
	k := utils.BigIntMul(bigInputSupply, bigOutputSupply)
	bigCalculatedAmountInWithoutFee := utils.BigIntSub(
		utils.BigIntDiv(k, utils.BigIntSub(bigOutputSupply, bigAssetOutAmount)),
		bigInputSupply,
	)
	bigAssetInAmount := utils.BigIntDiv(
		utils.BigIntMul(bigCalculatedAmountInWithoutFee, utils.ToBigUint(1000)),
		utils.ToBigUint(997),
	)
	bigSwapFees := utils.BigIntSub(bigAssetInAmount, bigCalculatedAmountInWithoutFee)
	amountIn := types.AssetAmount{
		Asset:  assetIn,
		Amount: bigAssetInAmount.Uint64(),
	}

	return &types.SwapQuote{
		SwapType:  constants.SwapFixedOutput,
		AmountIn:  amountIn,
		AmountOut: amountOut,
		SwapFee: types.AssetAmount{
			Asset:  amountIn.Asset,
			Amount: bigSwapFees.Uint64(),
		},
		Slippage: slippage,
	}, nil
}
