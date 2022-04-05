package pools

import (
	"context"

	"github.com/synycboom/tinyman-go-sdk/types"
	"github.com/synycboom/tinyman-go-sdk/utils"
)

// FetchPoolPosition fetches pool position of a user
func (p *Pool) FetchPoolPosition(ctx context.Context, userAddress string) (*types.PoolPosition, error) {
	if len(userAddress) == 0 {
		userAddress = p.UserAddress
	}

	accountInfo, err := p.ac.AccountInformation(userAddress).Do(ctx)
	if err != nil {
		return nil, err
	}

	var liquidityAssetAmount types.AssetAmount
	for _, asset := range accountInfo.Assets {
		if asset.AssetId == p.LiquidityAsset.ID {
			liquidityAssetAmount = types.AssetAmount{
				Asset: &types.Asset{
					ID: asset.AssetId,
				},
				Amount: asset.Amount,
			}

			break
		}
	}

	quote, err := p.FetchBurnQuote(ctx, liquidityAssetAmount, 0)
	if err != nil {
		return nil, err
	}

	share, _ := utils.BigFloatDiv(utils.ToBigFloat(liquidityAssetAmount.Amount), utils.ToBigFloat(p.IssuedLiquidity)).Float64()

	return &types.PoolPosition{
		Asset1:         quote.AmountsOut[p.Asset1.ID],
		Asset2:         quote.AmountsOut[p.Asset2.ID],
		LiquidityAsset: quote.LiquidityAssetAmount,
		Share:          share,
	}, nil
}
