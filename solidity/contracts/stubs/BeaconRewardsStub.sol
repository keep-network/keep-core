pragma solidity ^0.5.17;

import "../BeaconRewards.sol";

contract BeaconRewardsStub is BeaconRewards {
    constructor (
        uint256 _termLength,
        address _token,
        uint256 _minimumKeepsPerInterval,
        uint256 _firstIntervalStart,
        uint256[] memory _intervalWeights,
        address _operatorContract,
        address _stakingContract
    ) public BeaconRewards(
        _termLength,
        _token,
        _minimumKeepsPerInterval,
        _firstIntervalStart,
        _intervalWeights,
        _operatorContract,
        _stakingContract
    ) {
    }

    function getKeepCount() public view returns (uint256) {
        return _getKeepCount();
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

    function isTerminated(uint256 i) public view returns (bool) {
        return eligibleButTerminated(bytes32(i));
    }

    function recognizedByFactory(uint256 i) public view returns (bool) {
        return _recognizedByFactory(bytes32(i));
    }

    function getTotalRewards() public view returns (uint256) {
        return totalRewards;
    }

    function getUnallocatedRewards() public view returns (uint256) {
        return unallocatedRewards;
    }
}
