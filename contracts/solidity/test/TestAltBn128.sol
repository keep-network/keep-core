pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/utils/ModUtils.sol";
import "../contracts/AltBn128.sol";

contract TestAltBn128 {

    uint256[2] g1 = AltBn128.g1();
    uint256[4] g2 = AltBn128.g2();

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

    function isOnCurve(uint256 x, uint256 y) internal view returns (bool) {
        return ModUtils.modExp(y, 2, AltBn128.getP()) == (ModUtils.modExp(x, 3, AltBn128.getP()) + 3) % AltBn128.getP();
    }

    function testHashAndAdd() public {
        string memory hello = "hello!";
        string memory goodbye = "goodbye.";
        uint256 p_1_x;
        uint256 p_1_y;
        uint256 p_2_x;
        uint256 p_2_y;
        (p_1_x, p_1_y) = AltBn128.g1HashToPoint(bytes(hello));
        (p_2_x, p_2_y) = AltBn128.g1HashToPoint(bytes(goodbye));

        uint256 p_3_x;
        uint256 p_3_y;
        uint256 p_4_x;
        uint256 p_4_y;

        (p_3_x, p_3_y) = AltBn128.add([p_1_x, p_1_y], [p_2_x, p_2_y]);
        (p_4_x, p_4_y) = AltBn128.add([p_2_x, p_2_y], [p_1_x, p_1_y]);

        Assert.equal(p_3_x, p_4_x, "Point addition should be commutative.");
        Assert.equal(p_3_y, p_4_y, "Point addition should be commutative.");

        Assert.isTrue(isOnCurve(p_3_x, p_3_y), "Added points should be on the curve.");
    }

    function testHashAndScalarMultiply() public {
        string memory hello = "hello!";
        uint256 p_1_x;
        uint256 p_1_y;
        uint256 p_2_x;
        uint256 p_2_y;
        (p_1_x, p_1_y) = AltBn128.g1HashToPoint(bytes(hello));

        (p_2_x, p_2_y) = AltBn128.scalarMultiply([p_1_x, p_1_y], 12);

        Assert.isTrue(isOnCurve(p_2_x, p_2_y), "Multiplied point should be on the curve.");
    }

    uint256[2][] randomG1 = [
        [19985462441994274044747034318046506954527006367483173410210086890020894468080, 18435086518936643964830423002803816020906755142322386776411266459735121477493],
        [5020462286181323390508118928832214575747271080433231325350949221928552771006, 6185819871141660402526014503512646294685393996180865467240675728617780703293],
        [6443569433573553122968863711873646857259386088199162681359502284812487407640, 5378163071719228060986147516945215302807920666474786022068644421654848367565],
        [19923401560169709235429596406611561407855841155398968552083379042854674266499, 10525710049852251332517421942831871137061760270860191079861735894813959253055],
        [19708536568727021605314080150939514846681180496259214578335284564769723419938, 10692356507990003585226828000662725800587874779874382732784071410185542028439],
        [9372321588728408099991690236147178727235677209811335191981120085012199642559, 11142558497436993571688400857990083465974854590891619188837196531526982135288],
        [5965886725029153696599727822391947370059044516209856603831046549655428439060, 10175397348860086021525298362240324520688370458967443904613437789517586359962],
        [20748498912264019189558145442056089284703240490771913074152837182874426945993, 18057592905480302483449076150943157907511999106688668826058046434471622799474],
        [7477907739342510339540973467783537984932469471333402963930842749621045686487, 1179596217276931579251786249459263438406283313229247981371951224605996910316],
        [3386341017431964271492464889305868556498227248869025090652509478713128447791, 1836930069368635496176332910536803390892441983393373783218213609800061729358],
        [15369136978943154357167241223632015935727499997608268363280373457134516735375, 13963714284290182895189799343803541866405397472877283227980653081636863075815]
    ];

    function testAdd() public {
        uint i;
        uint8 j;

        uint256 p_1_x;
        uint256 p_1_y;
        uint256 p_2_x;
        uint256 p_2_y;

        for (i = 0; i < randomG1.length; i++) {
            for (j = 0; j < randomG1.length; j++) {

                (p_1_x, p_1_y) = AltBn128.add(randomG1[i], randomG1[j]);
                (p_2_x, p_2_y) = AltBn128.add(randomG1[j], randomG1[i]);

                Assert.equal(p_1_x, p_2_x, "Point addition should be commutative.");
                Assert.equal(p_1_y, p_2_y, "Point addition should be commutative.");

                Assert.isTrue(isOnCurve(p_1_x, p_1_y), "Added points should be on the curve.");
            }
        }
    }

    function testScalarMultiply() public {
        uint i;
        uint j;

        uint256 p_1_x;
        uint256 p_1_y;
        uint256 p_2_x;
        uint256 p_2_y;

        for (i = 1; i < randomG1.length; i++) {
            (p_1_x, p_1_y) = AltBn128.scalarMultiply(randomG1[i], i);

            Assert.isTrue(isOnCurve(p_1_x, p_1_y), "Multiplied point should be on the curve.");

            (p_2_x, p_2_y) = (randomG1[i][0], randomG1[i][1]);
            for (j = 1; j < i; j++) {
                (p_2_x, p_2_y) = AltBn128.add([p_2_x, p_2_y], randomG1[i]);
            }

            Assert.equal(p_1_x, p_2_x, "Scalar multiplication should match repeat addition.");
            Assert.equal(p_1_y, p_2_y, "Scalar multiplication should match repeat addition.");
        }
    }

    function testBasicPairing() public {
        bool result = AltBn128.pairing(g1, g2, [g1[0], AltBn128.getP() - g1[1]], g2);
        Assert.isTrue(result, "Basic pairing check should succeed.");
    }
}
