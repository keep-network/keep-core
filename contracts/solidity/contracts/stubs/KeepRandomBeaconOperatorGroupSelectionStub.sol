pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorGroupSelectionStub is KeepRandomBeaconOperator {
    constructor(
        address _serviceContract,
        address _stakingContract,
        address payable _groupContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract) public {
        relayEntryTimeout = 10;
        ticketInitialSubmissionTimeout = 20;
        ticketReactiveSubmissionTimeout = 65;
        resultPublicationBlockStep = 3;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelectionRelayEntry;
    }

    function callIsTicketValid(
        address staker,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public view returns(bool) {
        return super.isTicketValid(staker, ticketValue, stakerValue, virtualStakerIndex);
    }
}
