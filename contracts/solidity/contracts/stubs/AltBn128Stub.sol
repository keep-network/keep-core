pragma solidity ^0.5.4;

import "../cryptography/AltBn128.sol";

library AltBn128Stub {
    function sign(uint256 secretKey) public view returns(uint256, uint256) {
        return AltBn128.sign(abi.encodePacked(msg.sender), secretKey);
    }

    function g1HashToPoint(bytes memory message) public view returns(bytes memory) {
        AltBn128.G1Point memory point = AltBn128.g1HashToPoint(message);
        bytes memory messageHashed = new bytes(64);
        bytes32 x = bytes32(point.x);
        bytes32 y = bytes32(point.y);
        assembly {
            mstore(add(messageHashed, 32), x)
            mstore(add(messageHashed, 64), y)
        }

        return messageHashed;
    }
}
