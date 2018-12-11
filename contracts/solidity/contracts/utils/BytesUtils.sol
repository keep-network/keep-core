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
        require(
            b.length >= i + 32,
            "Input bytes array must be 32 length or more."
        );

        i += 32;

        /* solium-disable-next-line */
        assembly {
            result := mload(add(b, i))
        }
    }
}
