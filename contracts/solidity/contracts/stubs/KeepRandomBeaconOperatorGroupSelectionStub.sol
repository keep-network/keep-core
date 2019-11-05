pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorGroupSelectionStub is KeepRandomBeaconOperator {
    constructor(
        address _serviceContract,
        address _stakingContract,
        address payable _groupContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract) public {
        groupSize = 3;
        groupSelection.ticketSubmissionTimeout = 65;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }
}
