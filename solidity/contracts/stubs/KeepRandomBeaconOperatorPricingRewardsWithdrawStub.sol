pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";
import "../utils/BytesLib.sol";

contract KeepRandomBeaconOperatorPricingRewardsWithdrawStub is KeepRandomBeaconOperator {

    using BytesLib for bytes;

    constructor(
        address _serviceContract,
        address _stakingContract,
        address _registryContract,
        address _gasPriceOracle
    ) KeepRandomBeaconOperator(
        _serviceContract,
        _stakingContract,
        _registryContract,
        _gasPriceOracle
    ) public {
        groupSize = 3;
        groups.groupActiveTime = 3;
        groups.relayEntryTimeout = 10;
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
