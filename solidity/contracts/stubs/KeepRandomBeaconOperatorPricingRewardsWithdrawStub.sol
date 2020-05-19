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
        groupSize = 3;
        groups.groupActiveTime = 5;
        groups.relayEntryTimeout = 10;
    }

    function isExpiredGroup(bytes memory groupPubKey) public view returns(bool) {
        uint256 flaggedIndex = groups.groupIndices[groupPubKey];
        require(flaggedIndex > 0, "Group does not exist");
        uint256 i = flaggedIndex ^ (1 << 255);
        return groups.expiredGroupOffset > i;
    }

    function registerNewGroup(bytes memory groupPublicKey, address[] memory members) public {
        groups.addGroup(groupPublicKey);
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function relayEntry() public returns (uint256) {
        (uint256 groupMemberReward,,) = newEntryRewardsBreakdown();
        groups.addGroupMemberReward(
            groups.getGroupPublicKey(currentRequestGroupIndex),
            groupMemberReward
        );
        currentRequestStartBlock = 0;
    }
}
