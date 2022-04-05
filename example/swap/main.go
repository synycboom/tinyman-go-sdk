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

	// Note that 1000000 is equal to 1000000 / 10 ** decimals (=6), which is 1ALGO
	// Get a quote for a swap of 1 ALGO to USDC with 5% slippage tolerance
	algoAmount, err := types.NewAssetAmount(algo, 1000000)
	if err != nil {
		panic(err)
	}

	// Fetch swap quote used when submit swapping transactions
	quote, err := pool.FetchFixedInputSwapQuote(ctx, algoAmount, 0.05)
	if err != nil {
		panic(err)
	}

	// Calculate price after applying the slippage
	priceWithSlippage, err := quote.PriceWithSlippage()
	if err != nil {
		panic(err)
	}

	// Calculate output amount after applying the slippage
	amountOutWithSlippage, err := quote.AmountOutWithSlippage()
	if err != nil {
		panic(err)
	}

	fmt.Printf("USDC per ALGO: %v\n", quote.Price())
	fmt.Printf("USDC per ALGO (worst case): %v\n", priceWithSlippage)
	fmt.Printf("Swapping %v to %v\n", quote.AmountIn.Amount, amountOutWithSlippage)

	// Prepare a transaction group for swappingg
	// Note that some transactions need to be signed with LogicSig account, and they were signed in the function.
	txGroup, err := pool.PrepareSwapTransactionsFromQuote(ctx, quote, userAddress)
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

	fmt.Printf("Swapped with txid %s\n", txID)
}
