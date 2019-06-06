pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";


/**
 * @title KeepRandomBeaconStub
 * @dev A simplified random beacon contract to help local development. It
 * returns mocked response straight after calling `requestRelayEntry`
 */
contract KeepRandomBeaconStub is Ownable {

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry);

    uint256 internal _seq;
    uint256 internal _previousEntry;
    mapping (string => bool) internal _initialized;

     struct Request {
        address sender;
        uint256 payment;
        bytes groupPubKey;
        address callbackContract;
        string callbackMethod;
    }

    mapping(uint256 => Request) internal _requests;

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Random Beacon implementation contract.
     */
    function initialize()
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _initialized["KeepRandomBeaconStub"] = true;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconStub"];
    }

    /**
     * @dev Stub method to simulate successful request to generate a new relay entry,
     * which will include a random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @param callbackContract Callback contract address.
     * @param callbackMethod Callback contract method signature.
     * @return An uint256 representing uniquely generated relay request ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 seed, address callbackContract, string memory callbackMethod) public payable returns (uint256 requestID) {
        requestID = _seq++;

        // Return mocked data instead of interacting with relay.
        uint256 groupSignature = uint256(keccak256(abi.encodePacked(_previousEntry, block.timestamp, seed)));
        bytes memory groupPubKey = abi.encodePacked(keccak256(abi.encodePacked(block.timestamp, uint(1))));
        
        emit RelayEntryRequested(requestID, msg.value, _previousEntry, seed, groupPubKey);
        emit RelayEntryGenerated(requestID, groupSignature, groupPubKey, _previousEntry);
        _requests[requestID] = Request(msg.sender, msg.value, groupPubKey, callbackContract, callbackMethod);

        _previousEntry = groupSignature;
        return requestID;
    }

    /**
     * @dev Stub method to perform a contract callback.
     * @param requestID The request that started this generation - to tie the results back to the request.
     * @param groupSignature The generated random number.
     * @param groupPubKey Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 requestID, uint256 groupSignature, bytes memory groupPubKey, uint256 previousEntry, uint256 seed) public {

        address callbackContract = _requests[requestID].callbackContract;
        if (callbackContract != address(0)) {
            callbackContract.call(abi.encodeWithSignature(_requests[requestID].callbackMethod, groupSignature));
        }
    }
}
