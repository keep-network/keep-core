pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";


contract TestModUtils1 {

    using ModUtils for uint256;

    uint8[16] smallOddPrimes = [3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59];

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

    function testLegendre() public {
        uint256 i;
        uint256 j;
        int leg;
        for(i = 0; i < smallOddPrimes.length; i++) {
            for(j = 0; j < 50; j++) {
                leg = ModUtils.legendre(j, smallOddPrimes[i]);
                Assert.isTrue(leg == 0 || leg == 1 || leg == -1, "legendre() should only return [-1, 1]");
            }
        }
    }

    function testModSqrtOf0() public {
        uint256 p;
        uint256 i;
        uint256 zero = 0;

        // a = 0 mod p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            Assert.equal(zero, zero.modSqrt(p), "0 mod p should always equal 0");
        }
    }

    function testModSqrtMultipleOfP() public {
        uint256 p;
        uint256 pMult;
        uint256 i;
        uint256 j;
        uint256 zero = 0;

        // a = 0 mod p
        for(i = 0; i < smallOddPrimes.length; i++) {
            p = smallOddPrimes[i];
            for (j=0; j<20; j++) {
                pMult = p * i;
                Assert.equal(zero, pMult.modSqrt(p), "(n * p) mod p should always equal 0");
            }
        }
    }
}
