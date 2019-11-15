pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorGroupExpirationStub
 * @dev A simplified Random Beacon group contract to help local development.
 */
contract KeepRandomBeaconOperatorGroupExpirationStub is KeepRandomBeaconOperator {
    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groups.groupActiveTime = 20;
        groups.activeGroupsThreshold = 5;
        groups.relayEntryTimeout = 10;
    }

    using Groups for Groups.Group;

    function addGroup(bytes memory groupPubKey) public {
        groups.groups.push(Groups.Group(groupPubKey, block.number));
    }

    function getGroupRegistrationBlockHeight(uint256 groupIndex) public view returns(uint256) {
        return groups.groups[groupIndex].registrationBlockHeight;
    }

    function getGroupPublicKey(uint256 groupIndex) public view returns(bytes memory) {
        return groups.groups[groupIndex].groupPubKey;
    }

    function selectGroup(uint256 seed) public returns(uint256) {
        return groups.selectGroup(seed);
    }
}
