package pools

import (
	"context"
	"fmt"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareSwapTransactions prepares swap transaction and returns a transaction group
func (p *Pool) PrepareSwapTransactions(
	ctx context.Context,
	assetAmountIn,
	assetAmountOut *types.AssetAmount,
	swapType,
	swapperAddress string,
) (*utils.TransactionGroup, error) {
	if assetAmountIn == nil || assetAmountOut == nil {
		return nil, fmt.Errorf("assetAmountIn and assetAmountOut are required")
	}
	if len(swapperAddress) == 0 {
		swapperAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.SwapTransactions(
		p.ValidatorAppID,
		p.Asset1.ID,
		p.Asset2.ID,
		p.LiquidityAsset.ID,
		assetAmountIn.Asset.ID,
		assetAmountIn.Amount,
		assetAmountOut.Amount,
		swapType,
		swapperAddress,
		sp,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// PrepareSwapTransactionsFromQuote prepares swap transaction from a given swap quote and returns a transaction group
func (p *Pool) PrepareSwapTransactionsFromQuote(ctx context.Context, quote *types.SwapQuote, swapperAddress string) (*utils.TransactionGroup, error) {
	if quote == nil {
		return nil, fmt.Errorf("quote is required")
	}
	amountOut, err := quote.AmountOutWithSlippage()
	if err != nil {
		return nil, err
	}

	txGroup, err := p.PrepareSwapTransactions(
		ctx,
		quote.AmountIn,
		amountOut,
		quote.SwapType,
		swapperAddress,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}
