# OmniVRF Client CLI

A command-line interface for interacting with the OmniVRF smart contract.

## Installation

```bash
cd cmd/client-cli
go build -o omnivrf-client .
```

## Configuration

The client requires three main configuration parameters:

```bash
export OMNIVRF_RPC_URL="http://localhost:8545"
export OMNIVRF_PRIVATE_KEY="0x..."
export OMNIVRF_CONTRACT_ADDRESS="0x..."
```

Or provide them as flags:

```bash
./omnivrf-client --rpc-url http://localhost:8545 \
  --private-key 0x... \
  --contract-address 0x... \
  <command>
```

## Commands

### Request Randomness

Request verifiable randomness from the OmniVRF contract:

```bash
# Simple request (no callback)
./omnivrf-client request

# Request with callback
./omnivrf-client request \
  --callback-address 0x1234... \
  --callback-gas 200000

# Request and wait for result
./omnivrf-client request --wait

# Wait with custom timeout
./omnivrf-client request \
  --wait \
  --timeout 10m

# With custom gas settings
./omnivrf-client request \
  --gas-price 20 \
  --gas-limit 300000
```

Options:
- `--callback-address`: Contract to receive callback (default: 0x0)
- `--callback-gas`: Gas limit for callback (default: 200000)
- `--gas-price`: Gas price in gwei (default: 10)
- `--gas-limit`: Transaction gas limit (default: 500000)
- `--wait`: Wait for the randomness result via WebSocket/polling (default: false)
- `--timeout`: Timeout duration when waiting (default: 5m)

### Get Randomness Result

Check if a randomness request has been fulfilled:

```bash
./omnivrf-client get --task-hash 0xabc123...
```

Output shows:
- Whether the request is fulfilled
- The random value (if fulfilled)

### Withdraw Callback Gas

Withdraw unused callback gas deposits:

```bash
# Withdraw 1 ether worth of gas deposits
./omnivrf-client withdraw --amount 1000000000000000000

# Withdraw with custom gas settings
./omnivrf-client withdraw \
  --amount 500000000000000000 \
  --gas-price 15
```

### Contract Information

Get contract configuration and your current deposits:

```bash
./omnivrf-client info
```

Shows:
- Contract address and task mailbox
- Callback gas price and limits
- Your current callback gas deposits

## Examples

### Complete Workflow

```bash
# 1. Check contract info
./omnivrf-client info

# 2. Request randomness with callback
./omnivrf-client request \
  --callback-address 0x742d35Cc6634C0532925a3b844Bc9e7595f6E123 \
  --callback-gas 200000

# Output:
# Transaction submitted!
#   Hash: 0x123...
# Randomness requested successfully!
#   Task Hash: 0xabc...

# 3. Check result
./omnivrf-client get --task-hash 0xabc...

# 4. Withdraw unused gas deposits
./omnivrf-client withdraw --amount 100000000000000000
```

### Request and Wait for Result

```bash
# Request randomness and wait for the result in one command
./omnivrf-client request --wait

# Output:
# Requesting randomness...
# Transaction submitted!
#   Hash: 0x123...
# Transaction confirmed!
#   Block: 12345
# Randomness requested successfully!
#   Task Hash: 0xabc...
# 
# Waiting for randomness result (timeout: 5m0s)...
# ⠹ Polling for result...
# 
# ✓ Randomness fulfilled!
#   Task Hash: 0xabc...
#   Randomness: 123456789...
```

### Environment Setup for Different Networks

```bash
# Ethereum Mainnet
export OMNIVRF_RPC_URL="https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY"

# Base
export OMNIVRF_RPC_URL="https://base-mainnet.g.alchemy.com/v2/YOUR_KEY"

# Local DevNet
export OMNIVRF_RPC_URL="http://localhost:8545"
```

## Gas Costs

When requesting randomness with a callback:
- You must send ETH to cover: `callback_gas_limit * 10 gwei`
- Example: 200,000 gas limit = 0.002 ETH deposit
- Unused gas is refunded to your deposit balance
- Withdraw deposits using the `withdraw` command

## Security

- Never share or commit your private key
- Use environment variables or secure key management
- Consider using a dedicated wallet for interactions
- Monitor your callback gas deposits

## Features

### WebSocket Support
When using `--wait`, the client will attempt to connect via WebSocket for real-time event monitoring. If WebSocket is not available (e.g., using HTTP-only RPC), it will automatically fall back to polling every 2 seconds.

### Automatic URL Conversion
The client automatically converts HTTP URLs to WebSocket URLs:
- `http://` → `ws://`
- `https://` → `wss://`

## Troubleshooting

### "Insufficient callback gas"
You're not sending enough ETH for the callback. The contract requires:
```
Required ETH = callback_gas_limit * 10 gwei
```

### "Transaction failed"
Check:
- Contract address is correct
- You have enough ETH for gas + callback deposit
- RPC URL is accessible
- Private key has funds

### "Invalid task hash length"
Task hashes must be exactly 32 bytes (64 hex characters). Remove any "0x" prefix.

### "WebSocket not available"
If you see this message when using `--wait`, the client will automatically fall back to polling. This is normal for HTTP-only RPC endpoints.