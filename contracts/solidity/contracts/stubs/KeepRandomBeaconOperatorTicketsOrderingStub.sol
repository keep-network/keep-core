pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorTicketsOrderingStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorTicketsOrderingStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
    }

    function addTicket(uint64 newTicketValue) public {
        groupSelection.addTicket(newTicketValue);
    }

    /**
    * @dev Gets submitted group candidate tickets so far.
    */
    function getTickets() public view returns (uint64[] memory) {
        return groupSelection.tickets;
    }

    /**
    * @dev Gets an index of the highest ticket value (tail).
    */
    function getTail() public view returns (uint256) {
        return groupSelection.tail;
    }

    /**
    * @dev Gets an index of a ticket that a higherTicketValueIndex points to.
    * Ex. tickets[23, 5, 65]
    * getPreviousTicketIndex(2) = 0
    */
    function getPreviousTicketIndex(uint256 higherTicketValueIndex) public view returns (uint256) {
        return groupSelection.previousTicketIndex[higherTicketValueIndex];
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
        groupSelection.groupSize = size;
    }
}
