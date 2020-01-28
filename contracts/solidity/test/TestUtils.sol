pragma solidity 0.5.7;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/utils/AddressArrayUtils.sol";
import "../contracts/utils/UintArrayUtils.sol";


contract TestUtils {  
  
    using AddressArrayUtils for address[];
    using UintArrayUtils for uint256[];

    address[] public addresses;
    uint256[] public values;

    function testCanHandleEmptyArray() public {
        addresses.removeAddress(0x2222222222222222222222222222222222222222);
        Assert.equal(addresses.length, 0, "Empty array should stay unchanged on attempt to remove address from it.");
    }

    function testCanRemoveAddressFromSingleElementArray() public {
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.removeAddress(0x2222222222222222222222222222222222222222);
        Assert.equal(addresses.length, 0, "Occurrence of address in a single element array should be removed.");
    }

    function testCanRemoveIdenticalAddresses() public {
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.removeAddress(0x2222222222222222222222222222222222222222);
        Assert.equal(addresses.length, 0, "All occurrences should be removed and array length should become 0.");
    }

    function testCanRemoveAddress() public {
        bool exists = false;
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x1111111111111111111111111111111111111111);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x2222222222222222222222222222222222222222);
        addresses.push(0x3333333333333333333333333333333333333333);
        addresses.push(0x2222222222222222222222222222222222222222);

        addresses.removeAddress(0x2222222222222222222222222222222222222222);
        for (uint i = 0; i < addresses.length; i++) {
            if (addresses[i] == 0x2222222222222222222222222222222222222222) {
                exists = true;
            }
        }

        Assert.equal(exists, false, "All occurrences of the address should be removed from the array.");
        Assert.equal(addresses.length, 2, "Array length should change accordingly.");
    }

    function testCanHandleEmptyValueArray() public {
        values.removeValue(2);
        Assert.equal(values.length, 0, "Empty array should stay unchanged on attempt to remove value from it.");
    }

    function testCanRemoveValueFromSingleElementArray() public {
        values = [2];
        values.removeValue(2);
        Assert.equal(values.length, 0, "Occurrence of a value in a single element array should be removed.");
    }

    function testCanRemoveIdenticalValues() public {
        values = [2, 2, 2];
        values.removeValue(2);
        Assert.equal(values.length, 0, "All occurrences should be removed and array length should become 0.");
    }

    function testCanRemoveValue() public {
        bool exists = false;
        values = [2, 1, 2, 2, 3, 2];

        values.removeValue(2);
        for (uint i = 0; i < values.length; i++) {
            if (values[i] == 2) {
                exists = true;
            }
        }

        Assert.equal(exists, false, "All occurrences of the value should be removed from the array.");
        Assert.equal(addresses.length, 2, "Array length should change accordingly.");
    }
}
