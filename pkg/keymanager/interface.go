package keymanager

import (
	"context"
	"crypto/ecdsa"
)

// KeyManager provides an interface for securely managing VRF private keys
type KeyManager interface {
	// GetPrivateKey returns the VRF private key for signing
	GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error)
	
	// GetPublicKey returns the corresponding public key
	GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error)
}