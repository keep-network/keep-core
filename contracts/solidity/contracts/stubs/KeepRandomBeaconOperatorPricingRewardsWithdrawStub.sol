pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorPricingRewardsWithdrawStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groups.groupActiveTime = 5;
        groups.activeGroupsThreshold = 1;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function addGroupMember(bytes memory groupPublicKey, address member) public {
        groups.addGroupMember(groupPublicKey, member);
    }

    function relayEntry() public returns (uint256) {
        entryInProgress = false;
        bytes memory groupPubKey = groups.getGroupPublicKey(signingRequest.groupIndex);
        (uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) = newEntryRewardsBreakdown();
        submitterReward; // silence local var
        subsidy; // silence local var
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
    }
}
