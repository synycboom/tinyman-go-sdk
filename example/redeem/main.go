package main

import (
	"context"
	"fmt"

	exampleUtils "github.com/synycboom/tinyman-go-sdk/example/utils"
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

	// Fetch excess amount resulting from the swap
	redeemQuotes, err := tc.FetchExcessAmount(ctx, userAddress)
	if err != nil {
		panic(err)
	}

	// Filter redeem quotes matching with USDC
	redeemQuote, err := pool.GetRedeemQuoteMatchesAssetID(usdc.ID, redeemQuotes)
	if err != nil {
		panic(err)
	}

	if redeemQuote == nil {
		fmt.Println("No excess amount to be redeemed")

		return
	}

	fmt.Printf("There is %s to be redeemed\n", redeemQuote.Amount.String())

	// Prepare a transaction group for redeeming
	// Note that some transactions need to be signed with LogicSig account, and they were signed in the function.
	txGroup, err := pool.PrepareRedeemTransactionsFromQuote(ctx, redeemQuote, userAddress)
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

	fmt.Printf("Redeemed excess amount with txid %s\n", txID)
}
