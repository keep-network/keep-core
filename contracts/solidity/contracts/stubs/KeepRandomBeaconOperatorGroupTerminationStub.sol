pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperatorGroups.sol";

/**
 * @title KeepRandomBeaconOperatorGroupTerminationStub
 * @dev A simplified Random Beacon group contract to help local development.
 */
contract KeepRandomBeaconOperatorGroupTerminationStub is KeepRandomBeaconOperatorGroups {

    constructor() KeepRandomBeaconOperatorGroups() public {
        groupActiveTime = 5;
        activeGroupsThreshold = 1;
    }

    function addGroup(bytes memory groupPubKey) public {
        groups.push(Group(groupPubKey, block.number));
    }

    function registerNewGroups(uint256 groupsCount) public {
        for (uint i = 1; i <= groupsCount; i++) {
            addGroup(new bytes(i));
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

    function setOperatorContract(address _operatorContract) public {
        operatorContract = _operatorContract;
    }
}
