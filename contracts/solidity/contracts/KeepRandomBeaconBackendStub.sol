pragma solidity ^0.5.4;

import "./KeepRandomBeaconBackend.sol";

/**
 * @title KeepRandomBeaconBackendStub
 * @dev A simplified Random Beacon backend contract to help local development.
 */
contract KeepRandomBeaconBackendStub is KeepRandomBeaconBackend {

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return 1;
    }

    /**
     * @dev Returns public key of a group from available groups using modulo operator.
     */
    function selectGroup() public view returns(bytes memory) {
        // Compressed public key (G2 point) generated with Go client using secret key 123
        return hex"1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0";
    }

    /**
     * @dev Stub method to authorize frontend contract to help local development.
     */
    function authorizeFrontendContract(address _frontendContract) public {
        frontendContract = _frontendContract;
    }

    /**
     * @dev Stub relay entry to help local development.
     */
    function relayEntry(uint256 _requestID, uint256 _groupSignature, bytes memory _groupPubKey, uint256 _previousEntry, uint256 _seed) public {

        require(BLS.verify(_groupPubKey, abi.encodePacked(_previousEntry, _seed), bytes32(_groupSignature)), "Group signature failed to pass BLS verification.");

        groupSelectionSeed =_groupSignature;
        emit RelayEntryGenerated(_requestID, _groupSignature, _groupPubKey, _previousEntry, _seed);
        FrontendContract(frontendContract).relayEntry(_requestID, _groupSignature, _groupPubKey, _previousEntry, _seed);
    }

}
