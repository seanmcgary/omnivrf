// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Test, console} from "forge-std/Test.sol";
import {OperatorSet} from "@eigenlayer-contracts/src/contracts/libraries/OperatorSetLib.sol";
import {ITaskMailbox, ITaskMailboxTypes} from "@hourglass-monorepo/src/interfaces/core/ITaskMailbox.sol";

import {OmniVRF} from "@project/l2-contracts/OmniVRF.sol";
import {IOmniVRF} from "@project/interfaces/IOmniVRF.sol";

/**
 * @title MockTaskMailbox
 * @notice Mock TaskMailbox contract for testing
 */
contract MockTaskMailbox {
    mapping(bytes32 => bytes) public tasks;
    uint256 private taskCounter = 1;
    
    function createTask(ITaskMailboxTypes.TaskParams memory taskParams) external returns (bytes32) {
        bytes32 taskHash = keccak256(abi.encodePacked(taskCounter, taskParams.payload));
        tasks[taskHash] = taskParams.payload;
        taskCounter++;
        return taskHash;
    }
    
    function getTaskInfo(bytes32 taskHash) external view returns (
        ITaskMailboxTypes.Task memory task
    ) {
        task.payload = tasks[taskHash];
        task.status = ITaskMailboxTypes.TaskStatus.CREATED;
        return task;
    }
    
    function getTaskResult(bytes32 /* taskHash */) external pure returns (bytes memory) {
        return ""; // Mock empty result
    }
}

/**
 * @title MockVRFConsumer
 * @notice Mock contract for testing VRF callbacks
 */
contract MockVRFConsumer is IOmniVRF {
    bytes32 public lastTaskHash;
    uint256 public lastRandomness;
    bool public shouldRevert;
    
    function fulfillRandomness(bytes32 taskHash, uint256 randomness) external override {
        if (shouldRevert) {
            revert("Mock revert");
        }
        
        lastTaskHash = taskHash;
        lastRandomness = randomness;
    }
    
    function setShouldRevert(bool _shouldRevert) external {
        shouldRevert = _shouldRevert;
    }
}

/**
 * @title OmniVRFTest
 * @notice Test suite for OmniVRF contract
 */
contract OmniVRFTest is Test {
    OmniVRF public omniVRF;
    MockVRFConsumer public mockConsumer;
    MockTaskMailbox public mockTaskMailbox;
    
    address public constant USER = address(0x1234);
    address public constant OPERATOR = address(0x5678);
    
    uint256 public constant CALLBACK_GAS_LIMIT = 200000;
    uint256 public constant CALLBACK_GAS_PRICE = 10 gwei;
    
    // Events to test
    event RandomnessRequested(
        bytes32 indexed taskHash,
        address indexed requester,
        address callbackContract,
        uint256 seed
    );
    
    event RandomnessFulfilled(
        bytes32 indexed taskHash,
        uint256 randomness
    );
    
    event CallbackFailed(
        bytes32 indexed taskHash,
        address indexed callbackContract,
        string reason
    );

    function setUp() public {
        // Deploy contracts
        mockTaskMailbox = new MockTaskMailbox();
        omniVRF = new OmniVRF(address(mockTaskMailbox));
        mockConsumer = new MockVRFConsumer();
        
        // Setup test accounts
        vm.deal(USER, 10 ether);
    }

    function testRequestRandomnessWithoutCallback() public {
        vm.startPrank(USER);
        
        // Test request without callback
        bytes32 taskHash = omniVRF.requestRandomness(address(0), 0);
        
        // Verify request was created
        assertTrue(taskHash != bytes32(0));
        
        // Check request details
        (
            address requester,
            address callbackContract,
            uint256 callbackGasLimit,
            uint256 seed,
            uint256 blockNumber,
            bool fulfilled,
            uint256 result
        ) = omniVRF.requests(taskHash);
        
        assertEq(requester, USER);
        assertEq(callbackContract, address(0));
        assertEq(callbackGasLimit, 0);
        assertTrue(seed > 0); // Should have generated a seed
        assertEq(blockNumber, block.number);
        assertFalse(fulfilled);
        assertEq(result, 0);
        
        vm.stopPrank();
    }

    function testRequestRandomnessWithCallback() public {
        vm.startPrank(USER);
        
        uint256 requiredDeposit = CALLBACK_GAS_LIMIT * CALLBACK_GAS_PRICE;
        
        bytes32 taskHash = omniVRF.requestRandomness{value: requiredDeposit}(
            address(mockConsumer),
            CALLBACK_GAS_LIMIT
        );
        
        // Verify request was created
        assertTrue(taskHash != bytes32(0));
        
        // Check gas deposit
        assertEq(omniVRF.callbackGasDeposits(USER), requiredDeposit);
        
        // Check request details
        (
            address requester,
            address callbackContract,
            uint256 callbackGasLimit,
            ,,,
        ) = omniVRF.requests(taskHash);
        
        assertEq(requester, USER);
        assertEq(callbackContract, address(mockConsumer));
        assertEq(callbackGasLimit, CALLBACK_GAS_LIMIT);
        
        vm.stopPrank();
    }

    function testRequestRandomnessInsufficientGas() public {
        vm.startPrank(USER);
        
        uint256 requiredDeposit = CALLBACK_GAS_LIMIT * CALLBACK_GAS_PRICE;
        uint256 insufficientDeposit = requiredDeposit - 1;
        
        // Should revert with insufficient gas
        vm.expectRevert(OmniVRF.InsufficientCallbackGas.selector);
        omniVRF.requestRandomness{value: insufficientDeposit}(
            address(mockConsumer),
            CALLBACK_GAS_LIMIT
        );
        
        vm.stopPrank();
    }

    function testRequestRandomnessInvalidGasLimit() public {
        vm.startPrank(USER);
        
        // Test gas limit too low
        vm.expectRevert(OmniVRF.InvalidCallbackGasLimit.selector);
        omniVRF.requestRandomness{value: 1 ether}(
            address(mockConsumer),
            50000 // Below MIN_CALLBACK_GAS_LIMIT
        );
        
        // Test gas limit too high
        vm.expectRevert(OmniVRF.InvalidCallbackGasLimit.selector);
        omniVRF.requestRandomness{value: 1 ether}(
            address(mockConsumer),
            1000000 // Above MAX_CALLBACK_GAS_LIMIT
        );
        
        vm.stopPrank();
    }

    function testGetRandomness() public {
        vm.startPrank(USER);
        
        bytes32 taskHash = omniVRF.requestRandomness(address(0), 0);
        
        // Check unfulfilled request
        (bool fulfilled, uint256 randomness) = omniVRF.getRandomness(taskHash);
        assertFalse(fulfilled);
        assertEq(randomness, 0);
        
        vm.stopPrank();
    }

    function testGetRandomnessInvalidRequest() public {
        vm.expectRevert(OmniVRF.RequestNotFound.selector);
        omniVRF.getRandomness(bytes32(0));
    }

    function testWithdrawCallbackGas() public {
        vm.startPrank(USER);
        
        uint256 deposit = CALLBACK_GAS_LIMIT * CALLBACK_GAS_PRICE;
        
        // Make a request with callback
        omniVRF.requestRandomness{value: deposit}(
            address(mockConsumer),
            CALLBACK_GAS_LIMIT
        );
        
        // Verify deposit
        assertEq(omniVRF.callbackGasDeposits(USER), deposit);
        
        // Withdraw half
        uint256 withdrawAmount = deposit / 2;
        uint256 balanceBefore = USER.balance;
        
        omniVRF.withdrawCallbackGas(withdrawAmount);
        
        // Verify withdrawal
        assertEq(omniVRF.callbackGasDeposits(USER), deposit - withdrawAmount);
        assertEq(USER.balance, balanceBefore + withdrawAmount);
        
        vm.stopPrank();
    }

    function testWithdrawCallbackGasInsufficientDeposits() public {
        vm.startPrank(USER);
        
        // Try to withdraw without deposits
        vm.expectRevert("Insufficient deposits");
        omniVRF.withdrawCallbackGas(1 ether);
        
        vm.stopPrank();
    }

    function testTaskDataDecoding() public {
        OmniVRF.VRFTaskData memory taskData = OmniVRF.VRFTaskData({
            taskHash: bytes32(uint256(123)),
            seed: 456789
        });
        
        bytes memory encoded = abi.encode(taskData);
        OmniVRF.VRFTaskData memory decoded = omniVRF.decodeTaskData(encoded);
        
        assertEq(decoded.taskHash, bytes32(uint256(123)));
        assertEq(decoded.seed, 456789);
    }

    function testValidatePreTaskCreation() public {
        // Create valid task data
        OmniVRF.VRFTaskData memory taskData = OmniVRF.VRFTaskData({
            taskHash: bytes32(uint256(1)),
            seed: 12345
        });
        
        OperatorSet memory operatorSet = OperatorSet({
            avs: address(omniVRF),
            id: uint32(1)
        });
        
        ITaskMailboxTypes.TaskParams memory taskParams = ITaskMailboxTypes.TaskParams({
            refundCollector: address(omniVRF),
            avsFee: 0,
            executorOperatorSet: operatorSet,
            payload: abi.encode(taskData)
        });
        
        // Should not revert when called by the contract itself
        vm.prank(address(omniVRF));
        omniVRF.validatePreTaskCreation(address(omniVRF), taskParams);
        
        // Should revert when called by unauthorized address
        vm.prank(USER);
        vm.expectRevert(OmniVRF.UnauthorizedCaller.selector);
        omniVRF.validatePreTaskCreation(USER, taskParams);
    }

    function testValidatePreTaskCreationInvalidData() public {
        // Create invalid task data
        OperatorSet memory operatorSet2 = OperatorSet({
            avs: address(omniVRF),
            id: uint32(1)
        });
        
        ITaskMailboxTypes.TaskParams memory taskParams = ITaskMailboxTypes.TaskParams({
            refundCollector: address(omniVRF),
            avsFee: 0,
            executorOperatorSet: operatorSet2,
            payload: hex"deadbeef" // Invalid bytes that can't be decoded as VRFTaskData
        });
        
        vm.prank(address(omniVRF));
        vm.expectRevert(OmniVRF.InvalidTaskData.selector);
        omniVRF.validatePreTaskCreation(address(omniVRF), taskParams);
    }


    function testMultipleRequests() public {
        vm.startPrank(USER);
        
        // Make multiple requests
        bytes32 taskHash1 = omniVRF.requestRandomness(address(0), 0);
        bytes32 taskHash2 = omniVRF.requestRandomness(address(0), 0);
        bytes32 taskHash3 = omniVRF.requestRandomness(address(0), 0);
        
        // Verify all task hashes are different (unique)
        assertTrue(taskHash1 != taskHash2);
        assertTrue(taskHash2 != taskHash3);
        assertTrue(taskHash1 != taskHash3);
        
        vm.stopPrank();
    }

    function testSeedUniqueness() public {
        vm.startPrank(USER);
        
        bytes32 taskHash1 = omniVRF.requestRandomness(address(0), 0);
        
        // Move to next block to ensure different prevrandao
        vm.roll(block.number + 1);
        
        bytes32 taskHash2 = omniVRF.requestRandomness(address(0), 0);
        
        // Get seeds from both requests
        (, , , uint256 seed1, , ,) = omniVRF.requests(taskHash1);
        (, , , uint256 seed2, , ,) = omniVRF.requests(taskHash2);
        
        // Seeds should be different
        assertTrue(seed1 != seed2);
        
        vm.stopPrank();
    }

    function testContractConstants() public {
        assertEq(omniVRF.CALLBACK_GAS_PRICE(), 10 gwei);
        assertEq(omniVRF.MIN_CALLBACK_GAS_LIMIT(), 100000);
        assertEq(omniVRF.MAX_CALLBACK_GAS_LIMIT(), 500000);
    }

    // Helper function to simulate task completion
    function simulateTaskCompletion(bytes32 taskHash) internal {
        // This would normally be called by the Hourglass framework
        // For testing, we'll call it directly
        omniVRF.handlePostTaskResultSubmission(taskHash);
    }

    function testTaskCompletionFlow() public {
        vm.startPrank(USER);
        
        uint256 deposit = CALLBACK_GAS_LIMIT * CALLBACK_GAS_PRICE;
        
        // Request with callback
        bytes32 taskHash = omniVRF.requestRandomness{value: deposit}(
            address(mockConsumer),
            CALLBACK_GAS_LIMIT
        );
        
        vm.stopPrank();
        
        // Verify request is not fulfilled
        (bool fulfilled,) = omniVRF.getRandomness(taskHash);
        assertFalse(fulfilled);
        
        // Simulate task completion
        vm.expectEmit(true, false, false, false);
        emit RandomnessFulfilled(taskHash, 0); // randomness value is dynamic
        
        simulateTaskCompletion(taskHash);
        
        // Verify request is now fulfilled
        (fulfilled,) = omniVRF.getRandomness(taskHash);
        assertTrue(fulfilled);
        
        // Verify callback was executed
        assertEq(mockConsumer.lastTaskHash(), taskHash);
        assertTrue(mockConsumer.lastRandomness() > 0);
    }

    function testCallbackFailureHandling() public {
        vm.startPrank(USER);
        
        uint256 deposit = CALLBACK_GAS_LIMIT * CALLBACK_GAS_PRICE;
        uint256 initialBalance = omniVRF.callbackGasDeposits(USER);
        
        // Set mock consumer to revert
        mockConsumer.setShouldRevert(true);
        
        // Request with callback
        bytes32 taskHash = omniVRF.requestRandomness{value: deposit}(
            address(mockConsumer),
            CALLBACK_GAS_LIMIT
        );
        
        vm.stopPrank();
        
        // Simulate task completion - should handle callback failure
        vm.expectEmit(true, true, false, false);
        emit CallbackFailed(taskHash, address(mockConsumer), "Mock revert");
        
        simulateTaskCompletion(taskHash);
        
        // Verify gas was refunded
        assertEq(omniVRF.callbackGasDeposits(USER), initialBalance + deposit);
        
        // Verify request is still fulfilled despite callback failure
        (bool fulfilled,) = omniVRF.getRandomness(taskHash);
        assertTrue(fulfilled);
    }
}