// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "./interfaces/ICumulativeMerkleDrop128.sol";


contract CumulativeMerkleDrop128 is Ownable, ICumulativeMerkleDrop128 {
    using SafeERC20 for IERC20;

    address public immutable override token;

    bytes16 public override merkleRoot;
    mapping(address => uint256) public cumulativeClaimed;

    constructor(address token_) {
        token = token_;
    }

    function setMerkleRoot(bytes16 merkleRoot_) external override onlyOwner {
        emit MerkelRootUpdated(merkleRoot, merkleRoot_);
        merkleRoot = merkleRoot_;
    }

    function claim(
        address account,
        uint256 cumulativeAmount,
        bytes16 expectedMerkleRoot,
        bytes calldata merkleProof
    ) external override {
        require(merkleRoot == expectedMerkleRoot, "CMD: Merkle root was updated");

        // Verify the merkle proof
        bytes16 leaf = _keccak128(abi.encodePacked(account, cumulativeAmount));
        require(verifyAsm(merkleProof, expectedMerkleRoot, leaf), "CMD: Invalid proof");

        // Mark it claimed
        uint256 preclaimed = cumulativeClaimed[account];
        require(preclaimed < cumulativeAmount, "CMD: Nothing to claim");
        cumulativeClaimed[account] = cumulativeAmount;

        // Send the token
        uint256 amount = cumulativeAmount - preclaimed;
        IERC20(token).safeTransfer(account, amount);
        emit Claimed(account, amount);
    }

    // function verify(bytes calldata proof, bytes16 root, bytes16 leaf) public pure returns (bool) {
    //     for (uint256 i = 0; i < proof.length / 16; i++) {
    //         bytes16 node = _getBytes16(proof[i*16:(i+1)*16]);
    //         if (leaf < node) {
    //             leaf = _keccak128(abi.encodePacked(leaf, node));
    //         } else {
    //             leaf = _keccak128(abi.encodePacked(node, leaf));
    //         }
    //     }
    //     return leaf == root;
    // }
    //
    // function _getBytes16(bytes calldata input) internal pure returns(bytes16 res) {
    //     // solhint-disable-next-line no-inline-assembly
    //     assembly {
    //         res := calldataload(input.offset)
    //     }
    // }

    // Experimental assembly optimization
    function verifyAsm(bytes calldata proof, bytes16 root, bytes16 leaf) public pure returns (bool valid) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x10)
            let ptr := proof.offset

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
                }

                leaf := keccak256(mem1, 32)
            }

            valid := iszero(shr(128, xor(root, leaf)))
        }
    }

    function _keccak128(bytes memory input) internal pure returns(bytes16) {
        return bytes16(keccak256(input));
    }
}
