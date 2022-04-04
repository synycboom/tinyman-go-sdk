package types

import "context"

// AlgoClient represents the Algorand client
type AlgoClient interface {
	// SendRawTransaction submit a signed transaction group and wait for it to be confirmed if wait is true
	SendRawTransaction(ctx context.Context, signedGroup []byte, wait bool) (string, error)
}
