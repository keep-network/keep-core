pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./DelayedWithdrawal.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./cryptography/BLS.sol";


interface GroupContract {
    function runGroupSelection(uint256 randomBeaconValue) external;
    function numberOfGroups() external view returns(uint256);
    function selectGroup(uint256 previousEntry) external view returns(bytes memory);
}


/**
 * @title KeepRandomBeaconFrontendImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconFrontendImplV1 is Ownable, DelayedWithdrawal {

    using BytesLib for bytes;

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry, uint256 seed);

    uint256 internal _requestCounter;
    uint256 internal _minPayment;
    address internal _groupContract;
    uint256 internal _previousEntry;

    mapping (string => bool) internal _initialized;

    struct Request {
        address sender;
        uint256 payment;
        bytes groupPubKey;
    }

    mapping(uint256 => Request) internal _requests;

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Random Beacon implementaion contract.
     * @param minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     * @param genesisEntry Initial relay entry to create first group.
     * @param genesisGroupPubKey Group to respond to the initial relay entry request.
     * @param groupContract Group contract linked to this contract.
     */
    function initialize(uint256 minPayment, uint256 withdrawalDelay, uint256 genesisEntry, bytes memory genesisGroupPubKey, address groupContract)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconFrontendImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _previousEntry = genesisEntry;
        _groupContract = groupContract;

        // Create initial relay entry request. This will allow relayEntry to be called once
        // to trigger the creation of the first group. Requests are removed on successful
        // entries so genesis entry can only be called once.
        _requestCounter++;
        _requests[_requestCounter] = Request(msg.sender, 0, genesisGroupPubKey); 
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconFrontendImplV1"];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated relay request ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 seed) public payable returns (uint256) {
        require(
            msg.value >= _minPayment,
            "Payment is less than required minimum."
        );

        require(
            GroupContract(_groupContract).numberOfGroups() > 0,
            "At least one group needed to serve the request."
        );

        bytes memory groupPubKey = GroupContract(_groupContract).selectGroup(_previousEntry);

        _requestCounter++;

        _requests[_requestCounter] = Request(msg.sender, msg.value, groupPubKey);

        emit RelayEntryRequested(_requestCounter, msg.value, _previousEntry, seed, groupPubKey);
        return _requestCounter;
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 minPayment) public onlyOwner {
        _minPayment = minPayment;
    }

    /**
     * @dev Get the minimum payment that is required before a relay entry occurs.
     */
    function minimumPayment() public view returns(uint256) {
        return _minPayment;
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param requestID The request that started this generation - to tie the results back to the request.
     * @param groupSignature The generated random number.
     * @param groupPubKey Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 requestID, uint256 groupSignature, bytes memory groupPubKey, uint256 previousEntry, uint256 seed) public {

        require(_requests[requestID].groupPubKey.equalStorage(groupPubKey), "Provided group was not selected to produce entry for this request.");
        require(BLS.verify(groupPubKey, abi.encodePacked(previousEntry, seed), bytes32(groupSignature)), "Group signature failed to pass BLS verification.");

        delete _requests[requestID];
        _previousEntry = groupSignature;

        emit RelayEntryGenerated(requestID, groupSignature, groupPubKey, previousEntry, seed);
        GroupContract(_groupContract).runGroupSelection(groupSignature);
    }

    /**
     * @dev Gets the previous relay entry value.
     */
    function previousEntry() public view returns(uint256) {
        return _previousEntry;
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string memory) {
        return "V1";
    }
}
