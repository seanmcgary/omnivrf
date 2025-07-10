# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an EigenLayer AVS (Autonomous Verifiable Service) project named "omnivrf" built using the Hourglass framework. It's a task-based AVS that processes work distributed by the Hourglass infrastructure, with both on-chain (smart contracts) and off-chain (Go performer) components.

## Common Development Commands

### Build Commands
- `make build` - Builds the Go performer binary to `./bin/performer`
- `make deps` - Updates Go dependencies
- `make build/container` - Builds Docker container using `.hourglass/scripts/buildContainer.sh`
- `devkit avs build` - Builds both Go code and smart contracts

### Testing Commands
- `make test` - Runs both Go and Forge tests
- `make test-go` - Runs only Go tests with verbose output
- `make test-forge` - Runs only Solidity/Forge tests
- `devkit avs test` - Alternative way to run all tests

### Development Workflow Commands
- `devkit avs devnet start` - Starts local devnet with deployed contracts and infrastructure
- `devkit avs devnet start --skip-avs-run` - Starts devnet without running AVS performer locally
- `devkit avs run` - Runs AVS performer separately (if started with --skip-avs-run)
- `devkit avs devnet stop` - Stops the devnet
- `devkit avs call --signature="(uint256,string)" args='(5,"hello")'` - Simulates task execution

### Running the Performer Directly
The performer now supports a proper CLI interface with flags and environment variables:

```bash
# Run with default settings
./bin/performer

# Run with custom port and timeout
./bin/performer --port 9090 --timeout 60s

# Run with debug logging
./bin/performer --log-level debug

# Run with environment variables
export OMNIVRF_PORT=9090
export OMNIVRF_TIMEOUT=60s
export VRF_PRIVATE_KEY=0x...
./bin/performer

# View all available options
./bin/performer --help
```

Available flags:
- `--port` (env: `OMNIVRF_PORT`): RPC server port (default: 8080)
- `--timeout` (env: `OMNIVRF_TIMEOUT`): VRF computation timeout (default: 30s)
- `--key-manager-type` (env: `KEY_MANAGER_TYPE`): Key manager type - "env" or "aws" (default: env)
- `--vrf-private-key` (env: `VRF_PRIVATE_KEY`): VRF private key for env key manager
- `--aws-secret-name` (env: `AWS_SECRET_NAME`): AWS secret name for aws key manager
- `--aws-region` (env: `AWS_REGION`): AWS region (default: us-east-1)
- `--log-level` (env: `OMNIVRF_LOG_LEVEL`): Log level - debug, info, warn, error (default: info)

### Linting and Type Checking
Note: No explicit linting commands found in the Makefile. The project likely uses:
- `go fmt` for Go formatting
- Forge formatting for Solidity

## High-Level Architecture

### Core Components

1. **AVS Performer (Off-chain)**
   - Entry point: `cmd/main.go`
   - Implements `ValidateTask()` and `HandleTask()` methods
   - Runs as an RPC server on port 8080
   - Handles task validation and execution logic

2. **Smart Contracts Architecture**
   - **L1 Contracts** (`contracts/src/l1-contracts/`):
     - `TaskAVSRegistrar.sol` - Manages AVS registration and operator management
     - `HelloWorldL1.sol` - Example contract (can be replaced with custom logic)
   
   - **L2 Contracts** (`contracts/src/l2-contracts/`):
     - `AVSTaskHook.sol` - Implements task lifecycle hooks (validation, pre/post processing)
     - `HelloWorldL2.sol` - Example contract (can be replaced with custom logic)

3. **Task Flow**
   - Tasks are created through the TaskMailbox on L2
   - AVSTaskHook validates tasks before creation
   - Performer processes tasks off-chain
   - Results are submitted back through the Hourglass infrastructure
   - AVSTaskHook validates results before acceptance

### Key Integration Points

1. **EigenLayer Integration**
   - Uses EigenLayer's operator management system
   - Integrates with AllocationManager and KeyRegistrar
   - Supports both BLS (BN254) and ECDSA operator keys

2. **Hourglass Framework**
   - TaskMailbox handles task distribution
   - Aggregator collects and validates results
   - Executor runs the AVS Performer

3. **Configuration**
   - `config/config.yaml` - Project metadata
   - `config/contexts/devnet.yaml` - Environment-specific settings
   - `.env` - RPC endpoints for L1/L2 forks (Ethereum Sepolia, Base Sepolia)

### Development Patterns

1. **Extension Points**
   - `cmd/main.go`: Implement task validation and execution logic in `ValidateTask()` and `HandleTask()`
   - `contracts/src/l1-contracts/`: Add custom L1 contracts
   - `contracts/src/l2-contracts/`: Add custom L2 contracts
   - `contracts/script/DeployMyContracts.s.sol`: Wire up custom contract deployments

2. **Testing Structure**
   - Go tests: `cmd/main_test.go` - Unit tests for performer logic
   - Solidity tests: `contracts/test/` - Forge tests for smart contracts

3. **Deployment Flow**
   - L1 contracts deployed via `DeployMyL1Contracts.s.sol`
   - L2 contracts deployed via `DeployMyL2Contracts.s.sol`
   - Deployment outputs saved to JSON for environment tracking

## Important Notes

- The project uses Go 1.23.6 and Solidity ^0.8.27
- Requires Docker, Foundry, and the DevKit CLI for full development workflow
- Pre-generated operator keystores available in `keystores/` directory
- The codebase is in initial commit stage with template version v0.0.14
- Fork URLs for L1/L2 must be archive nodes (not full nodes) for proper devnet functionality