// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../libraries/Callback.sol";

// Stub contract used in tests
contract CallbackContractStub is IRandomBeaconConsumer {
    uint256 public _lastEntry;
    uint256 public _blockNumber;

    function __beaconCallback(uint256 entry, uint256 blockNumber)
        external
        override
    {
        _lastEntry = entry;
        _blockNumber = blockNumber;
    }
}
