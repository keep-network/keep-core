pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";

contract KeepRandomBeaconOperatorPricingRewardsWithdrawStub is KeepRandomBeaconOperator {

    using BytesLib for bytes;

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
        groups.groupActiveTime = 5;
        groups.relayEntryTimeout = 10;
    }

    function isExpiredGroup(bytes memory groupPubKey) public view returns(bool) {
        for (uint i = 0; i < groups.groups.length; i++) {
            if (groups.groups[i].groupPubKey.equalStorage(groupPubKey)) {
                bool isExpired = groups.expiredGroupOffset > i;
                return isExpired;
            }
        }

        revert("Group does not exist");
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function setGroupMembers(bytes memory groupPublicKey, address[] memory members) public {
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function relayEntry() public returns (uint256) {
        bytes memory groupPubKey = groups.getGroupPublicKey(signingRequest.groupIndex);
        (uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) = newEntryRewardsBreakdown();
        submitterReward; // silence local var
        subsidy; // silence local var
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
        currentEntryStartBlock = 0;
    }
}
