pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorRewardsStub is KeepRandomBeaconOperator {

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
        groups.groupActiveTime = 5;
        groups.relayEntryTimeout = 10;
    }

    function registerNewGroup(bytes memory groupPublicKey, address[] memory members) public {
        groups.addGroup(groupPublicKey);
        groups.setGroupMembers(groupPublicKey, members, hex"");
        emit DkgResultSubmittedEvent(0, groupPublicKey, "");
    }

    function addGroupMemberReward(bytes memory groupPubKey, uint256 groupMemberReward) public {
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
    }

    function reportUnauthorizedSigning(
        uint256 groupIndex
    ) public {
        // Marks the given group as terminated.
        groups.reportRelayEntryTimeout(groupIndex, groupSize);
        emit UnauthorizedSigningReported(groupIndex);
    }

    function reportRelayEntryTimeout(uint256 groupIndex) public {
        // Marks the given group as terminated.
        groups.reportRelayEntryTimeout(groupIndex, groupSize);
        emit RelayEntryTimeoutReported(groupIndex);
    }

    function isGroupTerminated(uint256 groupIndex) public view returns (bool) {
        return groups.isGroupTerminated(groupIndex);
    }

}
