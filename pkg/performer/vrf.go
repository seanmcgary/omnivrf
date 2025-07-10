package performer

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

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

// TaskData represents the structure of incoming VRF task data
type TaskData struct {
	RequestID *big.Int `json:"requestId"`
	Seed      string   `json:"seed"`
}

// TaskResult represents the VRF computation result
type TaskResult struct {
	RequestID *big.Int `json:"requestId"`
	VRFProof  string   `json:"vrfProof"`
	VRFOutput string   `json:"vrfOutput"`
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

	// Parse and validate task data
	var taskData TaskData
	if err := json.Unmarshal(t.Payload, &taskData); err != nil {
		p.logger.Sugar().Errorw("Failed to unmarshal task data",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return fmt.Errorf("invalid task data format: %w", err)
	}

	// Validate required fields
	if taskData.RequestID == nil {
		return fmt.Errorf("missing requestId in task data")
	}

	if taskData.Seed == "" {
		return fmt.Errorf("missing seed in task data")
	}

	// Validate seed is valid hex
	if _, err := hex.DecodeString(taskData.Seed); err != nil {
		return fmt.Errorf("invalid seed format, must be hex: %w", err)
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
		zap.String("requestId", taskData.RequestID.String()),
	)

	return nil
}

// HandleTask processes the VRF task and computes the VRF proof
func (p *VRFPerformer) HandleTask(t *performerV1.TaskRequest) (*performerV1.TaskResponse, error) {
	p.logger.Sugar().Infow("Processing VRF task",
		zap.String("taskId", string(t.TaskId)),
	)

	// Parse task data
	var taskData TaskData
	if err := json.Unmarshal(t.Payload, &taskData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task data: %w", err)
	}

	// Get VRF private key
	privateKey, err := p.keyManager.GetPrivateKey(context.Background())
	if err != nil {
		p.logger.Sugar().Errorw("Failed to get VRF private key",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return nil, fmt.Errorf("failed to get VRF private key: %w", err)
	}

	// Convert seed from hex to bytes
	seedBytes, err := hex.DecodeString(taskData.Seed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode seed: %w", err)
	}

	// Compute VRF proof
	p.logger.Sugar().Infow("Computing VRF proof",
		zap.String("taskId", string(t.TaskId)),
		zap.String("requestId", taskData.RequestID.String()),
		zap.String("seed", taskData.Seed),
	)

	beta, pi, err := p.suite.Prove(privateKey, seedBytes)
	if err != nil {
		p.logger.Sugar().Errorw("Failed to compute VRF proof",
			zap.Error(err),
			zap.String("taskId", string(t.TaskId)),
		)
		return nil, fmt.Errorf("failed to compute VRF proof: %w", err)
	}

	// Create result
	result := TaskResult{
		RequestID: taskData.RequestID,
		VRFProof:  hex.EncodeToString(pi),
		VRFOutput: hex.EncodeToString(beta),
	}

	// Marshal result to JSON
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	p.logger.Sugar().Infow("VRF task completed successfully",
		zap.String("taskId", string(t.TaskId)),
		zap.String("requestId", taskData.RequestID.String()),
		zap.String("vrfOutput", result.VRFOutput),
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