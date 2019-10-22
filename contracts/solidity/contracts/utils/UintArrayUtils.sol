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

    function sort(uint256[] memory data) internal pure returns(uint256[] memory) {
       quickSort(data, int(0), int(data.length - 1));
       return data;
    }

    // Quicksort is a divide-and-conquer algorithm for sorting. It works by partitioning
    // an array into two subarrays, then sorting the subarrays independently.
    function quickSort(uint256[] memory arr, int lo, int hi) internal pure {
        if (hi <= lo) return;
        int j = partition(arr, lo, hi);
        quickSort(arr, lo, j-1);
        quickSort(arr, j+1, hi);
    }

    function partition(uint256[] memory arr, int lo, int hi) internal pure returns(int) {
        int i = lo; // left scan indice
        int j = hi + 1; // right scan indice
        uint256 pivot = arr[uint256(lo)]; // partitioning item
        while(true) {
            // scan left
            while (arr[uint256(++i)] < pivot) {
                if (i == hi) break;
            }
            // scan right
            while (pivot < arr[uint256(--j)]) {
                if (j == lo) break;
            }
            if (i >= j) break;
            (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
        }
        (arr[uint256(lo)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(lo)]);

        return j;
    }

}
