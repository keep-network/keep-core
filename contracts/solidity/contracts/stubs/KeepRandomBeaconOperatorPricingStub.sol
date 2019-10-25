pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorPricingStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract,
        address payable _groupContract,
        address _ticketsContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract, _groupContract, _ticketsContract) public {}

    function registerNewGroup(bytes memory groupPublicKey) public {
        groupContract.addGroup(groupPublicKey);
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
}
