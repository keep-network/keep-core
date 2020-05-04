pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";
import "../utils/BytesLib.sol";

contract KeepRandomBeaconOperatorPricingRewardsWithdrawStub is KeepRandomBeaconOperator {

    using BytesLib for bytes;

    constructor(
        address _serviceContract,
        address _stakingContract,
        address _registryContract
    ) KeepRandomBeaconOperator(
        _serviceContract,
        _stakingContract,
        _registryContract
    ) public {
        groups.groupActiveTime = 5;
        groups.relayEntryTimeout = 10;
    }

    function isExpiredGroup(bytes memory groupPubKey) public view returns(bool) {
        uint256 flaggedIndex = groups.groupIndices[groupPubKey];
        require(flaggedIndex > 0, "Group does not exist");
        uint256 i = flaggedIndex ^ (1 << 255);
        return groups.expiredGroupOffset > i;
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
        bytes memory groupPubKey = groups.getGroupPublicKey(currentRequestGroupIndex);
        (uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) = newEntryRewardsBreakdown();
        submitterReward; // silence local var
        subsidy; // silence local var
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
        currentRequestStartBlock = 0;
    }
}
