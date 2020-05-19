pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorPricingDKGStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorPricingDKGStub is KeepRandomBeaconOperator {

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

        // setGroupSize
        groupSize = 20;
        groupSelection.groupSize = 20;
        dkgResultVerification.groupSize = 20;

        // setGroupThreshold
        groupThreshold = 11;
        dkgResultVerification.signatureThreshold = 11;
    }

    function registerNewGroup(bytes memory groupPublicKey) public {
        groups.addGroup(groupPublicKey);
    }

    function setGroupMembers(bytes memory groupPublicKey, address[] memory members) public {
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }

    function getTicketSubmissionStartBlock() public view returns (uint256) {
        return groupSelection.ticketSubmissionStartBlock;
    }

    function isGroupSelectionInProgress() public view returns (bool) {
        return groupSelection.inProgress;
    }

    function getGroupPublicKey(uint256 groupIndex) public view returns (bytes memory) {
        return groups.groups[groupIndex].groupPubKey;
    }

    function setGasPriceCeiling(uint256 _gasPriceCeiling) public {
        gasPriceCeiling = _gasPriceCeiling;
    }

    function timeDKG() public view returns (uint256) {
        return dkgResultVerification.timeDKG;
    }
}
