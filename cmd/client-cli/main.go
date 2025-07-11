package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/hourglass-avs-template/pkg/bindings/OmniVRF"
)

const (
	// Flag names
	FlagRPCURL          = "rpc-url"
	FlagPrivateKey      = "private-key"
	FlagContractAddress = "contract-address"
	FlagCallbackAddress = "callback-address"
	FlagCallbackGas     = "callback-gas"
	FlagTaskHash        = "task-hash"
	FlagAmount          = "amount"
	FlagGasPrice        = "gas-price"
	FlagGasLimit        = "gas-limit"
	FlagWaitForResult   = "wait"
	FlagTimeout         = "timeout"

	// Environment variable names
	EnvRPCURL          = "OMNIVRF_RPC_URL"
	EnvPrivateKey      = "OMNIVRF_PRIVATE_KEY"
	EnvContractAddress = "OMNIVRF_CONTRACT_ADDRESS"
)

func main() {
	app := &cli.App{
		Name:    "omnivrf-client",
		Usage:   "OmniVRF Contract CLI Client",
		Version: "0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     FlagRPCURL,
				Usage:    "Ethereum RPC URL",
				EnvVars:  []string{EnvRPCURL},
				Required: true,
			},
			&cli.StringFlag{
				Name:     FlagPrivateKey,
				Usage:    "Ethereum private key (with or without 0x prefix)",
				EnvVars:  []string{EnvPrivateKey},
				Required: true,
			},
			&cli.StringFlag{
				Name:     FlagContractAddress,
				Usage:    "OmniVRF contract address",
				EnvVars:  []string{EnvContractAddress},
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "request",
				Aliases: []string{"req"},
				Usage:   "Request randomness from OmniVRF",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  FlagCallbackAddress,
						Usage: "Callback contract address (use 0x0 for no callback)",
						Value: "0x0000000000000000000000000000000000000000",
					},
					&cli.Uint64Flag{
						Name:  FlagCallbackGas,
						Usage: "Callback gas limit",
						Value: 200000,
					},
					&cli.StringFlag{
						Name:  FlagGasPrice,
						Usage: "Gas price in gwei",
						Value: "10",
					},
					&cli.Uint64Flag{
						Name:  FlagGasLimit,
						Usage: "Transaction gas limit",
						Value: 500000,
					},
					&cli.BoolFlag{
						Name:  FlagWaitForResult,
						Usage: "Wait for the randomness result via websocket",
						Value: false,
					},
					&cli.DurationFlag{
						Name:  FlagTimeout,
						Usage: "Timeout for waiting for result",
						Value: 5 * time.Minute,
					},
				},
				Action: requestRandomness,
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get randomness result for a task hash",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     FlagTaskHash,
						Usage:    "Task hash to query",
						Required: true,
					},
				},
				Action: getRandomness,
			},
			{
				Name:    "withdraw",
				Aliases: []string{"w"},
				Usage:   "Withdraw callback gas deposits",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     FlagAmount,
						Usage:    "Amount to withdraw in wei",
						Required: true,
					},
					&cli.StringFlag{
						Name:  FlagGasPrice,
						Usage: "Gas price in gwei",
						Value: "10",
					},
					&cli.Uint64Flag{
						Name:  FlagGasLimit,
						Usage: "Transaction gas limit",
						Value: 100000,
					},
				},
				Action: withdrawCallbackGas,
			},
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Get contract information",
				Action:  getContractInfo,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getClient(c *cli.Context) (*ethclient.Client, *ecdsa.PrivateKey, error) {
	// Connect to Ethereum node
	client, err := ethclient.Dial(c.String(FlagRPCURL))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Parse private key
	privateKeyHex := c.String(FlagPrivateKey)
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return client, privateKey, nil
}

func getContract(c *cli.Context, client *ethclient.Client) (*OmniVRF.OmniVRF, error) {
	contractAddress := common.HexToAddress(c.String(FlagContractAddress))
	
	contract, err := OmniVRF.NewOmniVRF(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	return contract, nil
}

func requestRandomness(c *cli.Context) error {
	client, privateKey, err := getClient(c)
	if err != nil {
		return err
	}
	defer client.Close()

	contract, err := getContract(c, client)
	if err != nil {
		return err
	}

	// Get account address and nonce
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transaction options
	gasPrice, err := parseGasPrice(c.String(FlagGasPrice))
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = c.Uint64(FlagGasLimit)

	// Calculate value if callback is requested
	callbackAddress := common.HexToAddress(c.String(FlagCallbackAddress))
	callbackGasLimit := big.NewInt(int64(c.Uint64(FlagCallbackGas)))
	
	if callbackAddress != common.HexToAddress("0x0") {
		// Calculate required payment for callback gas
		callbackGasPrice := big.NewInt(10e9) // 10 gwei (contract constant)
		requiredValue := new(big.Int).Mul(callbackGasLimit, callbackGasPrice)
		auth.Value = requiredValue
		
		fmt.Printf("Callback requested. Sending %s wei for gas deposit\n", requiredValue.String())
	} else {
		auth.Value = big.NewInt(0)
		callbackGasLimit = big.NewInt(0)
	}

	// Request randomness
	fmt.Printf("Requesting randomness...\n")
	fmt.Printf("  From: %s\n", fromAddress.Hex())
	fmt.Printf("  Contract: %s\n", c.String(FlagContractAddress))
	fmt.Printf("  Callback: %s\n", callbackAddress.Hex())
	fmt.Printf("  Callback Gas: %s\n", callbackGasLimit.String())

	tx, err := contract.RequestRandomness(auth, callbackAddress, callbackGasLimit)
	if err != nil {
		return fmt.Errorf("failed to request randomness: %w", err)
	}

	fmt.Printf("\nTransaction submitted!\n")
	fmt.Printf("  Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  Gas Price: %s gwei\n", new(big.Int).Div(tx.GasPrice(), big.NewInt(1e9)).String())
	fmt.Printf("  Gas Limit: %d\n", tx.Gas())

	// Wait for transaction receipt
	fmt.Printf("\nWaiting for confirmation...\n")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	fmt.Printf("Transaction confirmed!\n")
	fmt.Printf("  Block: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)

	// Parse events to get task hash
	var taskHash common.Hash
	for _, log := range receipt.Logs {
		if len(log.Topics) > 0 && log.Topics[0] == crypto.Keccak256Hash([]byte("RandomnessRequested(bytes32,address,address,uint256)")) {
			if len(log.Topics) > 1 {
				taskHash = log.Topics[1]
				fmt.Printf("\nRandomness requested successfully!\n")
				fmt.Printf("  Task Hash: %s\n", taskHash.Hex())
				
				if !c.Bool(FlagWaitForResult) {
					fmt.Printf("\nUse 'omnivrf-client get --task-hash %s' to check the result\n", taskHash.Hex())
				}
			}
		}
	}

	// If wait flag is set, listen for the result
	if c.Bool(FlagWaitForResult) && taskHash != (common.Hash{}) {
		fmt.Printf("\nWaiting for randomness result (timeout: %s)...\n", c.Duration(FlagTimeout).String())
		
		err = waitForRandomnessResult(c, client, contract, taskHash)
		if err != nil {
			return fmt.Errorf("error waiting for result: %w", err)
		}
	}

	return nil
}

func getRandomness(c *cli.Context) error {
	client, _, err := getClient(c)
	if err != nil {
		return err
	}
	defer client.Close()

	contract, err := getContract(c, client)
	if err != nil {
		return err
	}

	// Parse task hash
	taskHashHex := c.String(FlagTaskHash)
	taskHashHex = strings.TrimPrefix(taskHashHex, "0x")
	
	var taskHash [32]byte
	if len(taskHashHex) != 64 {
		return fmt.Errorf("invalid task hash length: expected 32 bytes (64 hex chars)")
	}
	
	taskHashBytes := common.FromHex(taskHashHex)
	copy(taskHash[:], taskHashBytes)

	// Get randomness
	result, err := contract.GetRandomness(&bind.CallOpts{}, taskHash)
	if err != nil {
		return fmt.Errorf("failed to get randomness: %w", err)
	}

	fmt.Printf("Task Hash: 0x%x\n", taskHash)
	fmt.Printf("Fulfilled: %v\n", result.Fulfilled)
	fmt.Printf("Randomness: %s\n", result.Randomness.String())

	if !result.Fulfilled {
		fmt.Printf("\nNote: This request has not been fulfilled yet. Please wait for the VRF operators to process it.\n")
	}

	return nil
}

func withdrawCallbackGas(c *cli.Context) error {
	client, privateKey, err := getClient(c)
	if err != nil {
		return err
	}
	defer client.Close()

	contract, err := getContract(c, client)
	if err != nil {
		return err
	}

	// Parse amount
	amount, ok := new(big.Int).SetString(c.String(FlagAmount), 10)
	if !ok {
		return fmt.Errorf("invalid amount: %s", c.String(FlagAmount))
	}

	// Get account address and nonce
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transaction options
	gasPrice, err := parseGasPrice(c.String(FlagGasPrice))
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = c.Uint64(FlagGasLimit)

	// Withdraw gas
	fmt.Printf("Withdrawing callback gas...\n")
	fmt.Printf("  From: %s\n", fromAddress.Hex())
	fmt.Printf("  Amount: %s wei\n", amount.String())

	tx, err := contract.WithdrawCallbackGas(auth, amount)
	if err != nil {
		return fmt.Errorf("failed to withdraw gas: %w", err)
	}

	fmt.Printf("\nTransaction submitted!\n")
	fmt.Printf("  Hash: %s\n", tx.Hash().Hex())

	// Wait for confirmation
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	fmt.Printf("\nWithdrawal successful!\n")
	fmt.Printf("  Block: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)

	return nil
}

func getContractInfo(c *cli.Context) error {
	client, _, err := getClient(c)
	if err != nil {
		return err
	}
	defer client.Close()

	contract, err := getContract(c, client)
	if err != nil {
		return err
	}

	// Get contract constants
	callbackGasPrice, err := contract.CALLBACKGASPRICE(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get callback gas price: %w", err)
	}

	minCallbackGasLimit, err := contract.MINCALLBACKGASLIMIT(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get min callback gas limit: %w", err)
	}

	maxCallbackGasLimit, err := contract.MAXCALLBACKGASLIMIT(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get max callback gas limit: %w", err)
	}

	taskMailbox, err := contract.TaskMailbox(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get task mailbox: %w", err)
	}

	// Get current account's callback gas deposits
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(c.String(FlagPrivateKey), "0x"))
	if err == nil {
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		deposits, err := contract.CallbackGasDeposits(&bind.CallOpts{}, address)
		if err == nil {
			fmt.Printf("Your Callback Gas Deposits: %s wei\n", deposits.String())
			fmt.Printf("---\n")
		}
	}

	fmt.Printf("OmniVRF Contract Information\n")
	fmt.Printf("---\n")
	fmt.Printf("Contract Address: %s\n", c.String(FlagContractAddress))
	fmt.Printf("Task Mailbox: %s\n", taskMailbox.Hex())
	fmt.Printf("Callback Gas Price: %s wei (%s gwei)\n", 
		callbackGasPrice.String(), 
		new(big.Int).Div(callbackGasPrice, big.NewInt(1e9)).String())
	fmt.Printf("Min Callback Gas Limit: %s\n", minCallbackGasLimit.String())
	fmt.Printf("Max Callback Gas Limit: %s\n", maxCallbackGasLimit.String())

	return nil
}

func parseGasPrice(gasPriceGwei string) (*big.Int, error) {
	gasPrice, ok := new(big.Int).SetString(gasPriceGwei, 10)
	if !ok {
		return nil, fmt.Errorf("invalid gas price: %s", gasPriceGwei)
	}
	
	// Convert from gwei to wei
	return new(big.Int).Mul(gasPrice, big.NewInt(1e9)), nil
}

func waitForRandomnessResult(c *cli.Context, client *ethclient.Client, contract *OmniVRF.OmniVRF, taskHash common.Hash) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Duration(FlagTimeout))
	defer cancel()

	// Get the contract ABI for event parsing
	contractABI, err := abi.JSON(strings.NewReader(OmniVRF.OmniVRFABI))
	if err != nil {
		return fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	// Convert RPC URL to WebSocket URL if needed
	rpcURL := c.String(FlagRPCURL)
	wsURL := rpcURL
	if strings.HasPrefix(rpcURL, "http://") {
		wsURL = strings.Replace(rpcURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(rpcURL, "https://") {
		wsURL = strings.Replace(rpcURL, "https://", "wss://", 1)
	}

	// Try to connect with WebSocket for real-time updates
	wsClient, err := ethclient.DialContext(ctx, wsURL)
	if err != nil {
		// Fallback to polling if WebSocket is not available
		fmt.Printf("WebSocket not available, falling back to polling...\n")
		return pollForResult(ctx, contract, taskHash)
	}
	defer wsClient.Close()

	// Create filter query for RandomnessFulfilled events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(c.String(FlagContractAddress))},
		Topics: [][]common.Hash{
			{contractABI.Events["RandomnessFulfilled"].ID},
			{taskHash}, // Filter for our specific task hash
		},
	}

	// Subscribe to logs
	logs := make(chan types.Log)
	sub, err := wsClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		// Fallback to polling if subscription fails
		fmt.Printf("WebSocket subscription failed, falling back to polling...\n")
		return pollForResult(ctx, contract, taskHash)
	}
	defer sub.Unsubscribe()

	// Listen for events
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for randomness result")
			
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
			
		case log := <-logs:
			// Parse the event
			event := new(struct {
				TaskHash   [32]byte
				Randomness *big.Int
			})
			
			err := contractABI.UnpackIntoInterface(event, "RandomnessFulfilled", log.Data)
			if err != nil {
				return fmt.Errorf("failed to unpack event: %w", err)
			}

			fmt.Printf("\n✓ Randomness fulfilled!\n")
			fmt.Printf("  Task Hash: 0x%x\n", event.TaskHash)
			fmt.Printf("  Randomness: %s\n", event.Randomness.String())
			fmt.Printf("  Block: %d\n", log.BlockNumber)
			fmt.Printf("  Transaction: %s\n", log.TxHash.Hex())
			
			return nil
		}
	}
}

func pollForResult(ctx context.Context, contract *OmniVRF.OmniVRF, taskHash common.Hash) error {
	// Convert taskHash to [32]byte
	var taskHashArray [32]byte
	copy(taskHashArray[:], taskHash.Bytes())

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinIndex := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for randomness result")
			
		case <-ticker.C:
			// Check if the request has been fulfilled
			result, err := contract.GetRandomness(&bind.CallOpts{Context: ctx}, taskHashArray)
			if err != nil {
				// Continue polling on error
				continue
			}

			if result.Fulfilled {
				fmt.Printf("\n✓ Randomness fulfilled!\n")
				fmt.Printf("  Task Hash: %s\n", taskHash.Hex())
				fmt.Printf("  Randomness: %s\n", result.Randomness.String())
				return nil
			}

			// Show spinner
			fmt.Printf("\r%s Polling for result...", spinChars[spinIndex])
			spinIndex = (spinIndex + 1) % len(spinChars)
		}
	}
}