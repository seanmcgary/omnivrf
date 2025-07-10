package main

import (
	"context"
	"os"
	"time"

	"github.com/Layr-Labs/hourglass-monorepo/ponos/pkg/performer/server"
	"go.uber.org/zap"

	"github.com/Layr-Labs/hourglass-avs-template/pkg/keymanager"
	"github.com/Layr-Labs/hourglass-avs-template/pkg/performer"
)

// This offchain binary is run by Operators running the Hourglass Executor. It contains
// the business logic of the AVS and performs worked based on the tasked sent to it.
// The Hourglass Aggregator ingests tasks from the TaskMailbox and distributes work
// to Executors configured to run the AVS Performer. Performers execute the work and
// return the result to the Executor where the result is signed and return to the
// Aggregator to place in the outbox once the signing threshold is met.

func main() {
	ctx := context.Background()
	l, _ := zap.NewProduction()

	l.Info("Starting OmniVRF Performer")

	// Initialize key manager based on configuration
	var km performer.KeyManager
	var err error

	keyManagerType := os.Getenv("KEY_MANAGER_TYPE")
	switch keyManagerType {
	case "aws":
		// TODO: Implement AWS key manager in Phase 2
		l.Fatal("AWS key manager not yet implemented")
	default:
		l.Info("Using environment variable key manager")
		km, err = keymanager.NewEnvKeyManager()
		if err != nil {
			l.Fatal("Failed to initialize environment key manager", zap.Error(err))
		}
	}

	// Create VRF performer
	vrfPerformer := performer.NewVRFPerformer(km, l)

	// Initialize Hourglass server with VRF performer
	pp, err := server.NewPonosPerformerWithRpcServer(&server.PonosPerformerConfig{
		Port:    8080,
		Timeout: 30 * time.Second, // Increased timeout for VRF computation
	}, vrfPerformer, l)
	if err != nil {
		l.Fatal("Failed to create performer server", zap.Error(err))
	}

	l.Info("Starting VRF performer server on port 8080")

	if err := pp.Start(ctx); err != nil {
		l.Fatal("Failed to start performer server", zap.Error(err))
	}
}
