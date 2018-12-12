pragma solidity ^0.4.24;


library BytesUtils {

    /**
     * @dev Return bytes32 from the bytes array at specified index.
     * @param b Input bytes array.
     * @param i Index at which bytes32 value should be extracted.
     */
    function readBytes32(bytes memory b, uint256 i)
        internal
        pure
        returns (bytes32 result)
    {
        // Skip the first 32 bytes where the length of the bytes array is stored.
        uint256 position = i + 32;

        require(
            b.length >= position,
            "Input bytes array must be 32 length or more."
        );

        /* solium-disable-next-line */
        assembly {
            result := mload(add(b, position))
        }
    }
}
