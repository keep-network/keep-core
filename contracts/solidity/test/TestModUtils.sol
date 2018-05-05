pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";


contract TestUtils {

    using ModUtils for uint256;

    function testModExponent() public {
        // a simple test
        Assert.equals(21.modExp(2, 5), 1, "");
        // test for overflow - (2 ^ 256 - 1) ^ 2 % alt_bn128_P
        Assert.equals(
            (2 ** 256 - 1).modExp(2, 21888242871839275222246405745257275088696311157297823662689037894645226208583),
            12283109618583340521412061117291584720854994367414008739435419022702680857751,
            "modExp() should not overflow"
        )
    }

    function testModSqrt() public {
        // only work with odd primes p
        Assert.equals(0, 1.modSqrt(2))

        uint256[] smallOddPrimes = [1, 3, 5, 7, 11, 13, 17, 19, 23, 29];
        uint256 p, square;

        // a = 0 mod p
        for(uint256 i = 0; i < smallOddPrimes.legnth; i++) {
            p = smallOddPrimes[i];
            Assert.equals(0, 0.modSqrt(p), "0 mod p should always equal 0")
        }

        // a < p for small p
        for(uint256 i = 0; i < smallOddPrimes.legnth; i++) {
            p = smallOddPrimes[i];
            for(uint256 a = 1; a < p; a++) {
                square = (a * a) % p;
                Assert.equals(a, square.modSqrt(p), "Invalid modular square root for a < p")
            }
        }

        // a > p for small p
        for(uint256 i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for(uint256 a = p + 1; a < p + 10; a++) {
                square = (a * a) % p;
                Assert.equals(a, square.modSqrt(p), "Invalid modular square root for a > p");
            }
        }

        // TODO tests with larger p
    }
}
