pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/utils/AddressArrayUtils.sol";
import "../contracts/utils/UintArrayUtils.sol";


contract TestUtils {  
  
    using AddressArrayUtils for address[];
    using UintArrayUtils for uint256[];

    address[] public addresses;
    uint256[] public values;

    function testCanRemoveAddress() public {
        bool exists = false;
        addresses.push(0x1111111111111111111111111111111111111111);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x3333333333333333333333333333333333333333);

        addresses.removeAddress(0x2222222222222222222222222222222222222222);
        for (uint i = 0; i < addresses.length; i++) {
            if (addresses[i] == 0x2222222222222222222222222222222222222222) {
                exists = true;
            }
        }

        Assert.equal(exists, false, "All occurrences of the address should be removed from the array.");
    }

    function testCanRemoveValue() public {
        bool exists = false;
        values = [1, 2, 2, 3];

        values.removeValue(2);
        for (uint i = 0; i < values.length; i++) {
            if (values[i] == 2) {
                exists = true;
            }
        }

        Assert.equal(exists, false, "All occurrences of the value should be removed from the array.");
    }
}
