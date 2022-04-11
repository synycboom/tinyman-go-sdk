package pools

import (
	"context"
	"fmt"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareMintTransactions prepares mint transaction and returns a transaction group
func (p *Pool) PrepareMintTransactions(
	ctx context.Context,
	amountsIn map[uint64]types.AssetAmount,
	liquidityAssetAmount *types.AssetAmount,
	minterAddress string,
) (*utils.TransactionGroup, error) {
	if liquidityAssetAmount == nil {
		return nil, fmt.Errorf("liquidityAssetAmount is required")
	}

	asset1Amount := amountsIn[p.Asset1.ID]
	asset2Amount := amountsIn[p.Asset2.ID]
	if len(minterAddress) == 0 {
		minterAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.MintTransactions(
		p.ValidatorAppID,
		p.Asset1.ID,
		p.Asset2.ID,
		p.LiquidityAsset.ID,
		asset1Amount.Amount,
		asset2Amount.Amount,
		liquidityAssetAmount.Amount,
		minterAddress,
		sp,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// PrepareMintTransactionsFromQuote prepares mint transaction from a given mint quote and returns a transaction group
func (p *Pool) PrepareMintTransactionsFromQuote(ctx context.Context, quote *types.MintQuote, minterAddress string) (*utils.TransactionGroup, error) {
	if quote == nil {
		return nil, fmt.Errorf("quote is required")
	}

	return p.PrepareMintTransactions(
		ctx,
		quote.AmountsIn,
		&quote.LiquidityAssetAmount,
		minterAddress,
	)
}
