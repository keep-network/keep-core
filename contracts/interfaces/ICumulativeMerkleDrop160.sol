// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

// Allows anyone to claim a token if they exist in a merkle root.
interface ICumulativeMerkleDrop160 {
    // This event is triggered whenever a call to #claim succeeds.
    event MerkelRootUpdated(bytes20 oldMerkleRoot, bytes20 newMerkleRoot);
    // This event is triggered whenever a call to #claim succeeds.
    event Claimed(address account, uint256 amount);

    // Returns the address of the token distributed by this contract.
    function token() external view returns (address);
    // Returns the merkle root of the merkle tree containing cumulative account balances available to claim.
    function merkleRoot() external view returns (bytes20);
    // Sets the merkle root of the merkle tree containing cumulative account balances available to claim.
    function setMerkleRoot(bytes20 merkleRoot_) external;
    // Claim the given amount of the token to the given address. Reverts if the inputs are invalid.
    function claim(
        address account,
        uint256 cumulativeAmount,
        bytes20 expectedMerkleRoot,
        bytes calldata merkleProof
    ) external;
}
