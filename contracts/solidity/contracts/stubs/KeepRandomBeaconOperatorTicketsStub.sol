pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperatorTickets.sol";

/**
 * @title KeepRandomBeaconOperatorTicketsStub
 * @dev A simplified Random Beacon operator tickets contract to help local development.
 */
contract KeepRandomBeaconOperatorTicketsStub is KeepRandomBeaconOperatorTickets {

    constructor() KeepRandomBeaconOperatorTickets() public {
        ticketInitialSubmissionTimeout = 20;
        ticketReactiveSubmissionTimeout = 65;
    }

}
