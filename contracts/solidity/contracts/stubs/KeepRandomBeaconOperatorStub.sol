pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorStub is KeepRandomBeaconOperator {

    constructor(address _serviceContract, address _stakingContract) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {}

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.push(Group(groupPublicKey, block.number));
        address[] memory members = orderedParticipants();
        if (members.length > 0) {
            for (uint i = 0; i < groupSize; i++) {
                groupMembers[groupPublicKey].push(members[i]);
            }
        }
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

    function setRelayEntryTimeout(uint256 timeout) public {
        relayEntryTimeout = timeout;
    }

    function setMinimumStake(uint256 stake) public {
        minimumStake = stake;
    }

    function setGroupThreshold(uint256 threshold) public {
        groupThreshold = threshold;
    }

    function setTicketInitialSubmissionTimeout(uint256 timeout) public {
        ticketInitialSubmissionTimeout = timeout;
    }

    function setTicketReactiveSubmissionTimeout(uint256 timeout) public {
        ticketReactiveSubmissionTimeout = timeout;
    }

    function setTicketChallengeTimeout(uint256 timeout) public {
        ticketChallengeTimeout = timeout;
    }

    function setResultPublicationBlockStep(uint256 blockStep) public {
        resultPublicationBlockStep = blockStep;
    }
}
