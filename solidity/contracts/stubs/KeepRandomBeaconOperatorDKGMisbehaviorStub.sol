pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

/**
 * @title KeepRandomBeaconOperatorDKGMisbehaviorStub
 * @dev A simplified Random Beacon operator contract to help local development.
 */
contract KeepRandomBeaconOperatorDKGMisbehaviorStub is KeepRandomBeaconOperator {

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

        groupSize = 5;
        groupSelection.groupSize = 5;
        dkgResultVerification.groupSize = 5;

        groupThreshold = 3;
        dkgResultVerification.signatureThreshold = 3;
    }

    function getGroupSelectionRelayEntry() public view returns (uint256) {
        return groupSelection.seed;
    }

    function getTicketSubmissionStartBlock() public view returns (uint256) {
        return groupSelection.ticketSubmissionStartBlock;
    }

    function timeDKG() public view returns (uint256) {
        return dkgResultVerification.timeDKG;
    }
}
