package performer

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestABIEncodingDecoding(t *testing.T) {
	// Create a VRF performer for testing
	performer := &VRFPerformer{
		logger: zap.NewNop(),
	}

	t.Run("encode and decode VRFTaskData", func(t *testing.T) {
		// Create test data
		originalTaskHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		originalSeed := new(big.Int)
		originalSeed.SetString("12345678901234567890", 10)

		originalData := VRFTaskData{
			TaskHash: originalTaskHash,
			Seed:     originalSeed,
		}

		// Encode to bytes (simulate what Solidity would do)
		encoded, err := performer.EncodeTaskData(originalData)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)

		// Decode back to struct
		decoded, err := performer.decodeTaskData(encoded)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		// Verify data integrity
		assert.Equal(t, originalTaskHash, common.Hash(decoded.TaskHash))
		assert.Equal(t, originalSeed, decoded.Seed)
	})

	t.Run("encode and decode VRFTaskResult", func(t *testing.T) {
		// Create test result data
		taskHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
		vrfProof := []byte("test_vrf_proof_data_123456789")
		vrfOutput := big.NewInt(987654321098765432)

		originalResult := VRFTaskResult{
			TaskHash:  taskHash,
			VRFProof:  vrfProof,
			VRFOutput: vrfOutput,
		}

		// Encode to bytes
		encoded, err := performer.encodeTaskResult(originalResult)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)

		// Decode back to struct
		decoded, err := performer.DecodeTaskResult(encoded)
		require.NoError(t, err)
		require.NotNil(t, decoded)

		// Verify data integrity
		assert.Equal(t, taskHash, common.Hash(decoded.TaskHash))
		assert.Equal(t, vrfProof, decoded.VRFProof)
		assert.Equal(t, vrfOutput, decoded.VRFOutput)
	})

	t.Run("handle zero values", func(t *testing.T) {
		// Test with zero/empty values
		originalData := VRFTaskData{
			TaskHash: common.Hash{}, // Zero hash
			Seed:     big.NewInt(0),  // Zero seed
		}

		// Encode and decode
		encoded, err := performer.EncodeTaskData(originalData)
		require.NoError(t, err)

		decoded, err := performer.decodeTaskData(encoded)
		require.NoError(t, err)

		// Verify zero values are preserved
		assert.Equal(t, common.Hash{}, common.Hash(decoded.TaskHash))
		assert.True(t, decoded.Seed.Cmp(big.NewInt(0)) == 0)
	})

	t.Run("handle large values", func(t *testing.T) {
		// Test with maximum values
		maxSeed := new(big.Int)
		maxSeed.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10) // 2^256 - 1

		originalData := VRFTaskData{
			TaskHash: common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			Seed:     maxSeed,
		}

		// Encode and decode
		encoded, err := performer.EncodeTaskData(originalData)
		require.NoError(t, err)

		decoded, err := performer.decodeTaskData(encoded)
		require.NoError(t, err)

		// Verify large values are preserved
		assert.Equal(t, common.Hash(originalData.TaskHash), common.Hash(decoded.TaskHash))
		assert.True(t, maxSeed.Cmp(decoded.Seed) == 0)
	})

	t.Run("invalid data should fail decoding", func(t *testing.T) {
		// Test with invalid data
		invalidData := []byte("this is not valid ABI data")

		_, err := performer.decodeTaskData(invalidData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode ABI data")
	})

	t.Run("empty data should fail decoding", func(t *testing.T) {
		// Test with empty data
		emptyData := []byte{}

		_, err := performer.decodeTaskData(emptyData)
		assert.Error(t, err)
	})
}

func TestSolidityCompatibility(t *testing.T) {
	// This test verifies that our Go ABI encoding produces the same output
	// as Solidity's abi.encode() for the VRFTaskData struct
	performer := &VRFPerformer{
		logger: zap.NewNop(),
	}

	t.Run("matches Solidity abi.encode output format", func(t *testing.T) {
		// Create test data with known values
		taskHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		seed := big.NewInt(42)

		taskData := VRFTaskData{
			TaskHash: taskHash,
			Seed:     seed,
		}

		// Encode
		encoded, err := performer.EncodeTaskData(taskData)
		require.NoError(t, err)

		// Verify structure: should be 64 bytes (32 for taskHash + 32 for seed)
		assert.Equal(t, 64, len(encoded))

		// Verify taskHash is in the first 32 bytes
		encodedTaskHash := encoded[0:32]
		assert.Equal(t, taskHash.Bytes(), encodedTaskHash)

		// Verify seed is in the last 32 bytes (big-endian)
		encodedSeed := new(big.Int).SetBytes(encoded[32:64])
		assert.Equal(t, seed, encodedSeed)
	})
}