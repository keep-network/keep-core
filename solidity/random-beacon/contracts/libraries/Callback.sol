// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

interface IRandomBeaconConsumer {
    /// @notice Receives relay entry produced by Keep Random Beacon. This function
    /// should be called only by Keep Random Beacon.
    ///
    /// @param relayEntry Relay entry (random number) produced by Keep Random
    ///                   Beacon.
    /// @param blockNumber Block number at which the relay entry was submitted
    ///                    to the chain.
    function __beaconCallback(uint256 relayEntry, uint256 blockNumber) external;
}

/// @title Callback library
/// @dev Library for handling calls to random beacon consumer.
library Callback {
    struct Data {
      IRandomBeaconConsumer callbackContract;
      uint256 entrySubmittedBlock;
    }

    event CallbackExecuted(uint256 entry, uint256 entrySubmittedBlock);

    /// @notice Executes customer specified callback for the relay entry request.
    /// @param entry The generated random number.
    /// @param entryValidityBlocks Entry submitted is only valid for a certain number of blocks
    function executeCallback(Data storage self, uint256 entry, uint256 entryValidityBlocks) internal {
        require(
            address(self.callbackContract) != address(0),
            "Callback contract must be set"
        );

        require(
            block.number <= self.entrySubmittedBlock + entryValidityBlocks,
            "Entry is no longer valid"
        );

        self.callbackContract.__beaconCallback(entry, self.entrySubmittedBlock);

        emit CallbackExecuted(entry, self.entrySubmittedBlock);
    }
}
