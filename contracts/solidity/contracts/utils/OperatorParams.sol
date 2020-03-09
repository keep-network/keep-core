pragma solidity ^0.5.4;

library OperatorParams {
    // OperatorParams packs values that are commonly used together
    // into a single uint256 to reduce the cost functions
    // like querying eligibility.
    //
    // An OperatorParams uint256 contains:
    // - the operator's staked token amount (uint128)
    // - the operator's creation block (uint64)
    // - the operator's undelegation block (uint64)
    //
    // These are packed as [amount | createdAt | undelegatedAt]
    //
    // Staked KEEP is stored in an uint128,
    // which is sufficient ecause KEEP tokens have 18 decimals (2^60)
    // and there will be at most 10^9 KEEP in existence (2^30).
    //
    // Creation and undelegation times are stored in an uint64 each.
    // Because blocks are created every 10 to 15 seconds,
    // there are less than 4 million blocks per year (2^22).
    // Thus uint64s would be sufficient for a trillion years.
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
        // Check for staked amount overflow.
        // We shouldn't actually ever need this.
        require(
            amount <= AMOUNT_MAX,
            "amount uint128 overflow"
        );
        // Bitwise OR the blockheights together.
        // The resulting number is equal or greater than either,
        // and tells if we have a bit set outside the 64 available bits.
        require(
            (createdAt | undelegatedAt) <= BLOCKHEIGHT_MAX,
            "blockheight uint64 overflow"
        );
        uint256 a = amount << AMOUNT_SHIFT;
        uint256 c = createdAt << CREATION_SHIFT;
        uint256 u = undelegatedAt;
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
    }

    function getAmount(uint256 packedParams)
        internal pure returns (uint256) {
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

    function getCreationBlock(uint256 packedParams)
        internal pure returns (uint256) {
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

    function getUndelegationBlock(uint256 packedParams)
        internal pure returns (uint256) {
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
