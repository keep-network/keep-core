pragma solidity ^0.4.18;


library AddressArrayUtils {

    function isFound(address[] self, address _address) 
        internal
        pure
        returns(bool)
    {
        for (uint i = 0; i < self.length; i++) {
            if (_address == self[i]) {
                return true;
            }
        }
        return false;
    }

    function removeAddress(address[] storage self, address _address) 
        internal
        returns(address[])
    {
        for (uint i = 0; i < self.length; i++) {
            // If address is found in array.
            if (_address == self[i]) {
                // Delete element at index and shift array.
                for (uint j = i; j < self.length-1; j++) {
                    self[j] = self[j+1];
                }
                delete self[self.length-1];
                self.length--;
            }
        }
        return self;
    }
}
