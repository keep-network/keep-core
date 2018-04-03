pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";


/**
 * @title Random Beacon Proxy
 * @dev An ownable proxy contract to provide upgradable Random Beacon functionality.
 * Upgraded contract is set as the active contract.
 * All calls to this contract are delegated to the active contract.
 */
contract RandomBeaconProxy is Ownable {
    address public activeContract;

    event ContractUpdated(address indexed _contract);

    function RandomBeaconProxy(address _contract) {
        activeContract = _contract;
    }

    /**
     * @dev Delegate call to the active contract.
     */
    function() payable {
        require(activeContract.delegatecall(msg.data));
    }

    /**
     * @dev Upgrade current contract.
     * @param _contract The address of a new contract.
     */
    function updateContract(address _contract) 
        public
        onlyOwner
    {
        activeContract = _contract;
        ContractUpdated(activeContract);
    }
}
