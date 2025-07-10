package main

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

	"github.com/Layr-Labs/hourglass-avs-template/pkg/performer"
)

// MockKeyManager for testing
type MockKeyManager struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewMockKeyManager() *MockKeyManager {
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

func Test_VRFTaskRequestPayload(t *testing.T) {
	logger := zaptest.NewLogger(t)
	km := NewMockKeyManager()
	vrfPerformer := performer.NewVRFPerformer(km, logger)

	// Create valid VRF task data
	taskData := performer.TaskData{
		RequestID: big.NewInt(42),
		Seed:      hex.EncodeToString([]byte("test seed for integration")),
	}

	payload, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("Failed to marshal task data: %v", err)
	}

	taskRequest := &performerV1.TaskRequest{
		TaskId:   []byte("integration-test-task-id"),
		Payload:  payload,
		Metadata: []byte("test-metadata"),
	}

	// Test validation
	err = vrfPerformer.ValidateTask(taskRequest)
	if err != nil {
		t.Errorf("ValidateTask failed: %v", err)
	}

	// Test task handling
	resp, err := vrfPerformer.HandleTask(taskRequest)
	if err != nil {
		t.Errorf("HandleTask failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if string(resp.TaskId) != string(taskRequest.TaskId) {
		t.Errorf("Expected task ID %s, got %s", string(taskRequest.TaskId), string(resp.TaskId))
	}

	// Verify result structure
	var result performer.TaskResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result.RequestID.Cmp(taskData.RequestID) != 0 {
		t.Errorf("Expected request ID %s, got %s", taskData.RequestID.String(), result.RequestID.String())
	}

	if result.VRFProof == "" {
		t.Error("Expected non-empty VRF proof")
	}

	if result.VRFOutput == "" {
		t.Error("Expected non-empty VRF output")
	}

	t.Logf("VRF task completed successfully - Request ID: %s, VRF Output: %s", 
		result.RequestID.String(), result.VRFOutput)
}
