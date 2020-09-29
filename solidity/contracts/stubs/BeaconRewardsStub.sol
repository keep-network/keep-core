pragma solidity ^0.5.17;

import "../BeaconRewards.sol";

contract BeaconRewardsStub is BeaconRewards {
    constructor (
        address _token,
        uint256 _firstIntervalStart,
        address _operatorContract,
        address _stakingContract
    ) public BeaconRewards(
        _token,
        _firstIntervalStart,
        _operatorContract,
        _stakingContract
    ) {
    }

    function setIntervalWeights(uint256[] memory _intervalWeights) public {
        intervalWeights = _intervalWeights;
    }

    function setTermLength(uint256 _termLength) public {
        termLength = _termLength;
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

    function findEndpoint(uint256 intervalEndpoint) public view returns (uint256) {
        return _findEndpoint(intervalEndpoint);
    }
}
