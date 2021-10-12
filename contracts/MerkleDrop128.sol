// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;
pragma abicoder v1;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MerkleDrop128 {
    using SafeERC20 for IERC20;
    event Claimed(uint256 index, address account, uint256 amount);

    address public immutable token;
    bytes16 public immutable merkleRoot;
    uint256 public immutable depth;

    // This is a packed array of booleans.
    mapping(uint256 => uint256) private claimedBitMap;

    constructor(address token_, bytes16 merkleRoot_, uint256 depth_) {
        token = token_;
        merkleRoot = merkleRoot_;
        depth = depth_;
    }

    function isClaimed(uint256 index) public view returns (bool) {
        uint256 claimedWordIndex = index / 256;
        uint256 claimedBitIndex = index % 256;
        uint256 claimedWord = claimedBitMap[claimedWordIndex];
        uint256 mask = (1 << claimedBitIndex);
        return claimedWord & mask == mask;
    }

    function _setClaimed(uint256 index) private {
        uint256 claimedWordIndex = index / 256;
        uint256 claimedBitIndex = index % 256;
        claimedBitMap[claimedWordIndex] = claimedBitMap[claimedWordIndex] | (1 << claimedBitIndex);
    }

    function claim(address receiver, address account, uint256 amount, bytes calldata merkleProof, bytes calldata signature) external {
        // Verify the merkle proof.
        bytes16 node = bytes16(keccak256(abi.encodePacked(account, amount)));
        (bool valid, uint256 index) = _verifyAsm(merkleProof, merkleRoot, node);
        require(valid, "MD: Invalid proof");
        require(ECDSA.recover(keccak256(abi.encodePacked(receiver)), signature) == account, "MD: Invalid signature");
        require(!isClaimed(index), 'MD: Drop already claimed');

        // Mark it claimed and send the token.
        _setClaimed(index);
        IERC20(token).safeTransfer(receiver, amount);
        emit Claimed(index, account, amount);
    }

    function verify(bytes calldata proof, bytes16 root, bytes16 leaf) public view returns (bool valid, uint256 index) {
        return _verifyAsm(proof, root, leaf);
    }

    function _verifyAsm(bytes calldata proof, bytes16 root, bytes16 leaf) private view returns (bool valid, uint256 index) {
        // solhint-disable-next-line no-inline-assembly
        uint256 loopDepth = 0;
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x10)
            let ptr := proof.offset

            for { let end := add(ptr, proof.length) } lt(ptr, end) { ptr := add(ptr, 0x10) } {
                index := shl(1, index)
                let node := calldataload(ptr)

                switch lt(leaf, node)
                case 1 {
                    mstore(mem1, leaf)
                    mstore(mem2, node)
                }
                default {
                    mstore(mem1, node)
                    mstore(mem2, leaf)
                    index := or(1, index)
                }

                leaf := keccak256(mem1, 32)
                loopDepth := add(loopDepth, 1)
            }

            valid := iszero(shr(128, xor(root, leaf)))
        }

        if (loopDepth < depth) {
            index = index << (depth - loopDepth);
        }
    }
}
