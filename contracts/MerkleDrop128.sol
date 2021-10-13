// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;
pragma abicoder v1;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import "./interfaces/IMerkleDrop128.sol";

contract MerkleDrop128 is IMerkleDrop128 {
    using SafeERC20 for IERC20;

    address public immutable override token;
    bytes16 public immutable override merkleRoot;
    uint256 public immutable override depth;

    // This is a packed array of booleans.
    mapping(uint256 => uint256) private _claimedBitMap;

    constructor(address token_, bytes16 merkleRoot_, uint256 depth_) {
        token = token_;
        merkleRoot = merkleRoot_;
        depth = depth_;
    }

    function claim(address receiver, uint256 amount, bytes calldata merkleProof, bytes calldata signature) external override {
        bytes32 signedHash = ECDSA.toEthSignedMessageHash(keccak256(abi.encodePacked(receiver)));
        address account = ECDSA.recover(signedHash, signature);
        // Verify the merkle proof.
        bytes16 node = bytes16(keccak256(abi.encodePacked(account, amount)));
        (bool valid, uint256 index) = _verifyAsm(merkleProof, merkleRoot, node);
        require(valid, "MD: Invalid proof");
        _invalidate(index);
        IERC20(token).safeTransfer(receiver, amount);
    }

    function verify(bytes calldata proof, bytes16 root, bytes16 leaf) external view returns (bool valid, uint256 index) {
        return _verifyAsm(proof, root, leaf);
    }

    function _invalidate(uint256 index) private {
        uint256 claimedWordIndex = index >> 8;
        uint256 claimedBitIndex = index & 0xff;
        uint256 claimedWord = _claimedBitMap[claimedWordIndex];
        uint256 newClaimedWord = claimedWord | (1 << claimedBitIndex);
        require(claimedWord != newClaimedWord, "MD: Drop already claimed");
        _claimedBitMap[claimedWordIndex] = newClaimedWord;
    }

    function _verifyAsm(bytes calldata proof, bytes16 root, bytes16 leaf) private view returns (bool valid, uint256 index) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x10)
            let ptr := proof.offset
            let mask := 1

            for { let end := add(ptr, proof.length) } lt(ptr, end) { ptr := add(ptr, 0x10) } {
                let node := calldataload(ptr)

                switch lt(leaf, node)
                case 1 {
                    mstore(mem1, leaf)
                    mstore(mem2, node)
                }
                default {
                    mstore(mem1, node)
                    mstore(mem2, leaf)
                    index := or(mask, index)
                }

                leaf := keccak256(mem1, 32)
                mask := shl(1, mask)
            }

            valid := iszero(shr(128, xor(root, leaf)))
        }
        unchecked {
            index <<= depth - proof.length / 16;
        }
    }
}
