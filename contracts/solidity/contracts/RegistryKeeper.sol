pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";


/**
 * @title Registry Keeper
 * @dev An ownable contract to keep registry of approved contracts and roles.
 */
contract RegistryKeeper is Ownable {
    // The Panic Button can disable malicious or malfunctioning contracts
    // that have been previously approved by the Registry Keeper.
    address public panicButton;

    // Operator Contract Upgrader can add approved operator contracts to
    // the service contract and deprecate old ones.
    address public operatorContractUpgrader;

    // The registry of operator contracts
    // 0 - NULL (default), 1 - APPROVED, 2 - DISABLED
    mapping(address => uint256) public operatorContracts;

    modifier onlyPanicButton() {
        require(panicButton == msg.sender, "Not authorized");
        _;
    }

    constructor(address _panicButton) Ownable() public {
        panicButton = _panicButton;
    }

    function approveOperatorContract(address operatorContract) public onlyOwner {
        operatorContracts[operatorContract] = 1;
    }

    function disableOperatorContract(address operatorContract) public onlyPanicButton {
        operatorContracts[operatorContract] = 2;
    }
}
