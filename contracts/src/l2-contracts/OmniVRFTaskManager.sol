// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";
import {IAVSTaskHook} from "@hourglass-monorepo/src/interfaces/avs/l2/IAVSTaskHook.sol";
import {ITaskMailbox, ITaskMailboxTypes} from "@hourglass-monorepo/src/interfaces/core/ITaskMailbox.sol";

import {IOmniVRFConsumer} from "@project/interfaces/IOmniVRFConsumer.sol";

/**
 * @title OmniVRFTaskManager
 * @notice Verifiable Random Function service implemented as an Hourglass AVS
 * @dev Implements AVSTaskHook for lifecycle management and interacts with TaskMailbox
 */
contract OmniVRFTaskManager is IAVSTaskHook {
    
    // ============ STRUCTS ============
    
    /**
     * @notice Structure to store randomness request details
     */
    struct RandomnessRequest {
        address requester;           // Address that requested randomness
        address callbackContract;   // Optional contract to callback with result (address(0) if none)
        uint256 callbackGasLimit;   // Gas limit for callback execution
        uint256 seed;               // Deterministic seed for VRF computation
        uint256 blockNumber;        // Block number when request was made
        bool fulfilled;             // Whether the request has been fulfilled
        uint256 result;             // The generated random number (0 if not fulfilled)
    }

    /**
     * @notice Structure for task data sent to operators
     */
    struct VRFTaskData {
        bytes32 taskHash;           // Task hash from TaskMailbox (used as unique identifier)
        uint256 seed;               // Seed for VRF computation
    }

    // ============ STATE VARIABLES ============

    /// @notice Reference to the TaskMailbox contract for creating tasks
    ITaskMailbox public immutable taskMailbox;

    /// @notice Mapping from task hash to randomness request details
    mapping(bytes32 => RandomnessRequest) public requests;
    
    /// @notice Mapping from requester address to their gas deposits for callbacks
    mapping(address => uint256) public callbackGasDeposits;
    
    /// @notice Gas price used for callback cost calculations (10 gwei)
    uint256 public constant CALLBACK_GAS_PRICE = 10 gwei;
    
    /// @notice Minimum callback gas limit to prevent griefing
    uint256 public constant MIN_CALLBACK_GAS_LIMIT = 100000;
    
    /// @notice Maximum callback gas limit to prevent excessive costs
    uint256 public constant MAX_CALLBACK_GAS_LIMIT = 500000;

    // ============ EVENTS ============

    /**
     * @notice Emitted when randomness is requested
     * @param taskHash Unique task hash identifier from TaskMailbox
     * @param requester Address that requested randomness
     * @param callbackContract Optional callback contract address
     * @param seed Deterministic seed used for VRF computation
     */
    event RandomnessRequested(
        bytes32 indexed taskHash,
        address indexed requester,
        address callbackContract,
        uint256 seed
    );

    /**
     * @notice Emitted when randomness is fulfilled
     * @param taskHash Unique task hash identifier
     * @param randomness The generated random number
     */
    event RandomnessFulfilled(
        bytes32 indexed taskHash,
        uint256 randomness
    );

    /**
     * @notice Emitted when a callback execution fails
     * @param taskHash Unique task hash identifier
     * @param callbackContract Address of the contract that failed
     * @param reason Failure reason
     */
    event CallbackFailed(
        bytes32 indexed taskHash,
        address indexed callbackContract,
        string reason
    );

    // ============ ERRORS ============

    error InsufficientCallbackGas();
    error InvalidCallbackGasLimit();
    error RequestNotFound();
    error RequestAlreadyFulfilled();
    error InvalidTaskData();
    error UnauthorizedCaller();

    // ============ CONSTRUCTOR ============

    /**
     * @notice Initialize the OmniVRF Task Manager
     * @param _taskMailbox Address of the TaskMailbox contract
     */
    constructor(address _taskMailbox) {
        taskMailbox = ITaskMailbox(_taskMailbox);
    }

    // ============ EXTERNAL FUNCTIONS ============

    /**
     * @notice Request verifiable randomness
     * @param callbackContract Optional contract to callback with result (address(0) for no callback)
     * @param callbackGasLimit Gas limit for callback execution (ignored if callbackContract is address(0))
     * @return taskHash Unique task hash identifier from TaskMailbox
     */
    function requestRandomness(
        address callbackContract,
        uint256 callbackGasLimit
    ) external payable returns (bytes32 taskHash) {
        
        // Validate callback gas limit if callback is requested
        if (callbackContract != address(0)) {
            if (callbackGasLimit < MIN_CALLBACK_GAS_LIMIT || callbackGasLimit > MAX_CALLBACK_GAS_LIMIT) {
                revert InvalidCallbackGasLimit();
            }
            
            // Calculate required gas deposit
            uint256 requiredDeposit = callbackGasLimit * CALLBACK_GAS_PRICE;
            if (msg.value < requiredDeposit) {
                revert InsufficientCallbackGas();
            }
            
            // Store gas deposit
            callbackGasDeposits[msg.sender] += msg.value;
        }
        
        // Generate deterministic seed
        uint256 seed = uint256(keccak256(abi.encodePacked(
            msg.sender,
            block.prevrandao,
            block.timestamp,
            block.number
        )));
        
        // Create task in TaskMailbox first to get taskHash
        // For MVP, use a simple operator set with ID 1
        OperatorSet memory operatorSet = OperatorSet({
            avs: address(this),
            id: uint32(1)
        });
        
        // Create task data for operators (taskHash will be available to performer through task context)
        VRFTaskData memory taskData = VRFTaskData({
            taskHash: bytes32(0), // Placeholder - performers will get actual taskHash from TaskMailbox
            seed: seed
        });
        
        ITaskMailboxTypes.TaskParams memory taskParams = ITaskMailboxTypes.TaskParams({
            refundCollector: msg.sender,
            avsFee: 0,  // No fee for MVP
            executorOperatorSet: operatorSet,
            payload: abi.encode(taskData)
        });
        
        taskHash = taskMailbox.createTask(taskParams);
        
        // Store request with the actual taskHash
        requests[taskHash] = RandomnessRequest({
            requester: msg.sender,
            callbackContract: callbackContract,
            callbackGasLimit: callbackGasLimit,
            seed: seed,
            blockNumber: block.number,
            fulfilled: false,
            result: 0
        });
        
        emit RandomnessRequested(taskHash, msg.sender, callbackContract, seed);
    }

    /**
     * @notice Get randomness result for a given task hash
     * @param taskHash The task hash identifier
     * @return fulfilled Whether the request has been fulfilled
     * @return randomness The generated random number (0 if not fulfilled)
     */
    function getRandomness(bytes32 taskHash) external view returns (bool fulfilled, uint256 randomness) {
        RandomnessRequest storage request = requests[taskHash];
        if (request.requester == address(0)) {
            revert RequestNotFound();
        }
        
        return (request.fulfilled, request.result);
    }

    /**
     * @notice Withdraw unused callback gas deposits
     * @param amount Amount to withdraw
     */
    function withdrawCallbackGas(uint256 amount) external {
        require(callbackGasDeposits[msg.sender] >= amount, "Insufficient deposits");
        
        callbackGasDeposits[msg.sender] -= amount;
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }

    // ============ HOURGLASS TASK LIFECYCLE HOOKS ============

    /**
     * @notice Validate task creation parameters
     * @param caller Address creating the task
     * @param taskParams Task parameters
     */
    function validatePreTaskCreation(
        address caller,
        ITaskMailboxTypes.TaskParams memory taskParams
    ) external view override {
        // Only allow this contract to create VRF tasks
        if (caller != address(this)) {
            revert UnauthorizedCaller();
        }
        
        // Validate payload can be decoded as VRFTaskData
        try this.decodeTaskData(taskParams.payload) {
            // Validation successful
        } catch {
            revert InvalidTaskData();
        }
    }

    /**
     * @notice Handle post-task creation
     * @param taskHash Hash of the created task
     */
    function handlePostTaskCreation(bytes32 taskHash) external override {
        // No additional handling needed for MVP
    }

    /**
     * @notice Validate task result before submission
     * @param caller Address submitting the result
     * @param taskHash Hash of the task
     * @param cert Certificate data
     * @param result Task result data
     */
    function validatePreTaskResultSubmission(
        address caller,
        bytes32 taskHash,
        bytes memory cert,
        bytes memory result
    ) external view override {
        // For MVP, we trust the Hourglass aggregator
        // In production, we would verify VRF proofs here
    }

    /**
     * @notice Handle completed task results
     * @param taskHash Hash of the completed task
     */
    function handlePostTaskResultSubmission(bytes32 taskHash) external override {
        // Get task result from TaskMailbox
        // Note: This is a simplified approach for MVP
        // In production, we would get the actual result from the task
        
        // For now, we'll generate a mock result
        // TODO: Integrate with actual Hourglass result retrieval
        uint256 mockRandomness = uint256(keccak256(abi.encodePacked(taskHash, block.timestamp)));
        
        // Fulfill the request using the taskHash
        _fulfillRequest(taskHash, mockRandomness);
    }

    /**
     * @notice Calculate task fee for operators
     * @param operatorSet The operator set
     * @param payload Task payload
     * @return fee Task fee amount
     */
    function calculateTaskFee(
        OperatorSet memory operatorSet,
        bytes memory payload
    ) external view override returns (uint96 fee) {
        // No fees for MVP
        return 0;
    }

    // ============ INTERNAL FUNCTIONS ============

    /**
     * @notice Fulfill a randomness request
     * @param taskHash The task hash to fulfill
     * @param randomness The generated random number
     */
    function _fulfillRequest(bytes32 taskHash, uint256 randomness) internal {
        RandomnessRequest storage request = requests[taskHash];
        
        if (request.requester == address(0)) {
            revert RequestNotFound();
        }
        
        if (request.fulfilled) {
            revert RequestAlreadyFulfilled();
        }
        
        // Store result
        request.result = randomness;
        request.fulfilled = true;
        
        emit RandomnessFulfilled(taskHash, randomness);
        
        // Execute callback if requested
        if (request.callbackContract != address(0)) {
            _executeCallback(taskHash);
        }
    }

    /**
     * @notice Execute callback to consumer contract
     * @param taskHash The task hash to callback with
     */
    function _executeCallback(bytes32 taskHash) internal {
        RandomnessRequest memory request = requests[taskHash];
        
        // Calculate and deduct gas costs
        uint256 gasCost = request.callbackGasLimit * CALLBACK_GAS_PRICE;
        callbackGasDeposits[request.requester] -= gasCost;
        
        // Execute callback with gas limit protection
        try IOmniVRFConsumer(request.callbackContract).fulfillRandomness{
            gas: request.callbackGasLimit
        }(taskHash, request.result) {
            // Callback successful
        } catch Error(string memory reason) {
            // Refund gas on failure
            callbackGasDeposits[request.requester] += gasCost;
            emit CallbackFailed(taskHash, request.callbackContract, reason);
        } catch {
            // Refund gas on failure
            callbackGasDeposits[request.requester] += gasCost;
            emit CallbackFailed(taskHash, request.callbackContract, "Unknown error");
        }
    }

    // ============ UTILITY FUNCTIONS ============

    /**
     * @notice Decode task data (external for validation)
     * @param payload The encoded task data
     * @return taskData Decoded VRF task data
     */
    function decodeTaskData(bytes memory payload) external pure returns (VRFTaskData memory taskData) {
        return abi.decode(payload, (VRFTaskData));
    }
}