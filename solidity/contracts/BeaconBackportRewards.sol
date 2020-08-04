pragma solidity ^0.5.17;

import "./Rewards.sol";
import "./KeepRandomBeaconOperator.sol";
import "./TokenStaking.sol";

contract BeaconBackportRewards is Rewards {
    uint256[] lastGroupOfInterval;
    mapping(uint256 => bool) excludedGroups;
    uint256 excludedGroupCount;
    KeepRandomBeaconOperator operatorContract;
    TokenStaking tokenStaking;
    address excessRecipient;
    uint256 constant _minimumKeepsPerInterval = 1; // TODO define

    constructor (
        // The term length and first interval start are arbitrary,
        // and should be set so that `firstIntervalStart + termLength`
        // is sometime in the future.
        // Up to that point (when interval 0 finishes)
        // the contract can be funded freely,
        // with tokens being allocated correctly.
        // After that point, interval 0 may get allocated
        // so any further funding will end up on the remaining intervals.
        // If unallocated tokens are remaining after all intervals are allocated,
        // a withdrawal function is provided.
        address _token,
        uint256 _firstIntervalStart,
        // Interval weights. Does not define the number of intervals supported.
        uint256[] memory _intervalWeights,
        address _operatorContract,
        address _stakingContract,
        // The indices of the last group eligible for rewards in each interval.
        // Inclusive, defines the number of intervals supported by the contract.
        uint256[] memory _lastGroupOfInterval,
        // The indices of any groups below `lastEligibleGroup`
        // that should be excluded from the rewards.
        uint256[] memory _excludedGroups,
        // The address that receives any left-over tokens
        // from late funding or excluded groups.
        address _excessRecipient
    ) public Rewards(
        _token,
        _firstIntervalStart,
        _intervalWeights
    ) {
        operatorContract = KeepRandomBeaconOperator(_operatorContract);
        tokenStaking = TokenStaking(_stakingContract);
        lastGroupOfInterval = _lastGroupOfInterval;
        for (uint256 i = 0; i < _excludedGroups.length; i++) {
            excludedGroups[_excludedGroups[i]] = true;
        }
        excessRecipient = _excessRecipient;
    }

    function minimumKeepsPerInterval() public view returns (uint256) {
        return _minimumKeepsPerInterval;
    }

    function lastEligibleGroup() public view returns (uint256) {
        return lastGroupOfInterval[lastInterval()];
    }

    function lastInterval() public view returns (uint256) {
        return lastGroupOfInterval.length.sub(1);
    }

    function withdrawExcess() mustBeFinished(lastInterval()) public {
        if (lastInterval() >= intervalAllocations.length) {
            allocateRewards(lastInterval());
        }
        token.safeTransfer(excessRecipient, unallocatedRewards);
        unallocatedRewards = 0;
    }

    function _assignedInterval(bytes32 groupIndexBytes) internal view returns (uint256) {
        for (uint256 interval = 0; interval < lastGroupOfInterval.length; interval++) {
            if (lastGroupOfInterval[interval] >= uint256(groupIndexBytes)) {
                return interval;
            }
        }
    }

    function _getKeepCount() internal view returns (uint256) {
        return lastEligibleGroup().add(1);
    }

    function _getKeepAtIndex(uint256 i) internal view returns (bytes32) {
        return bytes32(i);
    }

    function _getCreationTime(bytes32 groupIndexBytes) internal view returns (uint256) {
        // Assign each group to the starting timestamp of its interval
        return startOf(_assignedInterval(groupIndexBytes));
    }

    function _isClosed(bytes32 groupIndexBytes) internal view returns (bool) {
        if (_isTerminated(groupIndexBytes)) { return false; }
        bytes memory groupPubkey = operatorContract.getGroupPublicKey(
            uint256(groupIndexBytes)
        );
        return operatorContract.isStaleGroup(groupPubkey);
    }

    function _isTerminated(bytes32 groupIndexBytes) internal view returns (bool) {
        return excludedGroups[uint256(groupIndexBytes)];
    }

    // A group is recognized if its index is at most `lastEligibleGroup`
    // and it isn't listed as excluded.
    function _recognizedByFactory(bytes32 groupIndexBytes) internal view returns (bool) {
        return lastEligibleGroup() >= uint256(groupIndexBytes);
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
