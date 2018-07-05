pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";


/**
 * @title KeepRandomBeaconStub
 * @dev A simplified random beacon contract to help local development. It
 * returns mocked response straight after calling `requestRelayEntry`
 */
contract KeepRandomBeaconStub is Ownable, EternalStorage {

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber);
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 activationBlockHeight);

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert();
    }

    /**
     * @dev Initialize Keep Random Beacon implementation contract.
     */
    function initialize()
        public
        onlyOwner
    {
        require(!initialized());
        boolStorage[keccak256("KeepRandomBeaconStub")] = true;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return boolStorage[keccak256("KeepRandomBeaconStub")];
    }

    /**
     * @dev Stub method to simulate successful request to generate a new relay entry,
     * which will include a random number (by signing the previous entry's random number).
     * @param _blockReward The value in KEEP for generating the signature.
     * @param _seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 _blockReward, uint256 _seed) public payable returns (uint256 requestID) {
        requestID = uintStorage[keccak256("seq")]++;
        emit RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);

        // Return mocked data instead of interacting with relay.
        uint256 _previousEntry = uintStorage[keccak256("previousEntry")];
        uint256 _groupSignature = uint256(keccak256(_previousEntry, block.timestamp, _seed));
        uint256 _groupID = uint256(keccak256(block.timestamp, 1));
        emit RelayEntryGenerated(requestID, _groupSignature, _groupID, _previousEntry, block.number);

        uintStorage[keccak256("previousEntry")] = _groupSignature;
        return requestID;
    }
}
