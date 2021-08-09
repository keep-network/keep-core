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
        uint256 index,
        address account,
        uint256 amountToClaim,
        uint256 cumulativeAmount,
        bytes20 targetMerkleRoot,
        bytes calldata merkleProof
    ) external override {
        require(amountToClaim > 0, "CMD: Amount should not be 0");
        require(merkleRoot == targetMerkleRoot, "CMD: Merkle root was updated");

        // Verify the merkle proof
        bytes20 node = _keccak160(abi.encodePacked(account, cumulativeAmount));
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

    function applyProof(uint256 index, bytes20 leaf, bytes calldata proof) public view returns (bytes20 computedHash) {
        this; // hide pure warning
        computedHash = leaf;

        for (uint256 i = 0; i < proof.length / 20; i++) {
            if ((index >> i) & 1 == 0) {
                computedHash = _keccak160(abi.encodePacked(computedHash, proof[i*20:(i+1)*20]));
            } else {
                computedHash = _keccak160(abi.encodePacked(proof[i*20:(i+1)*20], computedHash));
            }
        }
    }

    // Experimental assembly optimization
    function applyProof2(uint256 index, bytes20 leaf, bytes calldata proof) public pure returns (bytes20) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x14)
            let len := div(proof.length, 0x14)
            let ptr := proof.offset
            for { let i := 0 } lt(i, len) { i := add(i, 1) } {
                switch and(shr(i, index), 1)
                case 0 {
                    mstore(mem1, leaf)
                    mstore(mem2, calldataload(ptr))
                }
                default {
                    mstore(mem1, calldataload(ptr))
                    mstore(mem2, leaf)
                }

                ptr := add(ptr, 0x14)
                leaf := shl(96, keccak256(mem1, 40))
            }

            mstore(mem1, leaf)
            return(mem1, 32)
        }
    }

    function _keccak160(bytes memory input) internal pure returns(bytes20) {
        return bytes20(keccak256(input) << 96);
    }
}
