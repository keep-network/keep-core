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
        address payable _groupContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract) public {
        groupThreshold = 15;
        relayEntryTimeout = 10;
        groupSelection.ticketSubmissionTimeout = 65;
        resultPublicationBlockStep = 3;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groupContract.addGroup(groupPublicKey);
    }

    function addGroupMember(bytes memory groupPublicKey, address member) public {
        groupContract.addGroupMember(groupPublicKey, member);
    }

    function setGroupSize(uint8 size) public {
        groupSize = size;
        groupSelection.groupSize = size;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }

    function getTicketSubmissionStartBlock() public view returns (uint256) {
        return groupSelection.ticketSubmissionStartBlock;
    }

    function isGroupSelectionInProgress() public view returns (bool) {
        return groupSelection.inProgress;
    }

    function getRelayEntryTimeout() public view returns (uint256) {
        return relayEntryTimeout;
    }
}
