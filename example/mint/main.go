package main

import (
	"context"
	"fmt"

	exampleUtils "github.com/synycboom/tinyman-go-sdk/example/utils"
	"github.com/synycboom/tinyman-go-sdk/types"
)

// This sample is provided for demonstration purposes only.

func main() {
	ctx := context.Background()
	account, err := exampleUtils.Account()
	if err != nil {
		panic(err)
	}

	userAddress := account.Address.String()
	_, tc, err := exampleUtils.Clients(userAddress)
	if err != nil {
		panic(err)
	}

	// Check whether the user already opted in the app or not, if not let the user opt in
	if err := exampleUtils.OptInAppIfNeeded(ctx, tc, account); err != nil {
		panic(err)
	}

	usdc, err := tc.FetchAsset(ctx, 10458941)
	if err != nil {
		panic(err)
	}

	algo, err := tc.FetchAsset(ctx, 0)
	if err != nil {
		panic(err)
	}

	// Check whether the user already opted in USDC or not, if not let the user opt in
	if err := exampleUtils.OptInAssetIfNeeded(ctx, tc, account, usdc.ID); err != nil {
		panic(err)
	}

	// Fetch USDC-ALGO pool
	pool, err := tc.FetchPool(ctx, usdc, algo, true)
	if err != nil {
		panic(err)
	}

	// Check whether the user already opted in the liquidity pool or not, if not let the user opt in
	if err := exampleUtils.OptInAssetIfNeeded(ctx, tc, account, pool.LiquidityAsset.ID); err != nil {
		panic(err)
	}

	// Fetch mint quote used when submit minting transactions
	// Note that 500000 is equal to 500000 / 10 ** decimals (=6), which is 0.5USDC
	usdcAssetAmount, _ := types.NewAssetAmount(usdc, 500000)
	quote, err := pool.FetchMintQuote(ctx, usdcAssetAmount, nil, 0.05)
	if err != nil {
		panic(err)
	}

	// Calculate the liquidity asset amount after applying the slippage
	liquidityAssetAmountWithSlippage, err := quote.LiquidityAssetAmountWithSlippage()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity asset amount: %s\n", quote.LiquidityAssetAmount.String())
	fmt.Printf("Liquidity asset amount with slippage: %s\n", liquidityAssetAmountWithSlippage.String())

	// Prepare a transaction group for minting
	// Note that some transactions need to be signed with LogicSig account, and they were signed in the function.
	txGroup, err := pool.PrepareMintTransactionsFromQuote(ctx, quote, userAddress)
	if err != nil {
		panic(err)
	}

	// Some transactions that need the user signatures are signed here
	if err := txGroup.Sign(account); err != nil {
		panic(err)
	}

	// Submit a group of transaction to the blockchain
	txID, err := tc.Submit(ctx, txGroup, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity was added in txid %s\n", txID)

	// Fetch the pool position of the user
	info, err := pool.FetchPoolPosition(ctx, userAddress)
	if err != nil {
		panic(err)
	}

	share := info.Share * 100
	fmt.Printf("Pool Tokens: %s\n", info.LiquidityAsset.String())
	fmt.Printf("Asset1: %s\n", info.Asset1.String())
	fmt.Printf("Asset2: %s\n", info.Asset2.String())
	fmt.Printf("Share of pool: %f%%\n", share)
}
