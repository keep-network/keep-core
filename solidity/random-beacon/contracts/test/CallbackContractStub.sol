// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../libraries/Callback.sol";

// Stub contract used in tests
contract CallbackContractStub is IRandomBeaconConsumer {
    uint256 public lastEntry;
    uint256 public blockNumber;

    function __beaconCallback(uint256 _lastEntry, uint256 _blockNumber)
        external
        override
    {
        lastEntry = _lastEntry;
        blockNumber = _blockNumber;
    }
}
