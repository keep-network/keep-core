pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../utils/AddressArrayUtils.sol";
import "../../TokenStaking.sol";


library ContractReferences {
    using SafeMath for uint256;
    using AddressArrayUtils for address[];

    struct Storage {
        address owner;
        address[] serviceContracts;
        // TODO: replace with a secure authorization protocol (addressed in RFC 11).
        address stakingContract;
    }

    /**
     * @dev Adds service contract
     * @param serviceContract Address of the service contract.
     */
    function addServiceContract(
        Storage storage self,
        address serviceContract,
        address sender
    ) public {
        require(self.owner == sender, "Caller is not the owner");
        self.serviceContracts.push(serviceContract);
    }

    /**
     * @dev Removes service contract
     * @param serviceContract Address of the service contract.
     */
    function removeServiceContract(
        Storage storage self,
        address serviceContract,
        address sender
    ) public {
        require(self.owner == sender, "Caller is not the owner");
        self.serviceContracts.removeAddress(serviceContract);
    }

    /**
     * @dev Checks if service contract exist
     * @param serviceContract Address of the service contract.
     */
    function isServiceContract(
        Storage storage self,
        address serviceContract
    ) public view returns (bool) {
        return self.serviceContracts.contains(serviceContract);
    }

    /**
     * @dev Gets latest added service contract
     */
    function latestServiceContract(
        Storage storage self
    ) public view returns (address) {
        return self.serviceContracts[self.serviceContracts.length.sub(1)];
    }

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function stakeBalanceOf(
        Storage storage self,
        address _address
    ) public view returns (uint256) {
        return TokenStaking(self.stakingContract).balanceOf(_address);
    }

    /**
     * @dev Gets the magpie for the specified operator address.
     * @return Magpie address.
     */
    function magpieOf(
        Storage storage self,
        address _address
    ) public view returns (address payable) {
        return TokenStaking(self.stakingContract).magpieOf(_address);
    }
}
