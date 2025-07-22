// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

/**
 * @title IOmniVRF
 * @notice Interface for contracts that want to receive VRF results via callbacks
 * @dev Contracts implementing this interface can request VRF with automatic callbacks
 */
interface IOmniVRF {
    /**
     * @notice Called by OmniVRF when randomness is fulfilled
     * @param taskHash The unique task hash identifier from TaskMailbox
     * @param randomness The generated random number
     * @dev This function must be implemented by consuming contracts
     *      Gas limit for this function is specified during the request
     */
    function fulfillRandomness(bytes32 taskHash, uint256 randomness) external;
}