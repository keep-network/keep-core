pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorTicketsOrderingStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorTicketsOrderingStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract,
        address payable _groupContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract) public {
        groupSize = 10;
    }

    function addTicket(uint256 newTicketValue) public {
        groupSelection.addTicket(newTicketValue);
    }

    /**
    * @dev Gets submitted group candidate tickets so far.
    */
    function getTickets() public view returns (uint256[] memory) {
        return groupSelection.tickets;
    }

    /**
    * @dev Gets an index of a highest ticket value (tail).
    */
    function getTail() public view returns (uint256) {
        return groupSelection.tail;
    }

    /**
    * @dev Gets a highest ticket value from the tickets[] array.
    */
    function getTicketMaxValue() public view returns (uint256) {
        return groupSelection.tickets[groupSelection.tail];
    }

    /**
    * @dev Gets an index of a ticket that a higherTicketValueIndex points to.
    * Ex. tickets[23, 5, 65]
    * getOrderedLinkedTicketIndex(2) = 0
    */
    function getOrderedLinkedTicketIndex(uint256 higherTicketValueIndex) public view returns (uint256) {
        return groupSelection.orderedLinkedTicketIndices[higherTicketValueIndex];
    }
}
