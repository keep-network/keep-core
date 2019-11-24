pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorTicketsOrderingStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorTicketsPackingStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groupSelection.ticketSubmissionTimeout = 300;
    }

    function addTicket(uint64 newTicketValue) public {
        groupSelection.addTicket(newTicketValue);
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
        groupSelection.groupSize = size;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }

}
