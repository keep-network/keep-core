// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "./interfaces/ICumulativeMerkleDrop160.sol";


contract CumulativeMerkleDrop160 is Ownable, ICumulativeMerkleDrop160 {
    using SafeERC20 for IERC20;

    address public immutable override token;

    bytes20 public override merkleRoot;
    mapping(address => uint256) public cumulativeClaimed;

    constructor(address token_) {
        token = token_;
    }

    function setMerkleRoot(bytes20 merkleRoot_) external override onlyOwner {
        emit MerkelRootUpdated(merkleRoot, merkleRoot_);
        merkleRoot = merkleRoot_;
    }

    function claim(
        address account,
        uint256 cumulativeAmount,
        bytes20 expectedMerkleRoot,
        bytes calldata merkleProof
    ) external override {
        require(merkleRoot == expectedMerkleRoot, "CMD: Merkle root was updated");

        // Verify the merkle proof
        bytes20 leaf = _keccak160(abi.encodePacked(account, cumulativeAmount));
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

    // function verify(bytes calldata proof, bytes20 root, bytes20 leaf) public pure returns (bool) {
    //     for (uint256 i = 0; i < proof.length / 20; i++) {
    //         bytes20 node = _getBytes20(proof[i*20:(i+1)*20]);
    //         if (leaf < node) {
    //             leaf = _keccak160(abi.encodePacked(leaf, node));
    //         } else {
    //             leaf = _keccak160(abi.encodePacked(node, leaf));
    //         }
    //     }
    //     return leaf == root;
    // }
    //
    // function _getBytes20(bytes calldata input) internal pure returns(bytes20 res) {
    //     // solhint-disable-next-line no-inline-assembly
    //     assembly {
    //         res := calldataload(input.offset)
    //     }
    // }

    // Experimental assembly optimization
    function verifyAsm(bytes calldata proof, bytes20 root, bytes20 leaf) public pure returns (bool valid) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x14)
            let ptr := proof.offset

            for { let end := add(ptr, proof.length) } lt(ptr, end) { ptr := add(ptr, 0x14) } {
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

                leaf := keccak256(mem1, 0x28)
            }

            valid := iszero(shr(96, xor(root, leaf)))
        }
    }

    function _keccak160(bytes memory input) internal pure returns(bytes20) {
        return bytes20(keccak256(input));
    }
}
