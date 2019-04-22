pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";

contract TestModUtils2 {

    using ModUtils for uint256;

    function testModSqrtAgainstList() public {
        uint256 i;
        uint256 a;
        uint256 p;
        uint256 root;

        uint8[3][30] memory smallOddPrimesResults = [
            [3, 1, 1],
            [5, 1, 1],
            [5, 4, 3],
            [7, 1, 1],
            [7, 2, 4],
            [7, 4, 2],
            [11, 1, 1],
            [11, 3, 5],
            [11, 4, 9],
            [11, 5, 4],
            [11, 9, 3],
            [13, 1, 1],
            [13, 3, 9],
            [13, 4, 11],
            [13, 9, 3],
            [13, 10, 7],
            [13, 12, 8],
            [17, 1, 1],
            [17, 2, 6],
            [17, 4, 2],
            [17, 8, 12],
            [17, 9, 14],
            [17, 13, 8],
            [17, 15, 7],
            [17, 16, 4],
            [19, 1, 1],
            [19, 4, 17],
            [19, 5, 9],
            [19, 6, 5],
            [19, 7, 11]
        ];

        for(i = 0; i < smallOddPrimesResults.length; i++) {
            p = smallOddPrimesResults[i][0];
            a = smallOddPrimesResults[i][1];
            root = a.modSqrt(p);

            Assert.equal(root, smallOddPrimesResults[i][2], "modSqrt() result differed from list");
        }
    }
}
