pragma solidity ^0.5.4;


/**
 * @title Registry
 * @dev Governance owned registry of approved contracts and roles.
 */
contract Registry {

    enum ContractStatus { New, Approved, Disabled }

    // Governance role is to enable recovery from key compromise by rekeying other roles.
    address internal governance;

    // Registry Keeper maintains approved operator contracts. Each operator
    // contract must be approved before it can be authorized by a staker or
    // used by a service contract.
    address internal registryKeeper;

    // The Panic Button can disable malicious or malfunctioning contracts
    // that have been previously approved by the Registry Keeper.
    address internal panicButton;

    // Each service contract has a Operator Contract Upgrader whose purpose
    // is to manage operator contracts for that specific service contract.
    // The Operator Contract Upgrader can add new operator contracts to the
    // service contract’s operator contract list, and deprecate old ones.
    mapping(address => address) public operatorContractUpgraders;

    // The registry of operator contracts
    mapping(address => ContractStatus) public operatorContracts;

    event OperatorContractApproved(address operatorContract);

    modifier onlyGovernance() {
        require(governance == msg.sender, "Not authorized");
        _;
    }

    modifier onlyRegistryKeeper() {
        require(registryKeeper == msg.sender, "Not authorized");
        _;
    }

    modifier onlyPanicButton() {
        require(panicButton == msg.sender, "Not authorized");
        _;
    }

    constructor() public {
        governance = msg.sender;
        registryKeeper = msg.sender;
        panicButton = msg.sender;
    }

    function setGovernance(address _governance) public onlyGovernance {
        governance = _governance;
    }

    function setRegistryKeeper(address _registryKeeper) public onlyGovernance {
        registryKeeper = _registryKeeper;
    }

    function setPanicButton(address _panicButton) public onlyGovernance {
        panicButton = _panicButton;
    }

    function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) public onlyGovernance {
        operatorContractUpgraders[_serviceContract] = _operatorContractUpgrader;
    }

    function approveOperatorContract(address operatorContract) public onlyRegistryKeeper {
        require(
            isNewOperatorContract(operatorContract),
            "Only new operator contracts can be approved"
        );

        operatorContracts[operatorContract] = ContractStatus.Approved;
        emit OperatorContractApproved(operatorContract);
    }

    function disableOperatorContract(address operatorContract) public onlyPanicButton {
        require(
            isApprovedOperatorContract(operatorContract),
            "Only approved operator contracts can be disabled"
        );

        operatorContracts[operatorContract] = ContractStatus.Disabled;
    }

    function isNewOperatorContract(address operatorContract) public view returns (bool) {
        return operatorContracts[operatorContract] == ContractStatus.New;
    }

    function isApprovedOperatorContract(address operatorContract) public view returns (bool) {
        return operatorContracts[operatorContract] == ContractStatus.Approved;
    }

    function operatorContractUpgraderFor(address _serviceContract) public view returns (address) {
        return operatorContractUpgraders[_serviceContract];
    }
}
