package main

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

	// Create valid VRF task data using ABI encoding
	taskHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	seed := new(big.Int)
	seed.SetString("12345678901234567890", 10)
	
	taskData := performer.VRFTaskData{
		TaskHash: taskHash,
		Seed:     seed,
	}

	// Create a temporary performer instance to access encoding methods
	tempPerformer := &performer.VRFPerformer{}
	payload, err := tempPerformer.EncodeTaskData(taskData)
	if err != nil {
		t.Fatalf("Failed to encode task data: %v", err)
	}

	// Use the task hash as the TaskId (as TaskMailbox would provide)
	taskRequest := &performerV1.TaskRequest{
		TaskId:   taskHash.Bytes(),
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

	// Verify result structure using ABI decoding
	result, err := tempPerformer.DecodeTaskResult(resp.Result)
	if err != nil {
		t.Fatalf("Failed to decode result: %v", err)
	}

	if result.TaskHash != taskHash {
		t.Errorf("Expected task hash %s, got %s", taskHash.Hex(), common.BytesToHash(result.TaskHash[:]).Hex())
	}

	if len(result.VRFProof) == 0 {
		t.Error("Expected non-empty VRF proof")
	}

	if result.VRFOutput == nil || result.VRFOutput.Sign() <= 0 {
		t.Error("Expected non-zero VRF output")
	}

	t.Logf("VRF task completed successfully - Task Hash: %s, VRF Output: %s", 
		common.BytesToHash(result.TaskHash[:]).Hex(), result.VRFOutput.String())
}
