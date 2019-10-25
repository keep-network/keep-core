pragma solidity ^0.5.4;


interface OperatorContract {
    function relayEntryTimeout() external view returns(uint256);
}

/**
 * @title KeepRandomBeaconOperatorLinkedContract
 * @dev A base contract to allow contract to be linked with operator contract.
 */
contract KeepRandomBeaconOperatorLinkedContract {

    // Contract owner.
    address public owner;

    // Operator contract that is linked to this contract.
    address public operatorContract;

    // Duplicated constant from operator contract to avoid extra call.
    // The value is set when the operator contract is added.
    uint256 public relayEntryTimeout;

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner == msg.sender, "Caller is not the owner.");
        _;
    }

    /**
     * @dev Throws if called by any account other than the authorized address.
     */
    modifier onlyOperatorContract() {
        require(operatorContract == msg.sender, "Caller is not authorized.");
        _;
    }

    /**
     * @dev Initializes the contract with deployer as the contract owner.
     */
    constructor() public {
        owner = msg.sender;
    }

    /**
     * @dev Sets operator contract.
     */
    function setOperatorContract(address _operatorContract) public onlyOwner {
        require(operatorContract == address(0), "Operator contract can only be set once.");
        operatorContract = _operatorContract;
        relayEntryTimeout = OperatorContract(operatorContract).relayEntryTimeout();
    }
}
