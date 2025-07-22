package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Layr-Labs/hourglass-monorepo/ponos/pkg/performer/server"
	"github.com/urfave/cli/v2"
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

const (
	// Flag names
	FlagPort           = "port"
	FlagTimeout        = "timeout"
	FlagKeyManagerType = "key-manager-type"
	FlagVRFPrivateKey  = "vrf-private-key"
	FlagAWSSecretName  = "aws-secret-name"
	FlagAWSRegion      = "aws-region"
	FlagLogLevel       = "log-level"

	// Environment variable names
	EnvPort           = "OMNIVRF_PORT"
	EnvTimeout        = "OMNIVRF_TIMEOUT"
	EnvKeyManagerType = "KEY_MANAGER_TYPE"
	EnvVRFPrivateKey  = "VRF_PRIVATE_KEY"
	EnvAWSSecretName  = "AWS_SECRET_NAME"
	EnvAWSRegion      = "AWS_REGION"
	EnvLogLevel       = "OMNIVRF_LOG_LEVEL"
)

func main() {
	app := &cli.App{
		Name:    "omnivrf-performer",
		Usage:   "OmniVRF AVS Performer for EigenLayer",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    FlagPort,
				Value:   8080,
				Usage:   "Port for the RPC server",
				EnvVars: []string{EnvPort},
			},
			&cli.DurationFlag{
				Name:    FlagTimeout,
				Value:   30 * time.Second,
				Usage:   "Timeout for VRF computation",
				EnvVars: []string{EnvTimeout},
			},
			&cli.StringFlag{
				Name:    FlagKeyManagerType,
				Value:   "env",
				Usage:   "Key manager type (env, aws)",
				EnvVars: []string{EnvKeyManagerType},
			},
			&cli.StringFlag{
				Name:    FlagVRFPrivateKey,
				Usage:   "VRF private key (required for env key manager)",
				EnvVars: []string{EnvVRFPrivateKey},
			},
			&cli.StringFlag{
				Name:    FlagAWSSecretName,
				Usage:   "AWS secret name for VRF key (required for aws key manager)",
				EnvVars: []string{EnvAWSSecretName},
			},
			&cli.StringFlag{
				Name:    FlagAWSRegion,
				Value:   "us-east-1",
				Usage:   "AWS region for secrets manager",
				EnvVars: []string{EnvAWSRegion},
			},
			&cli.StringFlag{
				Name:    FlagLogLevel,
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
				EnvVars: []string{EnvLogLevel},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	// Setup logger
	logger, err := setupLogger(c.String(FlagLogLevel))
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	defer logger.Sync()

	logger.Info("Starting OmniVRF Performer",
		zap.String("version", c.App.Version),
		zap.Int("port", c.Int(FlagPort)),
		zap.Duration("timeout", c.Duration(FlagTimeout)),
		zap.String("key_manager_type", c.String(FlagKeyManagerType)),
	)

	// Initialize key manager
	km, err := initializeKeyManager(c, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize key manager: %w", err)
	}

	// Create VRF performer
	vrfPerformer := performer.NewVRFPerformer(km, logger)

	// Initialize Hourglass server with VRF performer
	pp, err := server.NewPonosPerformerWithRpcServer(&server.PonosPerformerConfig{
		Port:    c.Int(FlagPort),
		Timeout: c.Duration(FlagTimeout),
	}, vrfPerformer, logger)
	if err != nil {
		return fmt.Errorf("failed to create performer server: %w", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Starting VRF performer server", zap.Int("port", c.Int(FlagPort)))
		if err := pp.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
		// Give server time to gracefully shutdown
		time.Sleep(2 * time.Second)
		return nil
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

func setupLogger(level string) (*zap.Logger, error) {
	var config zap.Config
	
	switch level {
	case "debug":
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config = zap.NewProductionConfig()
	}

	return config.Build()
}

func initializeKeyManager(c *cli.Context, logger *zap.Logger) (performer.KeyManager, error) {
	keyManagerType := c.String(FlagKeyManagerType)
	
	switch keyManagerType {
	case "aws":
		// TODO: Implement AWS key manager in Phase 2
		secretName := c.String(FlagAWSSecretName)
		if secretName == "" {
			return nil, fmt.Errorf("aws-secret-name is required for AWS key manager")
		}
		return nil, fmt.Errorf("AWS key manager not yet implemented")
		
	case "env":
		logger.Info("Using environment variable key manager")
		// Check if VRF private key is provided
		if c.String(FlagVRFPrivateKey) == "" {
			return nil, fmt.Errorf("vrf-private-key is required for env key manager")
		}
		
		km, err := keymanager.NewEnvKeyManager()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize environment key manager: %w", err)
		}
		return km, nil
		
	default:
		return nil, fmt.Errorf("unknown key manager type: %s", keyManagerType)
	}
}
