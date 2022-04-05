package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/synycboom/tinyman-go-sdk/v1"
	"github.com/synycboom/tinyman-go-sdk/v1/constants"
	"golang.org/x/crypto/ed25519"
)

// OptInAssetIfNeeded lets a user address opt in an asset if needed
func OptInAssetIfNeeded(ctx context.Context, tc *tinyman.Client, account *crypto.Account, assetID uint64) error {
	userAddress := account.Address.String()
	isOptedIn, err := tc.IsAssetOptedIn(ctx, assetID, userAddress)
	if err != nil {
		return err
	}
	if !isOptedIn {
		fmt.Printf("%s was not opted in for asset %v, opting in...\n", userAddress, assetID)
		txGroup, err := tc.PrepareAssetOptInTransactions(ctx, assetID, userAddress)
		if err != nil {
			return err
		}

		if err := txGroup.Sign(account); err != nil {
			return err
		}

		txID, err := tc.Submit(ctx, txGroup, true)
		if err != nil {
			return err
		}

		fmt.Printf("Opted in for asset %v with txid %s\n", assetID, txID)
	}

	return nil
}

// OptInAppIfNeeded lets a user address opt in application if needed
func OptInAppIfNeeded(ctx context.Context, tinymanCli *tinyman.Client, account *crypto.Account) error {
	// userAddress param for IsOptedIn() is optional if it was set when initializing the client
	isOptedIn, err := tinymanCli.IsOptedIn(ctx, "")
	if err != nil {
		return err
	}

	if !isOptedIn {
		fmt.Println("Account is not opted into app, opting in now...")
		txGroup, err := tinymanCli.PrepareAppOptInTransaction(ctx, "")
		if err != nil {
			return err
		}

		if err := txGroup.Sign(account); err != nil {
			return err
		}

		txID, err := tinymanCli.Submit(ctx, txGroup, true)
		if err != nil {
			return err
		}

		fmt.Printf("Submitted opt-in tx %s", txID)
	}

	return nil
}

// Clients returns Algorand and Tinyman clients
func Clients(userAddress string) (*algod.Client, *tinyman.Client, error) {
	algodCli, err := algod.MakeClient(constants.AlgodTestnetHost, "")
	if err != nil {
		return nil, nil, err
	}

	client, err := tinyman.NewTestNetClient(algodCli, userAddress)
	if err != nil {
		return nil, nil, err
	}

	return algodCli, client, nil
}

// Account returns account from mnemonic or base64 private key from os environments
func Account() (*crypto.Account, error) {
	var privateKey ed25519.PrivateKey
	base64PrivateKey := os.Getenv("BASE64_PRIVATE_KEY")
	mnemonicWords := os.Getenv("MNEMONIC")
	if len(mnemonicWords) > 0 {
		key, err := mnemonic.ToPrivateKey(mnemonicWords)
		if err != nil {
			return nil, err
		}

		privateKey = key
	} else {
		decodedKey, err := base64.StdEncoding.DecodeString(base64PrivateKey)
		if err != nil {
			return nil, err
		}

		privateKey = decodedKey
	}

	account, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
