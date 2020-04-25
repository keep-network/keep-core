pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorRewardsStub is KeepRandomBeaconOperator {

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

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function setGroupMembers(bytes memory groupPublicKey, address[] memory members) public {
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function addGroupMemberReward(bytes memory groupPubKey, uint256 groupMemberReward) public {
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);
    }

    function emitRewardsWithdrawnEvent(address operator, uint256 groupIndex) public {
        emit GroupMemberRewardsWithdrawn(stakingContract.beneficiaryOf(operator), operator, 1000 wei, groupIndex);
    }

    function reportUnauthorizedSigning(
        uint256 groupIndex,
        bytes memory signedMsgSender
    ) public {
        uint256 minimumStake = stakingContract.minimumStake();
        stakingContract.seize(
            minimumStake,
            100,
            msg.sender,
            groups.getGroupMembers(0)
        );
        emit UnauthorizedSigningReported(groupIndex);
    }

    function reportRelayEntryTimeout() public {
        uint256 minimumStake = stakingContract.minimumStake();
        stakingContract.seize(
            minimumStake,
            100,
            msg.sender,
            groups.getGroupMembers(0)
        );
        emit RelayEntryTimeoutReported(0);
    }

}
