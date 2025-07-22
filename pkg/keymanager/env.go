package keymanager

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

// EnvKeyManager implements KeyManager using environment variables
type EnvKeyManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewEnvKeyManager creates a new EnvKeyManager that reads the VRF private key from environment variables
func NewEnvKeyManager() (*EnvKeyManager, error) {
	keyHex := os.Getenv("VRF_PRIVATE_KEY")
	if keyHex == "" {
		return nil, errors.New("VRF_PRIVATE_KEY environment variable not set")
	}

	// Remove 0x prefix if present
	if len(keyHex) > 2 && keyHex[:2] == "0x" {
		keyHex = keyHex[2:]
	}

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, errors.New("invalid hex format for VRF_PRIVATE_KEY")
	}

	privateKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, errors.New("invalid private key format")
	}

	return &EnvKeyManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// GetPrivateKey returns the VRF private key
func (m *EnvKeyManager) GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	if m.privateKey == nil {
		return nil, errors.New("private key not initialized")
	}
	return m.privateKey, nil
}

// GetPublicKey returns the VRF public key
func (m *EnvKeyManager) GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error) {
	if m.publicKey == nil {
		return nil, errors.New("public key not initialized")
	}
	return m.publicKey, nil
}