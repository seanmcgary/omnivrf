package performer

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

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
		taskData    TaskData
		expectError bool
	}{
		{
			name: "valid task",
			taskData: TaskData{
				RequestID: big.NewInt(1),
				Seed:      "deadbeef",
			},
			expectError: false,
		},
		{
			name: "missing request ID",
			taskData: TaskData{
				Seed: "deadbeef",
			},
			expectError: true,
		},
		{
			name: "missing seed",
			taskData: TaskData{
				RequestID: big.NewInt(1),
			},
			expectError: true,
		},
		{
			name: "invalid seed format",
			taskData: TaskData{
				RequestID: big.NewInt(1),
				Seed:      "not_hex",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal task data
			data, err := json.Marshal(tt.taskData)
			if err != nil {
				t.Fatalf("failed to marshal task data: %v", err)
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
	taskData := TaskData{
		RequestID: big.NewInt(42),
		Seed:      hex.EncodeToString([]byte("test seed")),
	}

	data, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("failed to marshal task data: %v", err)
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
	var result TaskResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	// Verify result structure
	if result.RequestID.Cmp(taskData.RequestID) != 0 {
		t.Errorf("expected request ID %s, got %s", taskData.RequestID.String(), result.RequestID.String())
	}

	if result.VRFProof == "" {
		t.Error("expected non-empty VRF proof")
	}

	if result.VRFOutput == "" {
		t.Error("expected non-empty VRF output")
	}

	// Verify VRF proof can be decoded
	proofBytes, err := hex.DecodeString(result.VRFProof)
	if err != nil {
		t.Errorf("failed to decode VRF proof: %v", err)
	}

	_, err = hex.DecodeString(result.VRFOutput)
	if err != nil {
		t.Errorf("failed to decode VRF output: %v", err)
	}

	// Verify VRF proof is correct
	publicKey, err := km.GetPublicKey(context.Background())
	if err != nil {
		t.Fatalf("failed to get public key: %v", err)
	}

	seedBytes, err := hex.DecodeString(taskData.Seed)
	if err != nil {
		t.Fatalf("failed to decode seed: %v", err)
	}

	verifiedOutput, err := performer.VerifyVRFProof(publicKey, seedBytes, proofBytes)
	if err != nil {
		t.Fatalf("failed to verify VRF proof: %v", err)
	}

	if hex.EncodeToString(verifiedOutput) != result.VRFOutput {
		t.Error("VRF proof verification failed - outputs don't match")
	}
}

func TestVRFPerformer_Deterministic(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	// Create test task data
	taskData := TaskData{
		RequestID: big.NewInt(100),
		Seed:      hex.EncodeToString([]byte("deterministic test")),
	}

	data, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("failed to marshal task data: %v", err)
	}

	// Create task request
	req := &performerV1.TaskRequest{
		TaskId:  []byte("deterministic-test"),
		Payload: data,
	}

	// Handle task multiple times
	results := make([]TaskResult, 3)
	for i := 0; i < 3; i++ {
		resp, err := performer.HandleTask(req)
		if err != nil {
			t.Fatalf("failed to handle task iteration %d: %v", i, err)
		}

		if err := json.Unmarshal(resp.Result, &results[i]); err != nil {
			t.Fatalf("failed to unmarshal result iteration %d: %v", i, err)
		}
	}

	// Verify all results are identical (VRF is deterministic)
	for i := 1; i < len(results); i++ {
		if results[0].VRFProof != results[i].VRFProof {
			t.Errorf("VRF proofs should be deterministic, but iteration %d differs", i)
		}
		if results[0].VRFOutput != results[i].VRFOutput {
			t.Errorf("VRF outputs should be deterministic, but iteration %d differs", i)
		}
	}
}

func TestVRFPerformer_DifferentSeeds(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	performer := NewVRFPerformer(km, logger)

	// Test different seeds produce different outputs
	seeds := []string{"seed1", "seed2", "seed3"}
	outputs := make([]string, len(seeds))

	for i, seed := range seeds {
		taskData := TaskData{
			RequestID: big.NewInt(int64(i + 1)),
			Seed:      hex.EncodeToString([]byte(seed)),
		}

		data, err := json.Marshal(taskData)
		if err != nil {
			t.Fatalf("failed to marshal task data for seed %s: %v", seed, err)
		}

		req := &performerV1.TaskRequest{
			TaskId:  []byte("different-seeds-test"),
			Payload: data,
		}

		resp, err := performer.HandleTask(req)
		if err != nil {
			t.Fatalf("failed to handle task for seed %s: %v", seed, err)
		}

		var result TaskResult
		if err := json.Unmarshal(resp.Result, &result); err != nil {
			t.Fatalf("failed to unmarshal result for seed %s: %v", seed, err)
		}

		outputs[i] = result.VRFOutput
	}

	// Verify all outputs are different
	for i := 0; i < len(outputs); i++ {
		for j := i + 1; j < len(outputs); j++ {
			if outputs[i] == outputs[j] {
				t.Errorf("different seeds should produce different outputs, but outputs %d and %d are identical", i, j)
			}
		}
	}
}