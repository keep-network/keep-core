pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";
import "../contracts/AltBn128.sol";

contract TestAltBn128 {

    function testHashing() public {
        string memory hello = "hello!";
        string memory goodbye = "goodbye.";
        uint256 p_1_x;
        uint256 p_1_y;
        uint256 p_2_x;
        uint256 p_2_y;
        (p_1_x, p_1_y) = AltBn128.g1HashToPoint(bytes(hello));
        (p_2_x, p_2_y) = AltBn128.g1HashToPoint(bytes(goodbye));

        Assert.isNotZero(p_1_x, "X should not equal 0 in a hashed point.");
        Assert.isNotZero(p_1_y, "Y should not equal 0 in a hashed point.");
        Assert.isNotZero(p_2_x, "X should not equal 0 in a hashed point.");
        Assert.isNotZero(p_2_y, "Y should not equal 0 in a hashed point.");

        Assert.isTrue(isOnCurve(p_1_x, p_1_y), "Hashed points should be on the curve.");
        Assert.isTrue(isOnCurve(p_2_x, p_2_y), "Hashed points should be on the curve.");
    }

    function isOnCurve(uint256 x, uint256 y) internal returns (bool) {

        return ModUtils.modExp(y, 2, AltBn128.getP()) == ModUtils.modExp(x, 3, AltBn128.getP()) + 3;
    }
}
