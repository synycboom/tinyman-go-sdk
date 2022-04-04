package utils

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/crypto"
	algoTypes "github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/ed25519"

	"github.com/synycboom/tinyman-go-sdk/types"
)

// TransactionGroup is a group of transaction that can be executed atomically after signing
type TransactionGroup struct {
	transactions       []algoTypes.Transaction
	signedTransactions [][]byte
}

// NewTransactionGroup creates a new transaction group
func NewTransactionGroup(itxs []any) (*TransactionGroup, error) {
	txs := make([]algoTypes.Transaction, len(itxs))
	for idx, itx := range itxs {
		tx, ok := itx.(algoTypes.Transaction)
		if !ok {
			return nil, fmt.Errorf("wrong transaction type")
		}

		txs[idx] = tx
	}
	gid, err := crypto.ComputeGroupID(txs)
	if err != nil {
		return nil, err
	}

	txsWithGroup := make([]algoTypes.Transaction, len(txs))
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
func (tg *TransactionGroup) Sign(acc types.CryptoAccount) error {
	address := acc.Address()
	privateKey := ed25519.PrivateKey(acc.PrivateKey())
	for idx, tx := range tg.transactions {
		if tx.Sender.String() == address {
			_, stx, err := crypto.SignTransaction(privateKey, tx)
			if err != nil {
				return err
			}

			tg.signedTransactions[idx] = stx
		}
	}

	return nil
}

// SignWithLogicSig signs a transaction group with logic sig account
func (tg *TransactionGroup) SignWithLogicSig(acc types.CryptoLogicSigAccount) error {
	address, err := acc.Address()
	if err != nil {
		return err
	}

	logicSig, ok := acc.LogicSig().(algoTypes.LogicSig)
	if !ok {
		return fmt.Errorf("wrong logic signature type")
	}

	for idx, tx := range tg.transactions {
		if tx.Sender.String() == address {
			_, stx, err := crypto.SignLogicsigTransaction(logicSig, tx)
			if err != nil {
				return err
			}

			tg.signedTransactions[idx] = stx
		}
	}

	return nil
}

// Submit sends a signed transaction groups to the blockchain
func (tg *TransactionGroup) Submit(ctx context.Context, client types.AlgoClient, wait bool) (string, error) {
	return client.SendRawTransaction(ctx, tg.ComputeSignedGroup(), wait)
}

// ComputeSignedGroup calculate a signed group
func (tg *TransactionGroup) ComputeSignedGroup() []byte {
	var signedGroup []byte
	for _, signedTx := range tg.signedTransactions {
		signedGroup = append(signedGroup, signedTx...)
	}

	return signedGroup
}
