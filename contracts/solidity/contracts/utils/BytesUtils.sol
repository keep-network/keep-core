pragma solidity ^0.5.4;

library BytesUtils {
    /**
    * @dev returs bytes slice of size 65
    * Modified from:https://github.com/GNSPS/solidity-bytes-utils/blob/master/contracts/BytesLib.sol
    * @param _bytes the bytyes to slice from
    * @param _start where to slice from
    */
    function slice(bytes memory _bytes,uint _start)public pure returns (bytes memory){
        require(_bytes.length >= (_start + 65),"signature not big enough");
        bytes memory tempBytes;
        assembly {
            tempBytes := mload(0x40)
            let lengthmod := and(0x41, 31)
            let mc := add(add(tempBytes, lengthmod), mul(0x20, iszero(lengthmod)))
            let end := add(mc, 0x41)
            for {
                let cc := add(add(add(_bytes, lengthmod), mul(0x20, iszero(lengthmod))), _start)
            } lt(mc, end) {
                mc := add(mc, 0x20)
                cc := add(cc, 0x20)
            } {
                mstore(mc, mload(cc))
            }
            mstore(tempBytes, 0x41)
            mstore(0x40, and(add(mc, 31), not(31)))
        }
        return tempBytes;
    }
}