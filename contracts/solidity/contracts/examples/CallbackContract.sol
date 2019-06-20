pragma solidity ^0.5.4;


/**
 * @title CallbackContract
 * @dev Example callback contract for Random Beacon.
 */
contract CallbackContract {

    uint256 internal _lastEntry;

    /**
     * @dev Example of a callback method. Method signature can be
     * calculated as bytes4(keccak256("callback(uint256)")
    */
    function callback(uint256 requestResponse)
        public
    {
        _lastEntry = requestResponse;
    }

    /**
     * @dev Returns previous entry.
     */
    function lastEntry() public view returns (uint256)
    {
        return _lastEntry;
    }
}
