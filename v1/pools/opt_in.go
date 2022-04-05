package pools

import (
	"context"

	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareLiquidityAssetOptInTransactions prepares liquidity asset opt-in transaction and returns a transaction group
func (p *Pool) PrepareLiquidityAssetOptInTransactions(ctx context.Context, userAddress string) (*utils.TransactionGroup, error) {
	if len(userAddress) == 0 {
		userAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.AssetOptInTransactions(p.LiquidityAsset.ID, userAddress, sp)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}
