pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";

contract TestModUtils4 {

    using ModUtils for uint256;

    uint8[16] smallOddPrimes = [3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59];

    function testModSqrtALessThanP() public {
        uint256 p;
        uint256 root;
        uint256 i;
        uint256 a;

        // a < p for small p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for(a = 1; a < p; a++) {
                root = a.modSqrt(p);
                if (root != 0) {
                    Assert.equal(a % p, (root * root) % p, "Invalid modular square root for a < p");
                }
            }
        }
    }

    function testModSqrtAGreaterThanP() public {
        uint256 p;
        uint256 root;
        uint8 i;
        uint256 a;

        // a > p for small p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for(a = p + 1; a < p + 10; a++) {
                root = a.modSqrt(p);
                if (root != 0) {
                    Assert.equal(a % p, (root * root) % p, "Invalid modular square root for a > p");
                }
            }
        }
    }
}
