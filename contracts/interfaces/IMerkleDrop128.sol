// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;
pragma abicoder v1;

// Allows anyone to claim a token if they exist in a merkle root.
interface IMerkleDrop128 {
    // This event is triggered whenever a call to #claim succeeds.
    event Claimed(uint256 index, address account, uint256 amount);

    // Returns the address of the token distributed by this contract.
    function token() external view returns (address);
    // Returns the merkle root of the merkle tree containing account balances available to claim.
    function merkleRoot() external view returns (bytes16);
    // Returns the tree depth of the merkle tree containing account balances available to claim.
    function depth() external view returns (uint256);
    // Returns true if the index has been marked claimed.
    function isClaimed(uint256 index) external view returns (bool);
    // Claim the given amount of the token to the given address. Reverts if the inputs are invalid.
    function claim(address receiver, address account, uint256 amount, bytes calldata merkleProof, bytes calldata signature) external;
}
