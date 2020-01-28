pragma solidity 0.5.7;


/**
 * @title Registry
 * @dev Governance owned registry of approved contracts and roles.
 */
contract Registry {
    // Governance role is to enable recovery from key compromise by rekeying other roles.
    address internal governance;

    // Registry Keeper maintains approved operator contracts. Each operator
    // contract must be approved before it can be authorized by a staker or
    // used by a service contract.
    address internal registryKeeper;

    // The Panic Button can disable malicious or malfunctioning contracts
    // that have been previously approved by the Registry Keeper.
    address internal panicButton;

    // Operator Contract Upgrader can add approved operator contracts to
    // the service contract and deprecate old ones.
    address public operatorContractUpgrader;

    // The registry of operator contracts
    // 0 - NULL (default), 1 - APPROVED, 2 - DISABLED
    mapping(address => uint256) public operatorContracts;

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
        operatorContractUpgrader = msg.sender;
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

    function setOperatorContractUpgrader(address _operatorContractUpgrader) public onlyGovernance {
        operatorContractUpgrader = _operatorContractUpgrader;
    }

    function approveOperatorContract(address operatorContract) public onlyRegistryKeeper {
        operatorContracts[operatorContract] = 1;
    }

    function disableOperatorContract(address operatorContract) public onlyPanicButton {
        operatorContracts[operatorContract] = 2;
    }

    function isApprovedOperatorContract(address operatorContract) public view returns (bool) {
        return operatorContracts[operatorContract] == 1;
    }
}
