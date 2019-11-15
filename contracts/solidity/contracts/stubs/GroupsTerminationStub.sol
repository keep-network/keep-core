pragma solidity ^0.5.4;
import "../libraries/Groups.sol";

contract GroupsTerminationStub {
    using Groups for Groups.Group;
    using Groups for Groups.Storage;
    Groups.Storage groups;

    constructor() public {
        groups.groupActiveTime = 5;
        groups.activeGroupsThreshold = 1;
    }

    function addGroup(bytes memory groupPubKey) public {
        groups.groups.push(Groups.Group(groupPubKey, block.number));
    }

    function registerNewGroups(uint256 groupsCount) public {
        for (uint i = 1; i <= groupsCount; i++) {
            groups.addGroup(new bytes(i));
        }
    }

    function terminateGroup(uint256 groupIndex) public {
        groups.terminatedGroups.push(groupIndex);
    }

    function clearGroups() public {
        for (uint i = 0; i < groups.groups.length; i++) {
            delete groups.groupMembers[groups.groups[i].groupPubKey];
        }
        groups.groups.length = 0;
        groups.terminatedGroups.length = 0;
        groups.expiredGroupOffset = 0;
    }

    function setActiveGroupsThreshold(uint256 threshold) public {
        groups.activeGroupsThreshold = threshold;
    }

    function selectGroup(uint256 seed) public returns(uint256) {
        return groups.selectGroup(seed);
    }
}
