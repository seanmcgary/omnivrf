package performer

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vechain/go-ecvrf"
	performerV1 "github.com/Layr-Labs/protocol-apis/gen/protos/eigenlayer/hourglass/v1/performer"
	"go.uber.org/zap"
)

// KeyManager interface for managing VRF keys
type KeyManager interface {
	GetPrivateKey(ctx context.Context) (*ecdsa.PrivateKey, error)
	GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error)
}

// VRFPerformer implements the Hourglass Performer interface for VRF computation
type VRFPerformer struct {
	keyManager KeyManager
	suite      ecvrf.VRF
	logger     *zap.Logger
}

// VRFTaskData represents the ABI-encoded structure of incoming VRF task data (matches Solidity struct)
type VRFTaskData struct {
	TaskHash [32]byte `abi:"taskHash"`
	Seed     *big.Int `abi:"seed"`
}

// VRFTaskResult represents the VRF computation result
type VRFTaskResult struct {
	TaskHash  [32]byte `abi:"taskHash"`
	VRFProof  []byte   `abi:"vrfProof"`
	VRFOutput *big.Int `abi:"vrfOutput"`
}

// NewVRFPerformer creates a new VRF performer instance
func NewVRFPerformer(keyManager KeyManager, logger *zap.Logger) *VRFPerformer {
	return &VRFPerformer{
		keyManager: keyManager,
		suite:      ecvrf.Secp256k1Sha256Tai, // Ethereum compatible VRF suite
		logger:     logger,
	}
}

// ValidateTask validates the incoming VRF task request
func (p *VRFPerformer) ValidateTask(t *performerV1.TaskRequest) error {
	p.logger.Sugar().Infow("Validating VRF task",
		zap.String("taskId", string(t.TaskId)),
		zap.Int("payloadLength", len(t.Payload)),
	)

	// Parse and validate ABI-encoded task data
	taskData, err := p.decodeTaskData(t.Payload)
	if err != nil {
		p.logger.Sugar().Errorw("Failed to decode task data",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return fmt.Errorf("invalid task data format: %w", err)
	}

	// Validate required fields
	if taskData.Seed == nil {
		return fmt.Errorf("missing seed in task data")
	}

	if taskData.Seed.Sign() <= 0 {
		return fmt.Errorf("invalid seed: must be positive")
	}

	// Verify we can access the VRF key
	if _, err := p.keyManager.GetPrivateKey(context.Background()); err != nil {
		p.logger.Sugar().Errorw("Failed to access VRF private key",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return fmt.Errorf("VRF key not accessible: %w", err)
	}

	p.logger.Sugar().Infow("VRF task validation successful",
		zap.String("taskId", string(t.TaskId)),
		zap.String("taskHash", hex.EncodeToString(taskData.TaskHash[:])),
		zap.String("seed", taskData.Seed.String()),
	)

	return nil
}

// HandleTask processes the VRF task and computes the VRF proof
func (p *VRFPerformer) HandleTask(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	p.logger.Sugar().Infow("Processing VRF task",
		zap.String("taskId", string(t.TaskId)),
	)

	// Parse ABI-encoded task data
	taskData, err := p.decodeTaskData(t.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode task data: %w", err)
	}

	// Get actual task hash from the TaskId (this is the real task hash from TaskMailbox)
	taskHash := common.BytesToHash(t.TaskId)

	// Get VRF private key
	privateKey, err := p.keyManager.GetPrivateKey(context.Background())
	if err != nil {
		p.logger.Sugar().Errorw("Failed to get VRF private key",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return nil, fmt.Errorf("failed to get VRF private key: %w", err)
	}

	// Convert seed to bytes for VRF computation
	seedBytes := taskData.Seed.Bytes()

	// Compute VRF proof
	p.logger.Sugar().Infow("Computing VRF proof",
		zap.String("taskId", string(t.TaskId)),
		zap.String("taskHash", taskHash.Hex()),
		zap.String("seed", hex.EncodeToString(seedBytes)),
	)

	beta, pi, err := p.suite.Prove(privateKey, seedBytes)
	if err != nil {
		p.logger.Sugar().Errorw("Failed to compute VRF proof",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return nil, fmt.Errorf("failed to compute VRF proof: %w", err)
	}

	// Convert VRF output (beta) to big.Int
	vrfOutput := new(big.Int).SetBytes(beta)

	// Create result
	result := VRFTaskResult{
		TaskHash:  taskHash,
		VRFProof:  pi,
		VRFOutput: vrfOutput,
	}

	// Encode result using ABI
	resultBytes, err := p.encodeTaskResult(result)
	if err != nil {
		return nil, fmt.Errorf("failed to encode result: %w", err)
	}

	p.logger.Sugar().Infow("VRF task completed successfully",
		zap.String("taskId", string(t.TaskId)),
		zap.String("taskHash", taskHash.Hex()),
		zap.String("vrfOutput", hex.EncodeToString(beta)),
		zap.Int("proofLength", len(pi)),
	)

	return &performerV1.TaskResponse{
		TaskId: t.TaskId,
		Result: resultBytes,
	}, nil
}

// VerifyVRFProof verifies a VRF proof (utility function for testing)
func (p *VRFPerformer) VerifyVRFProof(publicKey *ecdsa.PublicKey, seed []byte, proof []byte) ([]byte, error) {
	return p.suite.Verify(publicKey, seed, proof)
}

// decodeTaskData decodes ABI-encoded VRFTaskData from bytes
func (p *VRFPerformer) decodeTaskData(payload []byte) (*VRFTaskData, error) {
	// Define the ABI type for VRFTaskData struct
	vrfTaskDataType, err := abi.NewType("tuple", "VRFTaskData", []abi.ArgumentMarshaling{
		{Name: "taskHash", Type: "bytes32"},
		{Name: "seed", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	// Create arguments for decoding
	arguments := abi.Arguments{{Type: vrfTaskDataType}}

	// Decode the payload
	decoded, err := arguments.Unpack(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ABI data: %w", err)
	}

	if len(decoded) != 1 {
		return nil, fmt.Errorf("expected 1 decoded value, got %d", len(decoded))
	}

	// The ABI library returns a struct with the correct field names
	// Use reflection to extract the fields
	structValue := reflect.ValueOf(decoded[0])
	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", structValue.Kind())
	}

	// Extract TaskHash field
	taskHashField := structValue.FieldByName("TaskHash")
	if !taskHashField.IsValid() {
		return nil, fmt.Errorf("TaskHash field not found")
	}
	taskHashBytes, ok := taskHashField.Interface().([32]byte)
	if !ok {
		return nil, fmt.Errorf("failed to cast TaskHash to [32]byte, got %T", taskHashField.Interface())
	}

	// Extract Seed field
	seedField := structValue.FieldByName("Seed")
	if !seedField.IsValid() {
		return nil, fmt.Errorf("Seed field not found")
	}
	seed, ok := seedField.Interface().(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to cast Seed to *big.Int, got %T", seedField.Interface())
	}

	return &VRFTaskData{
		TaskHash: taskHashBytes,
		Seed:     seed,
	}, nil
}

// encodeTaskResult encodes VRFTaskResult to ABI bytes
func (p *VRFPerformer) encodeTaskResult(result VRFTaskResult) ([]byte, error) {
	// Define the ABI type for VRFTaskResult struct
	vrfTaskResultType, err := abi.NewType("tuple", "VRFTaskResult", []abi.ArgumentMarshaling{
		{Name: "taskHash", Type: "bytes32"},
		{Name: "vrfProof", Type: "bytes"},
		{Name: "vrfOutput", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	// Create arguments for encoding
	arguments := abi.Arguments{{Type: vrfTaskResultType}}

	// Create the struct data for encoding
	structData := struct {
		TaskHash  [32]byte `abi:"taskHash"`
		VRFProof  []byte   `abi:"vrfProof"`
		VRFOutput *big.Int `abi:"vrfOutput"`
	}{
		TaskHash:  result.TaskHash,
		VRFProof:  result.VRFProof,
		VRFOutput: result.VRFOutput,
	}

	// Encode the result
	encoded, err := arguments.Pack(structData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ABI data: %w", err)
	}

	return encoded, nil
}

// EncodeTaskData encodes VRFTaskData to ABI bytes (exported for testing)
func (p *VRFPerformer) EncodeTaskData(data VRFTaskData) ([]byte, error) {
	// Define the ABI type for VRFTaskData struct
	vrfTaskDataType, err := abi.NewType("tuple", "VRFTaskData", []abi.ArgumentMarshaling{
		{Name: "taskHash", Type: "bytes32"},
		{Name: "seed", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	// Create arguments for encoding
	arguments := abi.Arguments{{Type: vrfTaskDataType}}

	// Create the struct data for encoding
	structData := struct {
		TaskHash [32]byte `abi:"taskHash"`
		Seed     *big.Int `abi:"seed"`
	}{
		TaskHash: data.TaskHash,
		Seed:     data.Seed,
	}

	// Encode the data
	encoded, err := arguments.Pack(structData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ABI data: %w", err)
	}

	return encoded, nil
}

// DecodeTaskResult decodes ABI-encoded VRFTaskResult from bytes (exported for testing)
func (p *VRFPerformer) DecodeTaskResult(payload []byte) (*VRFTaskResult, error) {
	// Define the ABI type for VRFTaskResult struct
	vrfTaskResultType, err := abi.NewType("tuple", "VRFTaskResult", []abi.ArgumentMarshaling{
		{Name: "taskHash", Type: "bytes32"},
		{Name: "vrfProof", Type: "bytes"},
		{Name: "vrfOutput", Type: "uint256"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	// Create arguments for decoding
	arguments := abi.Arguments{{Type: vrfTaskResultType}}

	// Decode the payload
	decoded, err := arguments.Unpack(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ABI data: %w", err)
	}

	if len(decoded) != 1 {
		return nil, fmt.Errorf("expected 1 decoded value, got %d", len(decoded))
	}

	// The ABI library returns a struct with the correct field names
	// Use reflection to extract the fields
	structValue := reflect.ValueOf(decoded[0])
	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", structValue.Kind())
	}

	// Extract TaskHash field
	taskHashField := structValue.FieldByName("TaskHash")
	if !taskHashField.IsValid() {
		return nil, fmt.Errorf("TaskHash field not found")
	}
	taskHashBytes, ok := taskHashField.Interface().([32]byte)
	if !ok {
		return nil, fmt.Errorf("failed to cast TaskHash to [32]byte, got %T", taskHashField.Interface())
	}

	// Extract VRFProof field
	vrfProofField := structValue.FieldByName("VrfProof")
	if !vrfProofField.IsValid() {
		return nil, fmt.Errorf("VrfProof field not found")
	}
	vrfProof, ok := vrfProofField.Interface().([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to cast VrfProof to []byte, got %T", vrfProofField.Interface())
	}

	// Extract VRFOutput field
	vrfOutputField := structValue.FieldByName("VrfOutput")
	if !vrfOutputField.IsValid() {
		return nil, fmt.Errorf("VrfOutput field not found")
	}
	vrfOutput, ok := vrfOutputField.Interface().(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to cast VrfOutput to *big.Int, got %T", vrfOutputField.Interface())
	}

	return &VRFTaskResult{
		TaskHash:  taskHashBytes,
		VRFProof:  vrfProof,
		VRFOutput: vrfOutput,
	}, nil
}