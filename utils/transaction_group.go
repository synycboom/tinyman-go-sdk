package utils

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/synycboom/tinyman-go-sdk/v1/constants"
)

// TransactionGroup is a group of transaction that can be executed atomically after signing
type TransactionGroup struct {
	transactions       []types.Transaction
	signedTransactions [][]byte
}

// NewTransactionGroup creates a new transaction group
func NewTransactionGroup(txs []types.Transaction) (*TransactionGroup, error) {
	gid, err := crypto.ComputeGroupID(txs)
	if err != nil {
		return nil, err
	}

	txsWithGroup := make([]types.Transaction, len(txs))
	for idx, tx := range txs {
		tx.Group = gid
		txsWithGroup[idx] = tx
	}

	return &TransactionGroup{
		transactions:       txsWithGroup,
		signedTransactions: make([][]byte, len(txs)),
	}, nil
}

// Sign signs a transaction group with an account
func (tg *TransactionGroup) Sign(acc *crypto.Account) error {
	accAddr := acc.Address.String()
	for idx, tx := range tg.transactions {
		if tx.Sender.String() == accAddr {
			_, stx, err := crypto.SignTransaction(acc.PrivateKey, tx)
			if err != nil {
				return err
			}

			tg.signedTransactions[idx] = stx
		}
	}

	return nil
}

// SignWithPrivateKey signs a transaction group with a given private key if a given address matches
func (tg *TransactionGroup) SignWithPrivateKey(address string, sk ed25519.PrivateKey) error {
	for idx, tx := range tg.transactions {
		if tx.Sender.String() == address {
			_, stx, err := crypto.SignTransaction(sk, tx)
			if err != nil {
				return err
			}

			tg.signedTransactions[idx] = stx
		}
	}

	return nil
}

// SignWithLogicSig signs a transaction group with logic sig account
func (tg *TransactionGroup) SignWithLogicSig(account *crypto.LogicSigAccount) error {
	address, err := account.Address()
	if err != nil {
		return err
	}

	for idx, tx := range tg.transactions {
		if tx.Sender.String() == address.String() {
			_, stx, err := crypto.SignLogicsigTransaction(account.Lsig, tx)
			if err != nil {
				return err
			}

			tg.signedTransactions[idx] = stx
		}
	}

	return nil
}

// Submit sends a signed transaction groups to the blockchain
func (tg *TransactionGroup) Submit(ctx context.Context, client *algod.Client, wait bool) (string, error) {
	var signedGroup []byte
	for _, signedTx := range tg.signedTransactions {
		signedGroup = append(signedGroup, signedTx...)
	}

	pendingTxID, err := client.SendRawTransaction(signedGroup).Do(ctx)
	if err != nil {
		return "", err
	}

	if wait {
		_, err := future.WaitForConfirmation(client, pendingTxID, constants.MaxWaitRound, ctx)
		if err != nil {
			return pendingTxID, err
		}
	}

	return pendingTxID, nil
}

// Transactions returns transactions inside the transaction group
func (tg *TransactionGroup) Transactions() []types.Transaction {
	return tg.transactions
}

// SignedTransactions returns signed transactions inside the transaction group
func (tg *TransactionGroup) SignedTransactions() [][]byte {
	return tg.signedTransactions
}

// SetSignedTransactions sets a signed transaction at a given index
func (tg *TransactionGroup) SetSignedTransactions(index int, signedTx []byte) error {
	if index > len(tg.signedTransactions)-1 {
		return fmt.Errorf("index is out of bound")
	}

	tg.signedTransactions[index] = signedTx

	return nil
}
