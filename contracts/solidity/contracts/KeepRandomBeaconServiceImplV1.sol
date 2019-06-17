pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./utils/AddressArrayUtils.sol";
import "./DelayedWithdrawal.sol";


interface OperatorContract {
    function sign(uint256 entryId, uint256 seed, uint256 previousEntry) payable external;
    function numberOfGroups() external view returns(uint256);
    function selectGroup(uint256 previousEntry) external returns(bytes memory);
}


/**
 * @title KeepRandomBeaconServiceImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconServiceImplV1 is Ownable, DelayedWithdrawal {

    using AddressArrayUtils for address[];

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey); 
    event RelayEntryGenerated(uint256 entryId, uint256 entry);

    uint256 internal _minPayment;
    uint256 internal _previousEntry;
    uint256 internal _entryCounter;

    address[] internal _operatorContracts;
    mapping (address => uint256) internal _operatorContractNumberOfGroups;

    mapping (string => bool) internal _initialized;

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
     * @param operatorContract Operator contract linked to this contract.
     */
    function initialize(uint256 minPayment, uint256 withdrawalDelay, address operatorContract)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconServiceImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _operatorContracts.push(operatorContract);
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconServiceImplV1"];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated entry Id.
     */
    function requestRelayEntry(uint256 seed) public payable returns (uint256) {
        require(
            msg.value >= _minPayment,
            "Payment is less than required minimum."
        );

        _entryCounter++;

        // TODO: select operator contract
        // TODO: Figure out pricing, if we decide to pass payment to the backed use this instead:
        // OperatorContract(_operatorContracts[0]).sign.value(msg.value)(_entryCounter, seed, _previousEntry);
        OperatorContract(_operatorContracts[0]).sign(_entryCounter, seed, _previousEntry);
        return _entryCounter;
    }

    /**
     * @dev Store valid entry returned by operator contract and call customer specified callback if required.
     * @param entryId Entry id tracked internally by this contract.
     * @param entry The generated random number.
     */
    function entryCreated(uint256 entryId, uint256 entry) public {
        require(
            _operatorContracts.contains(msg.sender),
            "Only authorized operator contract can call relay entry."
        );

        _previousEntry = entry;
        emit RelayEntryGenerated(entryId, entry);
        // TODO: customer-specified callback
    }

    /**
     * @dev Store number of groups returned by operator contract.
     * @param numberOfGroups Number of groups.
     */
    function groupCreated(uint256 numberOfGroups) public {
        require(
            _operatorContracts.contains(msg.sender),
            "Only authorized operator contract can call groupCreated."
        );

        _operatorContractNumberOfGroups[msg.sender] = numberOfGroups;
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
