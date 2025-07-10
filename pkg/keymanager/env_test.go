package keymanager

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewEnvKeyManager(t *testing.T) {
	// Generate a test private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	// Convert to hex
	privateKeyBytes := crypto.FromECDSA(privateKey)
	hexKey := hex.EncodeToString(privateKeyBytes)

	tests := []struct {
		name        string
		envValue    string
		expectError bool
	}{
		{
			name:        "valid key without 0x prefix",
			envValue:    hexKey,
			expectError: false,
		},
		{
			name:        "valid key with 0x prefix",
			envValue:    "0x" + hexKey,
			expectError: false,
		},
		{
			name:        "empty environment variable",
			envValue:    "",
			expectError: true,
		},
		{
			name:        "invalid hex format",
			envValue:    "not_hex_data",
			expectError: true,
		},
		{
			name:        "invalid key length",
			envValue:    "deadbeef",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("VRF_PRIVATE_KEY", tt.envValue)
			} else {
				os.Unsetenv("VRF_PRIVATE_KEY")
			}

			// Create key manager
			km, err := NewEnvKeyManager()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Test getting private key
			ctx := context.Background()
			privKey, err := km.GetPrivateKey(ctx)
			if err != nil {
				t.Errorf("failed to get private key: %v", err)
				return
			}

			if privKey == nil {
				t.Error("expected non-nil private key")
				return
			}

			// Test getting public key
			pubKey, err := km.GetPublicKey(ctx)
			if err != nil {
				t.Errorf("failed to get public key: %v", err)
				return
			}

			if pubKey == nil {
				t.Error("expected non-nil public key")
				return
			}

			// Verify public key matches private key
			expectedPubKey := &privKey.PublicKey
			if pubKey.X.Cmp(expectedPubKey.X) != 0 || pubKey.Y.Cmp(expectedPubKey.Y) != 0 {
				t.Error("public key doesn't match private key")
			}
		})
	}

	// Clean up
	os.Unsetenv("VRF_PRIVATE_KEY")
}

func TestEnvKeyManager_KeyConsistency(t *testing.T) {
	// Generate a test private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	// Convert to hex and set environment variable
	privateKeyBytes := crypto.FromECDSA(privateKey)
	hexKey := hex.EncodeToString(privateKeyBytes)
	os.Setenv("VRF_PRIVATE_KEY", hexKey)
	defer os.Unsetenv("VRF_PRIVATE_KEY")

	// Create key manager
	km, err := NewEnvKeyManager()
	if err != nil {
		t.Fatalf("failed to create key manager: %v", err)
	}

	ctx := context.Background()

	// Get keys multiple times and verify consistency
	for i := 0; i < 5; i++ {
		privKey, err := km.GetPrivateKey(ctx)
		if err != nil {
			t.Errorf("iteration %d: failed to get private key: %v", i, err)
			continue
		}

		pubKey, err := km.GetPublicKey(ctx)
		if err != nil {
			t.Errorf("iteration %d: failed to get public key: %v", i, err)
			continue
		}

		// Verify the key is the same as the original
		if privKey.D.Cmp(privateKey.D) != 0 {
			t.Errorf("iteration %d: private key doesn't match original", i)
		}

		// Verify public key matches private key
		expectedPubKey := &privKey.PublicKey
		if pubKey.X.Cmp(expectedPubKey.X) != 0 || pubKey.Y.Cmp(expectedPubKey.Y) != 0 {
			t.Errorf("iteration %d: public key doesn't match private key", i)
		}
	}
}

func TestEnvKeyManager_SignVerify(t *testing.T) {
	// Generate a test private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	// Convert to hex and set environment variable
	privateKeyBytes := crypto.FromECDSA(privateKey)
	hexKey := hex.EncodeToString(privateKeyBytes)
	os.Setenv("VRF_PRIVATE_KEY", hexKey)
	defer os.Unsetenv("VRF_PRIVATE_KEY")

	// Create key manager
	km, err := NewEnvKeyManager()
	if err != nil {
		t.Fatalf("failed to create key manager: %v", err)
	}

	ctx := context.Background()

	// Get keys
	privKey, err := km.GetPrivateKey(ctx)
	if err != nil {
		t.Fatalf("failed to get private key: %v", err)
	}

	pubKey, err := km.GetPublicKey(ctx)
	if err != nil {
		t.Fatalf("failed to get public key: %v", err)
	}

	// Test signing and verification
	message := []byte("test message")
	hash := crypto.Keccak256Hash(message)

	// Sign with private key
	signature, err := crypto.Sign(hash.Bytes(), privKey)
	if err != nil {
		t.Fatalf("failed to sign message: %v", err)
	}

	// Verify with public key
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		t.Fatalf("failed to recover public key: %v", err)
	}

	// Convert to ecdsa.PublicKey for comparison
	recoveredPubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		t.Fatalf("failed to unmarshal recovered public key: %v", err)
	}

	// Verify the recovered public key matches our public key
	if pubKey.X.Cmp(recoveredPubKey.X) != 0 || pubKey.Y.Cmp(recoveredPubKey.Y) != 0 {
		t.Error("signature verification failed - public keys don't match")
	}
}