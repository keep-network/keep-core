pragma solidity ^0.5.4;
import "../libraries/operator/Groups.sol";

contract GroupsTerminationStub {
    using Groups for Groups.Storage;
    Groups.Storage groups;

    constructor() public {
        groups.groupActiveTime = 5;
    }

    function addGroup(bytes memory groupPubKey) public {
        groups.addGroup(groupPubKey);
    }

    function registerNewGroups(uint256 groupsCount) public {
        for (uint i = 1; i <= groupsCount; i++) {
            groups.addGroup(new bytes(i));
        }
    }

    function terminateGroup(uint256 groupIndex) public {
        groups.terminatedGroups.push(groupIndex);
    }

    function selectGroup(uint256 seed) public returns(uint256) {
        return groups.selectGroup(seed);
    }
}
