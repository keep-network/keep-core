pragma solidity ^0.5.4;

library OperatorUtils {
    uint256 constant BLOCKHEIGHT_WIDTH = 64;
    uint256 constant AMOUNT_WIDTH = 128;

    uint256 constant BLOCKHEIGHT_MAX = (2**BLOCKHEIGHT_WIDTH) - 1;
    uint256 constant AMOUNT_MAX = (2**AMOUNT_WIDTH) - 1;

    uint256 constant CREATION_SHIFT = BLOCKHEIGHT_WIDTH;
    uint256 constant AMOUNT_SHIFT = 2 * BLOCKHEIGHT_WIDTH;

    function pack(
        uint256 amount,
        uint256 createdAt,
        uint256 undelegatedAt
    ) internal pure returns (uint256) {
        uint256 a = (amount & AMOUNT_MAX) << AMOUNT_SHIFT;
        uint256 c = (createdAt & BLOCKHEIGHT_MAX) << CREATION_SHIFT;
        uint256 u = undelegatedAt & BLOCKHEIGHT_MAX;
        return (a | c | u);
    }

    function unpack(uint256 packedParams) internal pure returns (
        uint256 amount,
        uint256 createdAt,
        uint256 undelegatedAt
    ) {
        amount = getAmount(packedParams);
        createdAt = getCreationBlock(packedParams);
        undelegatedAt = getUndelegationBlock(packedParams);

        return (amount, createdAt, undelegatedAt);
    }

    function getAmount(uint256 packedParams) internal pure returns (uint256) {
        return (packedParams >> AMOUNT_SHIFT) & AMOUNT_MAX;
    }

    function setAmount(
        uint256 packedParams,
        uint256 amount
    ) internal pure returns (uint256) {
        return pack(
            amount,
            getCreationBlock(packedParams),
            getUndelegationBlock(packedParams)
        );
    }

    function getCreationBlock(uint256 packedParams) internal pure returns (uint256) {
        return (packedParams >> CREATION_SHIFT) & BLOCKHEIGHT_MAX;
    }

    function setCreationBlock(
        uint256 packedParams,
        uint256 creationBlockheight
    ) internal pure returns (uint256) {
        return pack(
            getAmount(packedParams),
            creationBlockheight,
            getUndelegationBlock(packedParams)
        );
    }

    function getUndelegationBlock(uint256 packedParams) internal pure returns (uint256) {
        return packedParams & BLOCKHEIGHT_MAX;
    }

    function setUndelegationBlock(
        uint256 packedParams,
        uint256 undelegationBlockheight
    ) internal pure returns (uint256) {
        return pack(
            getAmount(packedParams),
            getCreationBlock(packedParams),
            undelegationBlockheight
        );
    }
}
