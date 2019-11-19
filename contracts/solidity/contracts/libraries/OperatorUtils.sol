pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../utils/AddressArrayUtils.sol";

library OperatorUtils {
    using SafeMath for uint256;
    using AddressArrayUtils for address[];

    struct Storage {
        address owner;
        address[] serviceContracts;
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
}
