// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import {IRandomBeacon} from "@keep-network/random-beacon/contracts/RandomBeacon.sol";
import {IRandomBeaconConsumer} from "@keep-network/random-beacon/contracts/libraries/Callback.sol";

// TODO: get rid of this contract; use RandomBeacon implementation instead.
contract RandomBeaconStub is IRandomBeacon {
    IRandomBeaconConsumer public callbackContract;

    function requestRelayEntry(IRandomBeaconConsumer _callbackContract)
        external
    {
        callbackContract = _callbackContract;
    }
}
