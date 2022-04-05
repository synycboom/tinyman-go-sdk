package main

import (
	"context"
	"fmt"

	exampleUtils "github.com/synycboom/tinyman-go-sdk/example/utils"
	"github.com/synycboom/tinyman-go-sdk/v1/pools"
)

// This sample is provided for demonstration purposes only.

func main() {
	ctx := context.Background()
	account, err := exampleUtils.Account()
	if err != nil {
		panic(err)
	}

	userAddress := account.Address.String()
	ac, tc, err := exampleUtils.Clients(userAddress)
	if err != nil {
		panic(err)
	}

	// Check whether the user already opted in the app or not, if not let the user opt in
	if err := exampleUtils.OptInAppIfNeeded(ctx, tc, account); err != nil {
		panic(err)
	}

	// Create a new asset used to bootstraping one side of the liquidity pool
	assetID, err := exampleUtils.CreateAsset("sdk-test-2", "st2", 6, 1000000000, userAddress, account, ac, tc)
	if err != nil {
		panic(err)
	}

	token, err := tc.FetchAsset(ctx, assetID)
	if err != nil {
		panic(err)
	}

	algo, err := tc.FetchAsset(ctx, 0)
	if err != nil {
		panic(err)
	}

	// Fetch the created pool
	pool, err := pools.NewPool(ctx, ac, token, algo, nil, tc.ValidatorAppID, userAddress, true)
	if err != nil {
		panic(err)
	}

	// Prepare a transaction group for bootstrapping
	// Note that some transactions need to be signed with LogicSig account, and they were signed in the function.
	txGroup, err := pool.PrepareBootstrapTransactions(ctx, userAddress)
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

	fmt.Printf("Liquidity pool was bootstrapped with txid %s\n", txID)

	pool, err = tc.FetchPool(ctx, token, algo, true)
	if err != nil {
		panic(err)
	}

	// Get pool information
	info, err := pool.Info()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity pool info %#v\n", info)
}
