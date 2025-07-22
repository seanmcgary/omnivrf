# OmniVRF Performer CLI

The OmniVRF performer is the off-chain component that executes VRF computations for the AVS.

## Installation

```bash
make build
# Binary will be at ./bin/performer
```

## Usage

### Basic Usage

```bash
# Run with defaults
./bin/performer

# Run with custom configuration
./bin/performer \
  --port 9090 \
  --timeout 60s \
  --log-level debug
```

### Configuration Options

The performer can be configured using command-line flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--port` | `OMNIVRF_PORT` | `8080` | RPC server port |
| `--timeout` | `OMNIVRF_TIMEOUT` | `30s` | VRF computation timeout |
| `--key-manager-type` | `KEY_MANAGER_TYPE` | `env` | Key manager type (`env` or `aws`) |
| `--vrf-private-key` | `VRF_PRIVATE_KEY` | - | VRF private key (required for env key manager) |
| `--aws-secret-name` | `AWS_SECRET_NAME` | - | AWS secret name (required for aws key manager) |
| `--aws-region` | `AWS_REGION` | `us-east-1` | AWS region for secrets manager |
| `--log-level` | `OMNIVRF_LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |

### Key Management

#### Environment Variable Key Manager (Default)

```bash
export VRF_PRIVATE_KEY=0x1234567890abcdef...
./bin/performer --key-manager-type env
```

#### AWS Secrets Manager (Future)

```bash
./bin/performer \
  --key-manager-type aws \
  --aws-secret-name omnivrf-key \
  --aws-region us-west-2
```

### Docker Usage

```bash
# Build container
make build/container

# Run with environment variables
docker run -p 8080:8080 \
  -e VRF_PRIVATE_KEY=0x... \
  -e OMNIVRF_LOG_LEVEL=debug \
  my-avs:latest
```

### Development

For development with the devkit:

```bash
# Start devnet without AVS
devkit avs devnet start --skip-avs-run

# Run performer separately with custom config
./bin/performer --port 8080 --log-level debug

# Or use devkit to run with defaults
devkit avs run
```

### Signals

The performer handles graceful shutdown on:
- `SIGINT` (Ctrl+C)
- `SIGTERM`

## Architecture

The performer implements the Hourglass Performer interface:
- `ValidateTask()` - Validates incoming VRF task requests
- `HandleTask()` - Computes VRF proofs using the go-ecvrf library

Tasks are received via RPC from the Hourglass Executor and results are returned for aggregation.