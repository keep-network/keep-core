pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorStub is KeepRandomBeaconOperator {

    constructor(address _serviceContract, address _stakingContract) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groupThreshold = 15;
        relayEntryTimeout = 10;
        ticketInitialSubmissionTimeout = 20;
        ticketReactiveSubmissionTimeout = 100;
        ticketChallengeTimeout = 60;
        resultPublicationBlockStep = 3;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.push(Group(groupPublicKey, block.number));
    }

    function addGroupMember(bytes memory groupPublicKey, address member) public {
        groupMembers[groupPublicKey].push(member);
    }

    function getGroupPublicKey(uint256 groupIndex) public view returns(bytes memory) {
        return groups[groupIndex].groupPubKey;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

}
