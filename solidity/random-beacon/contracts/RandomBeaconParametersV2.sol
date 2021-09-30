pragma solidity ^0.8.6;

import "@openzeppelin/contracts/access/Ownable.sol";

contract RandomBeaconParametersV2 is Ownable {
    mapping(address => bytes) public initializationData;

    uint256 constant DELAY = 24 hours;

    string private constant PREFIX = "network.keep.random-beacon.params.";
    string private constant PREFIX_NEW =
        "network.keep.random-beacon.params.update.new.";
    string private constant PREFIX_UPDATE_TIMESTAMP =
        "network.keep.random-beacon.params.update.timestamp.";

    event UpdateStarted(string indexed parameter, uint256 newValue);
    event UpdateCompleted(string indexed parameter, uint256 newValue);

    error GovernanceUpdateNotInitiated();
    error GovernanceDelayNotElapsed();

    modifier onlyAfterGovernanceDelay(string memory parameter, uint256 delay) {
        uint256 slotTimestamp = calcSlot(PREFIX_UPDATE_TIMESTAMP, parameter);

        uint256 initiated;
        assembly {
            initiated := sload(slotTimestamp)
        }

        if (initiated == 0) revert GovernanceUpdateNotInitiated();

        if (block.timestamp - initiated < delay)
            revert GovernanceDelayNotElapsed();

        _;
    }

    function getParameter(string memory parameter)
        public
        view
        returns (uint256 value)
    {
        uint256 slot = calcSlot(PREFIX, parameter);

        assembly {
            value := sload(slot)
        }
    }

    function getParameterNewValue(string memory parameter)
        public
        view
        returns (uint256 value)
    {
        uint256 slot = calcSlot(PREFIX_NEW, parameter);

        assembly {
            value := sload(slot)
        }
    }

    function beginUpdate(string memory parameter, uint256 newValue)
        public
        onlyOwner
    {
        uint256 slotNew = calcSlot(PREFIX_NEW, parameter);
        uint256 slotTimestamp = calcSlot(PREFIX_UPDATE_TIMESTAMP, parameter);

        uint256 newValueCheck;
        assembly {
            sstore(slotNew, newValue)
            sstore(slotTimestamp, timestamp())

            newValueCheck := sload(slotNew)
        }

        emit UpdateStarted(parameter, newValue);
    }

    // TODO: onlyOwner is not needed
    function finalizeUpdate(string memory parameter)
        public
        onlyOwner
        onlyAfterGovernanceDelay(parameter, DELAY)
    {
        uint256 slot = calcSlot(PREFIX, parameter);
        uint256 slotNewValue = calcSlot(PREFIX_NEW, parameter);

        uint256 newValue;
        assembly {
            newValue := sload(slotNewValue)
            sstore(slot, newValue)
        }

        emit UpdateCompleted(parameter, newValue);

        // TODO: delete stored value for update
    }

    function calcSlot(string memory prefix, string memory parameter)
        private
        pure
        returns (uint256)
    {
        return uint256(keccak256(abi.encodePacked(prefix, parameter))) - 1;
    }
}
