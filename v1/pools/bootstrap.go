package pools

import (
	"context"

	"github.com/synycboom/tinyman-go-sdk/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/prepare"
)

// PrepareBootstrapTransactions prepares bootstrap transaction and returns a transaction group
func (p *Pool) PrepareBootstrapTransactions(ctx context.Context, bootstrapperAddress string) (*utils.TransactionGroup, error) {
	if len(bootstrapperAddress) == 0 {
		bootstrapperAddress = p.UserAddress
	}

	sp, err := p.ac.SuggestedParams().Do(ctx)
	if err != nil {
		return nil, err
	}

	txGroup, err := prepare.BootstrapTransactions(
		p.ValidatorAppID,
		p.Asset1.ID,
		p.Asset2.ID,
		p.Asset1.UnitName,
		p.Asset2.UnitName,
		bootstrapperAddress,
		sp,
	)
	if err != nil {
		return nil, err
	}

	return txGroup, nil
}
