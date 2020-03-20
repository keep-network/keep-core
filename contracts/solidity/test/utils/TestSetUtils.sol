pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "../../contracts/utils/SetUtils.sol";


contract TestSetUtils {
    using SetUtils for SetUtils.UintSet;

    mapping(uint256 => SetUtils.UintSet) sets;

    function testContainsWhenEmpty() public {
        SetUtils.UintSet storage set = sets[0];

        Assert.isFalse(set.contains(0), "Empty set should not contain 0");
        Assert.isFalse(set.contains(1), "Empty set should not contain 1");
    }

    function testContainsWithItems() public {
        SetUtils.UintSet storage set = sets[1];

        Assert.isFalse(set.contains(0), "Set should not contain 0 before insertion");
        set.insert(0);
        Assert.isTrue(set.contains(0), "Set should contain 0 after insertion");

        Assert.isFalse(set.contains(1), "Set should not contain 1 before insertion");
        set.insert(1);
        Assert.isTrue(set.contains(1), "Set should contain 1 after insertion");
    }

    function testInsertUnique() public {
        SetUtils.UintSet storage set = sets[2];

        set.insert(1);
        set.insert(0);
        set.insert(2);
        checkThisDamnEquality3(
            set.enumerate(),
            [1, 0, 2],
            "Unique elements should be inserted"
        );
    }

    function testInsertDuplicate() public {
        SetUtils.UintSet storage set = sets[3];

        set.insert(1);
        set.insert(1);
        checkThisDamnEquality1(
            set.enumerate(),
            [1],
            "Duplicate elements should not be inserted"
        );
    }

    function testRemoveLast() public {
        SetUtils.UintSet storage set = sets[4];

        set.insert(1);
        set.insert(0);
        set.insert(2);
        checkThisDamnEquality3(
            set.enumerate(),
            [1, 0, 2],
            "Elements should be inserted correctly"
        );

        set.remove(2);
        checkThisDamnEquality2(
            set.enumerate(),
            [1, 0],
            "Removing the last element should retain the order of others"
        );
    }

    function testRemoveNonLast() public {
        SetUtils.UintSet storage set = sets[4];

        set.insert(1);
        set.insert(0);
        set.insert(2);
        checkThisDamnEquality3(
            set.enumerate(),
            [1, 0, 2],
            "Elements should be inserted correctly"
        );

        set.remove(1);
        checkThisDamnEquality2(
            set.enumerate(),
            [2, 0],
            "Removing an element from the middle should replace it with the last"
        );
    }

    function checkThisDamnEquality1(
        uint256[] memory arrD,
        uint8[1] memory arrS,
        string memory err
    ) internal {
        uint256[] memory arrM = new uint256[](1);
        for (uint256 i = 0; i < arrS.length; i++) {
            arrM[i] = uint256(arrS[i]);
        }

        Assert.equal(arrD, arrM, err);
    }

    function checkThisDamnEquality2(
        uint256[] memory arrD,
        uint8[2] memory arrS,
        string memory err
    ) internal {
        uint256[] memory arrM = new uint256[](2);
        for (uint256 i = 0; i < arrS.length; i++) {
            arrM[i] = uint256(arrS[i]);
        }

        Assert.equal(arrD, arrM, err);
    }

    function checkThisDamnEquality3(
        uint256[] memory arrD,
        uint8[3] memory arrS,
        string memory err
    ) internal {
        uint256[] memory arrM = new uint256[](3);
        for (uint256 i = 0; i < arrS.length; i++) {
            arrM[i] = uint256(arrS[i]);
        }

        Assert.equal(arrD, arrM, err);
    }
}
