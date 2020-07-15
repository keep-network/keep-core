pragma solidity ^0.5.17;

import "../Rewards.sol";

contract RewardsStub is Rewards {
    uint256[] creationTimes;
    uint256 closedTime;
    mapping(uint256 => bool) terminated;

    constructor (
        uint256 _termLength,
        address _token,
        uint256 _minimumKeepsPerInterval,
        uint256 _firstIntervalStart,
        uint256[] memory _intervalWeights,
        uint256[] memory _creationTimes
    ) public Rewards(
        _termLength,
        _token,
        _minimumKeepsPerInterval,
        _firstIntervalStart,
        _intervalWeights
    ) {
        creationTimes = _creationTimes;
    }

    function receiveReward(uint256 i) public {
        receiveReward(bytes32(i));
    }

    function reportTermination(uint256 i) public {
        reportTermination(bytes32(i));
    }

    function eligibleForReward(uint256 i) public view returns (bool) {
        return eligibleForReward(bytes32(i));
    }

    function getTotalRewards() public view returns (uint256) {
        return totalRewards;
    }

    function getUnallocatedRewards() public view returns (uint256) {
        return unallocatedRewards;
    }

    function getPaidOutRewards() public view returns (uint256) {
        return paidOutRewards;
    }

    function findEndpoint(uint256 intervalEndpoint) public view returns (uint256) {
        return _findEndpoint(intervalEndpoint);
    }

    function getEndpoint(uint256 interval) public returns (uint256) {
        return _getEndpoint(interval);
    }

    function keepsInInterval(uint256 interval) public returns (uint256) {
        return _keepsInInterval(interval);
    }

    function getIntervalWeight(uint256 interval) public view returns (uint256) {
        return _getIntervalWeight(interval);
    }

    function getIntervalCount() public view returns (uint256) {
        return _getIntervalCount();
    }

    function baseAllocation(uint256 interval) public view returns (uint256) {
        return _baseAllocation(interval);
    }

    function adjustedAllocation(uint256 interval) public returns (uint256) {
        return _adjustedAllocation(interval);
    }

    function rewardPerKeep(uint256 interval) public returns (uint256) {
        uint256 __adjustedAllocation = _adjustedAllocation(interval);
        if (__adjustedAllocation == 0) {
            return 0;
        }
        uint256 keepCount = _keepsInInterval(interval);
        // Adjusted allocation would be zero if keep count was zero
        assert(keepCount > 0);
        return __adjustedAllocation.div(keepCount);
    }

    function allocateRewards(uint256 interval) public {
        _allocateRewards(interval);
    }

    function getAllocatedRewards(uint256 interval) public view returns (uint256) {
        return _getAllocatedRewards(interval);
    }

    function isAllocated(uint256 interval) public view returns (bool) {
        return _isAllocated(interval);
    }

    function terminate(uint256 i) public {
        terminated[i] = true;
    }

    function setCloseTime(uint256 i) public {
        closedTime = i;
    }

    function _getKeepCount() internal view returns (uint256) {
        return creationTimes.length;
    }

    function _getKeepAtIndex(uint256 i) internal view returns (bytes32) {
        return bytes32(i);
    }

    function _getCreationTime(bytes32 groupIndexBytes) internal view returns (uint256) {
        return creationTimes[uint256(groupIndexBytes)];
    }

    function _isClosed(bytes32 groupIndexBytes) internal view returns (bool) {
        return _getCreationTime(groupIndexBytes) <= closedTime;
    }

    function _isTerminated(bytes32 groupIndexBytes) internal view returns (bool) {
        return terminated[uint256(groupIndexBytes)];
    }

    function _recognizedByFactory(bytes32 groupIndexBytes) internal view returns (bool) {
        return _getKeepCount() > uint256(groupIndexBytes);
    }

    function _distributeReward(bytes32 groupIndexBytes, uint256 _value) internal {
        token.safeTransfer(
            msg.sender,
            _value
        );
    }
}
