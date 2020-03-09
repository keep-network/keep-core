pragma solidity ^0.5.4;
import "../utils/OperatorParams.sol";

contract OperatorParamsStub {
    using OperatorParams for uint256;

    function publicPack(
        uint256 amount,
        uint256 createdAt,
        uint256 undelegatedAt
    ) public pure returns (uint256) {
        return OperatorParams.pack(amount, createdAt, undelegatedAt);
    }

    function publicGetAmount(uint256 packed) public pure returns (uint256) {
        return packed.getAmount();
    }

    function publicSetAmount(uint256 packed, uint256 amount) public pure returns (uint256) {
        return packed.setAmount(amount);
    }

    function publicGetCreationBlock(uint256 packed) public pure returns (uint256) {
        return packed.getCreationBlock();
    }

    function publicSetCreationBlock(uint256 packed, uint256 creationBlock) public pure returns (uint256) {
        return packed.setCreationBlock(creationBlock);
    }

    function publicGetUndelegationBlock(uint256 packed) public pure returns (uint256) {
        return packed.getUndelegationBlock();
    }

    function publicSetUndelegationBlock(uint256 packed, uint256 undelegationBlock) public pure returns (uint256) {
        return packed.setUndelegationBlock(undelegationBlock);
    }
}
