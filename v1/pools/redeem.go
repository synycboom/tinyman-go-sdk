package pools

import (
	"context"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareRedeemTransactions prepares redeem transaction and returns a transaction group
func (p *Pool) PrepareRedeemTransactions(ctx context.Context, amountOut types.AssetAmount, redeemerAddress string) (*utils.TransactionGroup, error) {
	if len(redeemerAddress) == 0 {
		redeemerAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.RedeemTransactions(
		p.ValidatorAppID,
		p.Asset1.ID,
		p.Asset2.ID,
		p.LiquidityAsset.ID,
		amountOut.Asset.ID,
		amountOut.Amount,
		redeemerAddress,
		sp,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}

// PrepareRedeemTransactionsFromQuote prepares redeem transactions and return a transaction group from quote
func (p *Pool) PrepareRedeemTransactionsFromQuote(ctx context.Context, quote types.RedeemQuote, redeemerAddress string) (*utils.TransactionGroup, error) {
	return p.PrepareRedeemTransactions(ctx, quote.Amount, redeemerAddress)
}

// FilterRedeemQuotes filters redeem quotes belonging to this pool
func (p *Pool) FilterRedeemQuotes(quotes []types.RedeemQuote) ([]types.RedeemQuote, error) {
	poolAddress, err := p.Address()
	if err != nil {
		return nil, err
	}

	var res []types.RedeemQuote
	for _, r := range quotes {
		if r.PoolAddress == poolAddress {
			res = append(res, r)
		}
	}

	return res, nil
}

// GetRedeemQuoteMatchesAssetID filters redeem quote belonging to this pool which matches an asset id
func (p *Pool) GetRedeemQuoteMatchesAssetID(assetID uint64, quotes []types.RedeemQuote) (*types.RedeemQuote, error) {
	quotes, err := p.FilterRedeemQuotes(quotes)
	if err != nil {
		return nil, err
	}

	for _, r := range quotes {
		if r.Amount.Asset.ID == assetID {
			return &r, nil
		}
	}

	return nil, nil
}
