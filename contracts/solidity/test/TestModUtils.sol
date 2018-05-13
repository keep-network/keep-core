pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";


contract TestModUtils {

    using ModUtils for uint256;

    function testModExponent() public {
        uint256 a = 21;
        // a simple test
        Assert.equal(a.modExp(2, 5), 1, "");
        // test for overflow - (2 ^ 256 - 1) ^ 2 % alt_bn128_P
        uint256 almostOverflow = (2 ** 256 - 1);
        Assert.equal(
            almostOverflow.modExp(2, 21888242871839275222246405745257275088696311157297823662689037894645226208583),
            12283109618583340521412061117291584720854994367414008739435419022702680857751,
            "modExp() should not overflow"
        );
    }



    function testModSqrt() public {
        uint8[10] memory smallOddPrimes = [1, 3, 5, 7, 11, 13, 17, 19, 23, 29];
        uint256 p;
        uint256 square;
        uint256 i;
        uint256 a;
        uint256 zero = 0;

        // a = 0 mod p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            Assert.equal(zero, zero.modSqrt(p), "0 mod p should always equal 0");
        }

        // a < p for small p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for(a = 1; a < p; a++) {
                square = (a * a) % p;
                Assert.equal(a, square.modSqrt(p), "Invalid modular square root for a < p");
            }
        }

        // a > p for small p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for(a = p + 1; a < p + 10; a++) {
                square = (a * a) % p;
                Assert.equal(a, square.modSqrt(p), "Invalid modular square root for a > p");
            }
        }

        // TODO tests with larger p
        // TODO test throws with non-odd prime p
    }
}
