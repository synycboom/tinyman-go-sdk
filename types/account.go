package types

// CryptoAccount represents a cryptography account
type CryptoAccount interface {
	// Address returns an address derived from the private key
	Address() string

	// PublicKey returns a public key derived from the private key
	PublicKey() []byte

	// PrivateKey returns the private key
	PrivateKey() []byte
}

// CryptoLogicSigAccount is a logic signature account
type CryptoLogicSigAccount interface {
	// Address returns an address of the logic signature account
	Address() (string, error)

	// LogicSig returns the underlying logic signature which needs type casting
	LogicSig() interface{}
}
