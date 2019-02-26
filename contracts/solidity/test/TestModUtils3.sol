pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";

contract TestModUtils3 {

    using ModUtils for uint256;

    uint8[][] smallOddPrimesResults = [
        [3, 2, 0],
        [5, 2, 0],
        [5, 3, 0],
        [7, 3, 0],
        [7, 5, 0],
        [7, 6, 0],
        [11, 2, 0],
        [11, 6, 0],
        [11, 7, 0],
        [11, 8, 0],
        [13, 2, 0],
        [13, 5, 0],
        [13, 6, 0],
        [13, 7, 0],
        [13, 8, 0],
        [13, 11, 0],
        [17, 3, 0],
        [17, 5, 0],
        [17, 6, 0],
        [17, 7, 0],
        [17, 11, 0],
        [17, 12, 0],
        [17, 14, 0]
    ];

    function testModSqrtAgainstNonSquares() public {
        uint8 i;
        uint256 a;
        uint256 p;
        uint256 root;

        for(i = 0; i < smallOddPrimesResults.length; i++) {
            p = smallOddPrimesResults[i][0];
            a = smallOddPrimesResults[i][1];
            root = a.modSqrt(p);

            Assert.equal(root, smallOddPrimesResults[i][2], "modSqrt() result differed from list");
        }
    }
}
