pragma solidity ^0.5.4;

import "./KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorStub is KeepRandomBeaconOperator {

    function authorizeServiceContract(address _serviceContract) public {
        serviceContracts.push(_serviceContract);
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.push(Group(groupPublicKey, block.number));
        address[] memory members = orderedParticipants();
        if (members.length > 0) {
            for (uint i = 0; i < groupSize; i++) {
                groupMembers[groupPublicKey].push(members[i]);
            }
        }
    }

    function terminateGroup(uint256 groupIndex) public {
        terminatedGroups.push(groupIndex);
    }

    function getGroupRegistrationBlockHeight(uint256 groupIndex) public view returns(uint256) {
        return groups[groupIndex].registrationBlockHeight;
    }

    function getGroupPublicKey(uint256 groupIndex) public view returns(bytes memory) {
        return groups[groupIndex].groupPubKey;
    }

    function getTicketSubmissionStartBlock() public view returns(uint256) {
        return ticketSubmissionStartBlock;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function setGroupThreshold(uint256 threshold) public {
        groupThreshold = threshold;
    }

    function setActiveGroupsThreshold(uint256 threshold) public {
        activeGroupsThreshold = threshold;
    }

    function setGroupActiveTime(uint256 time) public {
        groupActiveTime = time;
    }

    function setRelayEntryTimeout(uint256 timeout) public {
        relayEntryTimeout = timeout;
    }

    function setMinimumStake(uint256 stake) public {
        minimumStake = stake;
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
