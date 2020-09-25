pragma solidity ^0.5.17;

import "../BeaconBackportRewards.sol";

contract BeaconBackportRewardsStub is BeaconBackportRewards {
    constructor (
        address _token,
        uint256 _firstIntervalStart,
        uint256[] memory _intervalWeights,
        address _operatorContract,
        address _stakingContract,
        uint256[] memory _lastGroupOfInterval,
        uint256[] memory _excludedGroups,
        address _excessRecipient
    ) public BeaconBackportRewards(
        _token,
        _firstIntervalStart,
        _intervalWeights,
        _operatorContract,
        _stakingContract,
        _lastGroupOfInterval,
        _excludedGroups,
        _excessRecipient
    ) {}

    function receiveReward(uint256 i) public {
        receiveReward(bytes32(i));
    }

    function reportTermination(uint256 i) public {
        reportTermination(bytes32(i));
    }

    function eligibleForReward(uint256 i) public view returns (bool) {
        return eligibleForReward(bytes32(i));
    }

    function eligibleButTerminatedWithUint(uint256 i) public view returns (bool) {
        return eligibleButTerminated(bytes32(i));
    }

    function rewardClaimedWithUint(uint256 i) public view returns (bool) {
        return rewardClaimed(bytes32(i));
    }

    function findEndpoint(uint256 i) public view returns (uint256) {
        return _findEndpoint(i);
    }

    function getKeepCount() public view returns (uint256) {
        return _getKeepCount();
    }

    function recognizedByFactory(uint256 i) public view returns (bool) {
        return _recognizedByFactory(bytes32(i));
    }

    function isExcluded(uint256 i) public view returns (bool) {
        return excludedGroups[i];
    }
}
