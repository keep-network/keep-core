pragma solidity 0.5.17;

import "../../contracts/KeepRandomBeaconOperator.sol";


contract KeepRandomBeaconOperatorPricingStub is KeepRandomBeaconOperator {
    constructor(
        address _serviceContract,
        address _stakingContract,
        address _registryContract
    )
        public
        KeepRandomBeaconOperator(
            _serviceContract,
            _stakingContract,
            _registryContract
        )
    {}

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function setDkgGasEstimate(uint256 gasEstimate) public {
        dkgGasEstimate = gasEstimate;
    }

    function setEntryVerificationGasEstimate(uint256 gasEstimate) public {
        entryVerificationGasEstimate = gasEstimate;
    }

    function setGroupMemberBaseReward(uint256 reward) public {
        groupMemberBaseReward = reward;
    }

    function setGroupSize(uint256 size) public {
        groupSize = size;
    }

    function setGroupSelectionGasEstimate(uint256 gas) public {
        groupSelectionGasEstimate = gas;
    }

    function setGasPriceCeiling(uint256 _gasPriceCeiling) public {
        gasPriceCeiling = _gasPriceCeiling;
    }

    function getNewEntryRewardsBreakdown()
        public
        view
        returns (
            uint256 groupMemberReward,
            uint256 submitterReward,
            uint256 subsidy
        )
    {
        return super.newEntryRewardsBreakdown();
    }

    function delayFactor() public view returns (uint256) {
        return
            DelayFactor.calculate(currentRequestStartBlock, relayEntryTimeout);
    }
}
