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
    }

    event CallbackExecuted(uint256 entry, uint256 entrySubmittedBlock);
    event CallbackFailed(uint256 entry, uint256 entrySubmittedBlock);

    /// @notice Sets callback contract.
    /// @param callbackContract Callback contract.
    function setCallbackContract(
      Data storage self,
      IRandomBeaconConsumer callbackContract
    ) internal {
      if (address(callbackContract) != address(0)) {
          self.callbackContract = callbackContract;
      }
    }

    /// @notice Executes consumer specified callback for the relay entry request.
    /// @param entry The generated random number.
    /// @param callbackGasLimit Callback gas limit.
    function executeCallback(Data storage self, uint256 entry, uint256 callbackGasLimit) internal {
      if (address(self.callbackContract) != address(0)) {
        try self.callbackContract.__beaconCallback{gas: callbackGasLimit}(entry, block.number) {
          emit CallbackExecuted(entry, block.number);
        } catch {
          emit CallbackFailed(entry, block.number);
        }
      }
    }
}
