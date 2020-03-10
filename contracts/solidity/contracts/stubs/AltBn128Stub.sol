pragma solidity ^0.5.4;

import "../cryptography/AltBn128.sol";

library AltBn128Stub {
    function sign(uint256 secretKey) public view returns(uint256, uint256) {
        AltBn128.G1Point memory p_1 = AltBn128.g1HashToPoint(abi.encodePacked(msg.sender));
        AltBn128.G1Point memory p_2 = AltBn128.scalarMultiply(p_1, secretKey);
        return (p_2.x, p_2.y);
    }
}
