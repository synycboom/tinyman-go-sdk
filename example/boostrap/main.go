package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"golang.org/x/crypto/ed25519"

	"github.com/synycboom/tinyman-go-sdk/v1"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"github.com/synycboom/tinyman-go-sdk/v1/pools"
)

// This sample is provided for demonstration purposes only.

func main() {
	var privateKey ed25519.PrivateKey
	base64PrivateKey := os.Getenv("BASE64_PRIVATE_KEY")
	mnemonicWords := os.Getenv("MNEMONIC")
	if len(mnemonicWords) > 0 {
		key, err := mnemonic.ToPrivateKey(mnemonicWords)
		if err != nil {
			panic(err)
		}

		privateKey = key
	} else {
		decodedKey, err := base64.StdEncoding.DecodeString(base64PrivateKey)
		if err != nil {
			panic(err)
		}

		privateKey = decodedKey
	}

	account, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	userAddress := account.Address.String()
	algodCli, err := algod.MakeClient(constants.AlgodTestnetHost, "")
	if err != nil {
		panic(err)
	}

	client, err := tinyman.NewTestNetClient(algodCli, userAddress)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	isOptedIn, err := client.IsOptedIn(ctx, userAddress)
	if err != nil {
		panic(err)
	}

	if !isOptedIn {
		fmt.Println("Account is not opted into app, opting in now...")
		txGroup, err := client.PrepareAppOptInTransaction(ctx, userAddress)
		if err != nil {
			panic(err)
		}

		if err := txGroup.Sign(&account); err != nil {
			panic(err)
		}

		txID, err := client.Submit(ctx, txGroup, true)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Submitted opt-in tx %s", txID)
	}

	assetID, err := createAsset("sdk-test-2", "st2", 6, userAddress, account, algodCli, client)
	if err != nil {
		panic(err)
	}

	token, err := client.FetchAsset(ctx, assetID)
	if err != nil {
		panic(err)
	}

	algo, err := client.FetchAsset(ctx, 0)
	if err != nil {
		panic(err)
	}

	pool, err := pools.NewPool(ctx, algodCli, token, algo, nil, client.ValidatorAppID, userAddress, true)
	if err != nil {
		panic(err)
	}

	txGroup, err := pool.PrepareBootstrapTransactions(ctx, userAddress)
	if err != nil {
		panic(err)
	}

	if err := txGroup.Sign(&account); err != nil {
		panic(err)
	}

	txID, err := client.Submit(ctx, txGroup, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity pool was bootstrapped with txid %s\n", txID)

	pool, err = client.FetchPool(ctx, token, algo, true)
	if err != nil {
		panic(err)
	}

	info, err := pool.Info()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Liquidity pool info %#v\n", info)
}
