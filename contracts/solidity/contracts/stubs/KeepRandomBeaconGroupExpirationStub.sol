pragma solidity ^0.5.4;

import "../KeepRandomBeaconGroups.sol";

/**
 * @title KeepRandomBeaconGroupExpirationStub
 * @dev A simplified Random Beacon group contract to help local development.s
 */
contract KeepRandomBeaconGroupExpirationStub is KeepRandomBeaconGroups {

    constructor() KeepRandomBeaconGroups() public {
        groupActiveTime = 300;
        activeGroupsThreshold = 5;
    }

    function addGroup(bytes memory groupPubKey) public {
        groups.push(Group(groupPubKey, block.number));
    }

    function getGroupRegistrationBlockHeight(uint256 groupIndex) public view returns(uint256) {
        return groups[groupIndex].registrationBlockHeight;
    }

    function getGroupPublicKey(uint256 groupIndex) public view returns(bytes memory) {
        return groups[groupIndex].groupPubKey;
    }

}
