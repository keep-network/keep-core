pragma solidity ^0.5.17;

import "./Rewards.sol";
import "./KeepRandomBeaconOperator.sol";
import "./TokenStaking.sol";

contract BeaconBackportRewards is Rewards {
    uint256 lastEligibleGroup;
    mapping(uint256 => bool) excludedGroups;
    uint256 excludedGroupCount;
    KeepRandomBeaconOperator operatorContract;
    TokenStaking tokenStaking;

    constructor (
        // Rewards can be allocated at `firstIntervalStart + termLength`.
        // Exact values are arbitrary.
        // After allocation, the contract should no longer be funded.
        uint256 _termLength,
        address _token,
        uint256 _firstIntervalStart,
        address _operatorContract,
        address _stakingContract,
        // The index of the last group eligible for rewards. Inclusive.
        uint256 _lastEligibleGroup,
        // The indices of any groups below `lastEligibleGroup`
        // that should be excluded from the rewards.
        uint256[] memory _excludedGroups
    ) public Rewards(
        _termLength,
        _token,
        1, // _minimumKeepsPerInterval,
        _firstIntervalStart,
        [] // _intervalWeights
    ) {
        operatorContract = KeepRandomBeaconOperator(_operatorContract);
        tokenStaking = TokenStaking(_stakingContract);
        lastEligibleGroup = _lastEligibleGroup;
        excludedGroupCount = _excludedGroups.length;
        for (uint256 i = 0; i < _excludedGroups.length; i++) {
            excludedGroups[_excludedGroups[i]] = true;
        }
    }

    function _getKeepCount() internal view returns (uint256) {
        return lastEligibleGroup.add(1).sub(excludedGroupCount);
    }

    function _getKeepAtIndex(uint256 i) internal view returns (bytes32) {
        return bytes32(i);
    }

    // All eligible groups are in interval 0
    function _getCreationTime(bytes32 groupIndexBytes) internal view returns (uint256) {
        return 0;
    }

    function _isClosed(bytes32 groupIndexBytes) internal view returns (bool) {
        bytes memory groupPubkey = operatorContract.getGroupPublicKey(
            uint256(groupIndexBytes)
        );
        return operatorContract.isStaleGroup(groupPubkey);
    }

    function _isTerminated(bytes32 groupIndexBytes) internal view returns (bool) {
        return false;
    }

    // A group is recognized if its index is at most `lastEligibleGroup`
    // and it isn't listed as excluded.
    function _recognizedByFactory(bytes32 groupIndexBytes) internal view returns (bool) {
        uint256 groupIndex = uint256(groupIndexBytes);
        return (lastEligibleGroup >= groupIndex) && !excludedGroups[groupIndex];
    }

    function _distributeReward(bytes32 groupIndexBytes, uint256 _value) internal {
        bytes memory groupPubkey = operatorContract.getGroupPublicKey(
            uint256(groupIndexBytes)
        );
        address[] memory members = operatorContract.getGroupMembers(groupPubkey);

        uint256 memberCount = members.length;
        uint256 dividend = _value.div(memberCount);

        // Only pay other members if dividend is nonzero.
        if(dividend > 0) {
            for (uint256 i = 0; i < memberCount - 1; i++) {
                token.safeTransfer(
                    tokenStaking.beneficiaryOf(members[i]),
                    dividend
                );
            }
        }

        // Transfer of dividend for the last member. Remainder might be equal to
        // zero in case of even distribution or some small number.
        uint256 remainder = _value.mod(memberCount);
        token.safeTransfer(
            tokenStaking.beneficiaryOf(members[memberCount - 1]),
            dividend.add(remainder)
        );
    }
}
