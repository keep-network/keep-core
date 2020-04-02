pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorGroupSelectionStub is KeepRandomBeaconOperator {
    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groupSelection.ticketSubmissionTimeout = 65;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
        groupSelection.groupSize = size;
    }

    function startGroupSelection(uint256 seed) public {
        groupSelection.start(seed);
    }
}
