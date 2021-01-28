pragma solidity ^0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

library BinarySearch {
    using SafeMath for uint256;

    /// @notice Find the index of the first item
    /// whose value equals or exceeds the target,
    /// or alternatively the number of items whose value is below the target.
    /// @param getValue A function taking an index as the argument
    /// and returning the corresponding value.
    /// Values are expected to be in ascending order,
    /// i.e. with indices `i` and `j` where `i < j`,
    /// `getValue(i) <= getValue(j)`.
    /// @param itemCount The number of items the search is performed over.
    /// `getValue(index)` must return a valid value for any `index < itemCount`.
    /// @return The lowest index where `getValue(index) >= target`.
    /// If no such index exists in the given range, return `itemCount`.
    function find(
        function (uint256) view returns (uint256) getValue,
        uint256 itemCount,
        uint256 target
    ) internal view returns (uint256) {
        if (itemCount == 0) { return 0; }

        uint256 lowerBound = 0;
        uint256 lowerBoundValue = getValue(0);
        if (lowerBoundValue >= target) { return 0; }

        uint256 upperBound = itemCount.sub(1);
        uint256 upperBoundValue = getValue(upperBound);
        if (upperBoundValue < target) { return itemCount; }

        uint256 length = upperBound.sub(lowerBound);

        while (length > 1) {
            // upper bound >= lower bound + 2
            // middle > lower bound
            uint256 middle = lowerBound.add(length.div(2));
            uint256 middleValue = getValue(middle);

            if (middleValue >= target) {
                upperBound = middle;
                upperBoundValue = middleValue;
            } else {
                lowerBound = middle;
                lowerBoundValue = middleValue;
            }
            length = upperBound.sub(lowerBound);
        }
        return upperBound;
    }
}
