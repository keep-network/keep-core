pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorRewardsStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groups.groupActiveTime = 5;
        groups.activeGroupsThreshold = 1;
        groups.relayEntryTimeout = 10;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function addGroupMember(bytes memory groupPublicKey, address member) public {
        groups.addGroupMember(groupPublicKey, member);
    }

    function addGroupMemberReward(bytes memory groupPubKey, uint256 groupMemberReward) public {
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
    }

}