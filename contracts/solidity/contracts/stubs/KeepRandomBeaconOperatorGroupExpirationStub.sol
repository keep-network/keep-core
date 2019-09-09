pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperatorGroups.sol";

/**
 * @title KeepRandomBeaconOperatorGroupExpirationStub
 * @dev A simplified Random Beacon group contract to help local development.
 */
contract KeepRandomBeaconOperatorGroupExpirationStub is KeepRandomBeaconOperatorGroups {

    constructor() KeepRandomBeaconOperatorGroups() public {
        groupActiveTime = 300;
        activeGroupsThreshold = 5;
        relayEntryTimeout = 24;
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

    function setOperatorContract(address _operatorContract) public {
        operatorContract = _operatorContract;
    }
}
