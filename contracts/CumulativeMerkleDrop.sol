// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import "./interfaces/ICumulativeMerkleDrop.sol";


contract CumulativeMerkleDrop is Ownable, ICumulativeMerkleDrop {
    using SafeERC20 for IERC20;
    using MerkleProof for bytes32[];

    address public immutable override token;
    address public rewardsHolder;

    bytes32 public override merkleRoot;
    mapping(address => uint256) public cumulativeClaimed;
    struct Claim {
        address account;
        uint256 amount;
        bytes32[] proof;
    }

    event RewardsHolderUpdated(address oldRewardsHolder, address newRewardsHolder);

    constructor(address token_, address rewardsHolder_) {
        require(IERC20(token_).totalSupply() > 0, "Token contract must be set");
        require(rewardsHolder_ != address(0), "Rewards GHolder must be an address");
        token = token_;
        rewardsHolder = rewardsHolder_;
    }

    function setMerkleRoot(bytes32 merkleRoot_) external override onlyOwner {
        emit MerkelRootUpdated(merkleRoot, merkleRoot_);
        merkleRoot = merkleRoot_;
    }

    function setRewardsHolder(address rewardsHolder_) external onlyOwner {
        require(rewardsHolder_ != address(0));
        emit RewardsHolderUpdated(rewardsHolder, rewardsHolder_);
        rewardsHolder = rewardsHolder_;
    }

    function claim(
        address account,
        uint256 cumulativeAmount,
        bytes32 expectedMerkleRoot,
        bytes32[] calldata merkleProof
    ) public override {
        require(merkleRoot == expectedMerkleRoot, "CMD: Merkle root was updated");

        // Verify the merkle proof
        bytes32 leaf = keccak256(abi.encodePacked(account, cumulativeAmount));
        require(_verifyAsm(merkleProof, expectedMerkleRoot, leaf), "CMD: Invalid proof");

        // Mark it claimed
        uint256 preclaimed = cumulativeClaimed[account];
        require(preclaimed < cumulativeAmount, "CMD: Nothing to claim");
        cumulativeClaimed[account] = cumulativeAmount;

        // Send the token
        unchecked {
            uint256 amount = cumulativeAmount - preclaimed;
            IERC20(token).safeTransferFrom(rewardsHolder, account, amount);
            emit Claimed(account, amount);
        }
    }

    function batchClaim(
        bytes32 expectedMerkleRoot,
        Claim[] calldata Claims
    ) external {
        for (uint i; i < Claims.length; i++) {
            claim(
                Claims[i].account,
                Claims[i].amount,
                expectedMerkleRoot,
                Claims[i].proof
            );
        }
     }

    function verify(bytes32[] calldata merkleProof, bytes32 root, bytes32 leaf) public pure returns (bool) {
        return merkleProof.verify(root, leaf);
    }

    function _verifyAsm(bytes32[] calldata proof, bytes32 root, bytes32 leaf) private pure returns (bool valid) {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let mem1 := mload(0x40)
            let mem2 := add(mem1, 0x20)
            let ptr := proof.offset

            for { let end := add(ptr, mul(0x20, proof.length)) } lt(ptr, end) { ptr := add(ptr, 0x20) } {
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

                leaf := keccak256(mem1, 0x40)
            }

            valid := eq(root, leaf)
        }
    }
}
