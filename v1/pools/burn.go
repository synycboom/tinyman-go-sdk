package pools

import (
	"context"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareBurnTransactions prepares burn transaction and returns a transaction group
func (p *Pool) PrepareBurnTransactions(
	ctx context.Context,
	assetsOut map[uint64]types.AssetAmount,
	liquidityAssetAmount types.AssetAmount,
	burnerAddress string,
) (*utils.TransactionGroup, error) {
	asset1Amount := assetsOut[p.Asset1.ID]
	asset2Amount := assetsOut[p.Asset2.ID]
	if len(burnerAddress) == 0 {
		burnerAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.BurnTransactions(
		p.ValidatorAppID,
		p.Asset1.ID,
		p.Asset2.ID,
		p.LiquidityAsset.ID,
		asset1Amount.Amount,
		asset2Amount.Amount,
		liquidityAssetAmount.Amount,
		burnerAddress,
		sp,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// PrepareBurnTransactionsFromQuote prepares burn transaction from a given burn quote and returns a transaction group
func (p *Pool) PrepareBurnTransactionsFromQuote(ctx context.Context, quote types.BurnQuote, burnerAddress string) (*utils.TransactionGroup, error) {
	amountsOut, err := quote.AmountsOutWithSlippage()
	if err != nil {
		return nil, err
	}

	return p.PrepareBurnTransactions(
		ctx,
		amountsOut,
		quote.LiquidityAssetAmount,
		burnerAddress,
	)
}
