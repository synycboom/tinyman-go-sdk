package main

import (
	"context"
	"fmt"

	exampleUtils "github.com/synycboom/tinyman-go-sdk/example/utils"
	"github.com/synycboom/tinyman-go-sdk/types"
)

// This sample is provided for demonstration purposes only.
// It is not intended for production use.
// This example does not constitute trading advice.

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

	pool, err := tc.FetchPool(ctx, usdc, algo, true)
	if err != nil {
		panic(err)
	}

	// Check whether the user already opted in the liquidity pool or not, if not let the user opt in
	if err := exampleUtils.OptInAssetIfNeeded(ctx, tc, account, pool.LiquidityAsset.ID); err != nil {
		panic(err)
	}

	// Fetch total balance of liquidtiy asset that the user has
	balance, err := tc.Balance(ctx, pool.LiquidityAsset, userAddress)
	if err != nil {
		panic(err)
	}

	if balance.Amount == 0 {
		fmt.Println("User does not have liquidity balance")

		return
	}

	fmt.Printf("Current balance of liquidity \n\t - ID:%v = %v\n", pool.LiquidityAsset.ID, balance.Amount)

	// Fetch burn quote used when buring liquidity asset
	quote, err := pool.FetchBurnQuote(ctx, balance, 0.05)
	if err != nil {
		panic(err)
	}

	// Calculate the output assets after applying the slippage
	outputAmountsWithSlippage, err := quote.AmountsOutWithSlippage()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity asset amount \n\t - %v\n", quote.LiquidityAssetAmount.String())
	fmt.Println("Output amounts")
	displayAmountOut(quote.AmountsOut)
	fmt.Println("Output amounts with slippage")
	displayAmountOut(outputAmountsWithSlippage)

	// Prepare a transaction group for burning
	// Note that some transactions need to be signed with LogicSig account, and they were signed in the function.
	txGroup, err := pool.PrepareBurnTransactionsFromQuote(ctx, quote, userAddress)
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

	fmt.Printf("Liquidity was removed in txid %s\n", txID)

	// Fetch pool position of the user
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

func displayAmountOut(out map[uint64]types.AssetAmount) {
	for _, v := range out {
		fmt.Printf("\t - %v\n", v.String())
	}
}
