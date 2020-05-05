pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorSlashingStub is KeepRandomBeaconOperator {

    constructor(
        address _serviceContract,
        address _stakingContract,
        address _registryContract
    ) KeepRandomBeaconOperator(
        _serviceContract,
        _stakingContract,
        _registryContract
    ) public {
        relayEntryTimeout = 10;
        groupSelection.ticketSubmissionTimeout = 69;
        resultPublicationBlockStep = 3;
    }

    function registerNewGroup(
        bytes memory groupPublicKey,
        address[] memory members
    ) public {
        groups.addGroup(groupPublicKey);
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }
}
