package main

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"

	"github.com/synycboom/tinyman-go-sdk/v1"
)

func createAsset(
	assetName,
	unitName string,
	decimals uint32,
	userAddress string,
	account crypto.Account,
	algodClient *algod.Client,
	tinymanClient *tinyman.Client,
) (uint64, error) {
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return 0, err
	}

	var note []byte = nil
	addr := userAddress
	defaultFrozen := false
	totalIssuance := uint64(1000000000000)
	reserve := addr
	freeze := addr
	clawback := addr
	manager := addr
	assetURL := "http://someurl"
	assetMetadataHash := "thisIsSomeLength32HashCommitment"

	txn, err := future.MakeAssetCreateTxn(addr, note, txParams,
		totalIssuance, decimals, defaultFrozen, manager, reserve, freeze, clawback,
		unitName, assetName, assetURL, assetMetadataHash,
	)
	if err != nil {
		return 0, err
	}

	txid, stx, err := crypto.SignTransaction(account.PrivateKey, txn)
	if err != nil {
		return 0, err
	}

	if _, err := algodClient.SendRawTransaction(stx).Do(context.Background()); err != nil {
		return 0, err
	}

	confirmedTxn, err := future.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	isTokenOptedIn, err := tinymanClient.IsAssetOptedIn(ctx, confirmedTxn.AssetIndex, userAddress)
	if err != nil {
		panic(err)
	}
	if !isTokenOptedIn {
		fmt.Printf("%s was not opted in for asset %v, opting in...\n", userAddress, confirmedTxn.AssetIndex)
		txGroup, err := tinymanClient.PrepareAssetOptInTransactions(ctx, confirmedTxn.AssetIndex, userAddress)
		if err != nil {
			panic(err)
		}

		if err := txGroup.Sign(&account); err != nil {
			panic(err)
		}

		txID, err := tinymanClient.Submit(ctx, txGroup, true)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Opted in for asset %v with txid %s\n", confirmedTxn.AssetIndex, txID)
	}

	return confirmedTxn.AssetIndex, nil
}
