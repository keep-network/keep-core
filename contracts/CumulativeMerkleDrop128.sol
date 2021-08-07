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
        uint256 index,
        address account,
        uint256 amountToClaim,
        uint256 cumulativeAmount,
        bytes16 targetMerkleRoot,
        bytes calldata merkleProof
    ) external override {
        require(amountToClaim > 0, "CMD: Amount should not be 0");
        require(merkleRoot == targetMerkleRoot, "CMD: Merkle root was updated");

        // Verify the merkle proof
        bytes16 node = _keccak128(abi.encodePacked(account, cumulativeAmount));
        require(targetMerkleRoot == applyProof(index, node, merkleProof), "CMD: Invalid proof");

        // Mark it claimed
        uint256 claimed = cumulativeClaimed[account] + amountToClaim;
        cumulativeClaimed[account] = claimed;
        if (claimed - amountToClaim == cumulativeAmount) {
            revert("CMD: Drop already claimed");
        }
        else if (claimed > cumulativeAmount) {
            revert("CMD: Claiming amount is too high");
        }

        // Send the token
        IERC20(token).safeTransfer(account, amountToClaim);
        emit Claimed(index, account, amountToClaim);
    }

    function applyProof(uint256 index, bytes16 leaf, bytes calldata proof) public view returns (bytes16 computedHash) {
        this; // hide pure warning
        computedHash = leaf;

        for (uint256 i = 0; i < proof.length / 16; i++) {
            if ((index >> i) & 1 == 0) {
                computedHash = _keccak128(abi.encodePacked(computedHash, proof[i*16:(i+1)*16]));
            } else {
                computedHash = _keccak128(abi.encodePacked(proof[i*16:(i+1)*16], computedHash));
            }
        }
    }

    // Experimental assembly optimization
    function applyProof2(uint256 index, bytes16 leaf, bytes calldata proof) public pure returns (bytes16) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            leaf := shr(128, leaf)
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x10)
            let len := div(proof.length, 0x10)
            let ptr := add(proof.offset, 0x20)
            for { let i := 0 } lt(i, len) { i := add(i, 1) } {
                switch and(shr(i, index), 1)
                case 0 {
                    mstore(mem2, calldataload(ptr))
                    mstore(mem1, leaf)
                }
                default {
                    mstore(mem2, leaf)
                    mstore(mem1, calldataload(ptr))
                }

                ptr := add(ptr, 0x10)
                leaf := keccak256(mem2, 32)
            }

            mstore(mem1, leaf)
            return(mem2, 32)
        }
    }

    function _keccak128(bytes memory input) internal pure returns(bytes16) {
        return bytes16(bytes32(uint256(keccak256(input)) << 128));
    }
}
