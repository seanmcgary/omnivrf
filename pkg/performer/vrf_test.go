package performer

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"go.uber.org/zap/zaptest"
)

// MockKeyManager for testing
type MockKeyManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewMockKeyManager() *MockKeyManager {
	// Generate a test private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return &MockKeyManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func (m *MockKeyManager) GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	return m.privateKey, nil
}

func (m *MockKeyManager) GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error) {
	return m.publicKey, nil
}

func TestVRFPerformer_ValidateTask(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	tests := []struct {
		name        string
		taskData    VRFTaskData
		expectError bool
	}{
		{
			name: "valid task",
			taskData: VRFTaskData{
				TaskHash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Seed:     big.NewInt(1),
			},
			expectError: false,
		},
		{
			name: "missing seed",
			taskData: VRFTaskData{
				TaskHash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Seed:     big.NewInt(0), // Use zero instead of nil for ABI encoding
			},
			expectError: true,
		},
		{
			name: "large seed (should pass - uint256 allows large values)",
			taskData: VRFTaskData{
				TaskHash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Seed:     new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)), // 2^256 - 1 (max uint256)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode task data using ABI
			data, err := performer.EncodeTaskData(tt.taskData)
			if err != nil {
				t.Fatalf("failed to encode task data: %v", err)
			}

			// Create task request
			req := &performerV1.TaskRequest{
				TaskId:  []byte("test-task-123"),
				Payload: data,
			}

			// Test validation
			err = performer.ValidateTask(req)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestVRFPerformer_HandleTask(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	// Create test task data
	taskData := VRFTaskData{
		TaskHash: common.HexToHash("0x456789abcdef456789abcdef456789abcdef456789abcdef456789abcdef456789"),
		Seed:     big.NewInt(42),
	}

	data, err := performer.EncodeTaskData(taskData)
	if err != nil {
		t.Fatalf("failed to encode task data: %v", err)
	}

	// Create task request
	req := &performerV1.TaskRequest{
		TaskId:  []byte("test-task-456"),
		Payload: data,
	}

	// Handle task
	resp, err := performer.HandleTask(req)
	if err != nil {
		t.Fatalf("failed to handle task: %v", err)
	}

	// Verify response
	if string(resp.TaskId) != string(req.TaskId) {
		t.Errorf("expected task ID %s, got %s", string(req.TaskId), string(resp.TaskId))
	}

	// Parse result
	result, err := performer.DecodeTaskResult(resp.Result)
	if err != nil {
		t.Fatalf("failed to decode result: %v", err)
	}

	// Verify result structure
	if result.TaskHash != common.BytesToHash(req.TaskId) {
		t.Errorf("expected task hash %s, got %s", common.BytesToHash(req.TaskId).Hex(), common.Hash(result.TaskHash).Hex())
	}

	if len(result.VRFProof) == 0 {
		t.Error("expected non-empty VRF proof")
	}

	if result.VRFOutput == nil || result.VRFOutput.Sign() == 0 {
		t.Error("expected non-zero VRF output")
	}

	// Verify VRF proof is correct
	publicKey, err := km.GetPublicKey(context.Background())
	if err != nil {
		t.Fatalf("failed to get public key: %v", err)
	}

	seedBytes := taskData.Seed.Bytes()

	verifiedOutput, err := performer.VerifyVRFProof(publicKey, seedBytes, result.VRFProof)
	if err != nil {
		t.Fatalf("failed to verify VRF proof: %v", err)
	}

	expectedOutput := new(big.Int).SetBytes(verifiedOutput)
	if result.VRFOutput.Cmp(expectedOutput) != 0 {
		t.Error("VRF proof verification failed - outputs don't match")
	}
}

func TestVRFPerformer_Deterministic(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	// Create test task data
	taskData := VRFTaskData{
		TaskHash: common.HexToHash("0x789abcdef789abcdef789abcdef789abcdef789abcdef789abcdef789abcdef789a"),
		Seed:     big.NewInt(100),
	}

	data, err := performer.EncodeTaskData(taskData)
	if err != nil {
		t.Fatalf("failed to encode task data: %v", err)
	}

	// Create task request
	req := &performerV1.TaskRequest{
		TaskId:  []byte("deterministic-test"),
		Payload: data,
	}

	// Handle task multiple times
	results := make([]*VRFTaskResult, 3)
	for i := 0; i < 3; i++ {
		resp, err := performer.HandleTask(req)
		if err != nil {
			t.Fatalf("failed to handle task iteration %d: %v", i, err)
		}

		results[i], err = performer.DecodeTaskResult(resp.Result)
		if err != nil {
			t.Fatalf("failed to decode result iteration %d: %v", i, err)
		}
	}

	// Verify all results are identical (VRF is deterministic)
	for i := 1; i < len(results); i++ {
		if string(results[0].VRFProof) != string(results[i].VRFProof) {
			t.Errorf("VRF proofs should be deterministic, but iteration %d differs", i)
		}
		if results[0].VRFOutput.Cmp(results[i].VRFOutput) != 0 {
			t.Errorf("VRF outputs should be deterministic, but iteration %d differs", i)
		}
	}
}

func TestVRFPerformer_DifferentSeeds(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	// Test different seeds produce different outputs
	seeds := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	outputs := make([]*big.Int, len(seeds))

	for i, seed := range seeds {
		taskData := VRFTaskData{
			TaskHash: common.HexToHash("0xabcdef123456abcdef123456abcdef123456abcdef123456abcdef123456abcdef"),
			Seed:     seed,
		}

		data, err := performer.EncodeTaskData(taskData)
		if err != nil {
			t.Fatalf("failed to encode task data for seed %s: %v", seed.String(), err)
		}

		req := &performerV1.TaskRequest{
			TaskId:  []byte("different-seeds-test"),
			Payload: data,
		}

		resp, err := performer.HandleTask(req)
		if err != nil {
			t.Fatalf("failed to handle task for seed %s: %v", seed.String(), err)
		}

		result, err := performer.DecodeTaskResult(resp.Result)
		if err != nil {
			t.Fatalf("failed to decode result for seed %s: %v", seed.String(), err)
		}

		outputs[i] = result.VRFOutput
	}

	// Verify all outputs are different
	for i := 0; i < len(outputs); i++ {
		for j := i + 1; j < len(outputs); j++ {
			if outputs[i].Cmp(outputs[j]) == 0 {
				t.Errorf("different seeds should produce different outputs, but outputs %d and %d are identical", i, j)
			}
		}
	}
}