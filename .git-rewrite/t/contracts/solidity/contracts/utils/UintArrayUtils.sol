pragma solidity ^0.5.4;


library UintArrayUtils {

    function removeValue(uint256[] storage self, uint256 _value) 
        internal
        returns(uint256[] storage)
    {
        for (uint i = 0; i < self.length; i++) {
            // If value is found in array.
            if (_value == self[i]) {
                // Delete element at index and shift array.
                for (uint j = i; j < self.length-1; j++) {
                    self[j] = self[j+1];
                }
                self.length--;
                i--;
            }
        }
        return self;
    }

    function sort(uint256[] memory array)
        internal pure
        returns (uint256[] memory)
    {
        for (uint i = 0; i < array.length; i++) {
            for (uint j = i+1; j < array.length; j++) {
                if (array[i] > array[j]) {
                    uint256 temp = array[i];
                    array[i] = array[j];
                    array[j] = temp;
                }
            }
        }
        return array;
    }
}
