pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract,
        address payable _groupContract,
        address _ticketsContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract, _ticketsContract) public {
        groupThreshold = 15;
        relayEntryTimeout = 10;
        resultPublicationBlockStep = 3;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groupContract.addGroup(groupPublicKey);
    }

    function addGroupMember(bytes memory groupPublicKey, address member) public {
        groupContract.addGroupMember(groupPublicKey, member);
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelectionRelayEntry;
    }

    function getTicketSubmissionStartBlock() public view returns (uint256) {
        return ticketContract.ticketSubmissionStartBlock();
    }
}
