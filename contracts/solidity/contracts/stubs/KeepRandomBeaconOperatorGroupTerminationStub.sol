pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorGroupTerminationStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorGroupTerminationStub is KeepRandomBeaconOperator {

    constructor(address _serviceContract, address _stakingContract) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        relayEntryTimeout = 10;
        groupActiveTime = 5;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.push(Group(groupPublicKey, block.number));
    }

    function registerNewGroups(uint256 groupsCount) public {
        for (uint i = 1; i <= groupsCount; i++) {
            registerNewGroup(new bytes(i));
        }
    }

    function terminateGroup(uint256 groupIndex) public {
        terminatedGroups.push(groupIndex);
    }

    function clearGroups() public {
        for (uint i = 0; i < groups.length; i++) {
            delete groupMembers[groups[i].groupPubKey];
        }
        groups.length = 0;
        terminatedGroups.length = 0;
        expiredGroupOffset = 0;
    }

    function setActiveGroupsThreshold(uint256 threshold) public {
        activeGroupsThreshold = threshold;
    }
}
