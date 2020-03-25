pragma solidity ^0.5.4;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorMisbihaveEventsStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract
    ) KeepRandomBeaconOperator(_serviceContract, _stakingContract) public {
    }

    function reportUnauthorizedSigning(
        uint256 groupIndex,
        bytes memory signedMsgSender
    ) public {
        uint256 minimumStake = stakingContract.minimumStake();
        address[] memory operators;
        stakingContract.seize(minimumStake, 100, msg.sender, operators);
        emit UnauthorizedSigningReported(groupIndex);
    }

    function reportRelayEntryTimeout() public {
        uint256 minimumStake = stakingContract.minimumStake();
        address[] memory operators;
        stakingContract.seize(minimumStake, 100, msg.sender, operators);
        emit RelayEntryTimeoutReported(signingRequest.groupIndex);
    }

}