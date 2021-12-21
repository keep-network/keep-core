// SPDX-License-Identifier: MIT
//
// ▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
//   ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
// ▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
//
//                           Trust math, not hardware.

pragma solidity ^0.8.9;

import "./BytesLib.sol";
import "./Groups.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";

library Heartbeat {
    using BytesLib for bytes;
    using ECDSA for bytes32;

    struct FailureClaim {
        // ID of the group raising the claim.
        uint64 groupId;
        // Indices of members accused of failed heartbeat. Indices must be in
        // range [1, groupMembers.length], unique, and sorted in ascending order.
        uint256[] failedMembersIndices;
        // Concatenation of signatures from members supporting the claim.
        // The message to be signed by each member is failed heartbeat nonce
        // for given group, keccak256 hash of the group public key, and failed
        // members indices. The calculated hash should be prefixed with
        // `\x19Ethereum signed message:\n` before signing, so the message
        // to sign is:
        // `\x19Ethereum signed message:\n${keccak256(
        //    nonce, groupPubKey, failedMembersIndices
        // )}`
        bytes signatures;
        // Indices of members corresponding to each signature. Indices must be
        // in range [1, groupMembers.length], unique, and sorted in ascending
        // order.
        uint256[] signingMembersIndices;
    }

    /// @notice The minimum number of group members needed to interact according
    ///         to the protocol to produce a valid failure claim.
    uint256 public constant groupThreshold = 33;

    /// @notice Size in bytes of a single signature produced by member
    ///         supporting the failure claim.
    uint256 public constant signatureByteSize = 65;

    /// @notice Verifies the failure claim according to rules mentioned in
    ///         `FailureClaim` struct documentation. Reverts if verification
    ///         fails.
    /// @param sortitionPool Sortition pool used by the application performing
    ///        claim verification.
    /// @param claim Failure claim.
    /// @param group Group raising the claim.
    /// @param nonce Current nonce for group used in the claim.
    /// @return failedMembers Identifiers of members who failed the heartbeat.
    function verifyFailureClaim(
        SortitionPool sortitionPool,
        FailureClaim calldata claim,
        Groups.Group storage group,
        uint256 nonce,
        uint32[] calldata members
    ) external view returns (uint32[] memory failedMembers) {
        // Validate failed members indices. Maximum indices count is equal to
        // the group size and is not limited deliberately to leave a theoretical
        // possibility to accuse more members than `groupSize - groupThreshold`.
        validateMembersIndices(claim.failedMembersIndices, members.length);

        // Validate signatures array is properly formed and number of
        // signatures and signers is correct.
        uint256 signaturesCount = claim.signatures.length / signatureByteSize;
        require(claim.signatures.length != 0, "No signatures provided");
        require(
            claim.signatures.length % signatureByteSize == 0,
            "Malformed signatures array"
        );
        require(
            signaturesCount == claim.signingMembersIndices.length,
            "Unexpected signatures count"
        );
        require(signaturesCount >= groupThreshold, "Too few signatures");
        require(signaturesCount <= members.length, "Too many signatures");

        // Validate signing members indices. Note that `signingMembersIndices`
        // were already partially validated during `signatures` parameter
        // validation.
        validateMembersIndices(claim.signingMembersIndices, members.length);

        // Each signing member needs to sign the hash of packed `groupPubKey`
        // and `failedMembersIndices` parameters. Usage of group public key
        // and not group ID is important because it provides uniqueness of
        // signed messages and prevent against reusing them in future.
        bytes32 signedMessageHash = keccak256(
            abi.encodePacked(
                nonce,
                group.groupPubKey,
                claim.failedMembersIndices
            )
        ).toEthSignedMessageHash();

        address[] memory groupMembers = sortitionPool.getIDOperators(members);

        // Verify each signature.
        bytes memory checkedSignature;
        bool senderSignatureExists = false;
        for (uint256 i = 0; i < signaturesCount; i++) {
            uint256 memberIndex = claim.signingMembersIndices[i];
            checkedSignature = claim.signatures.slice(
                signatureByteSize * i,
                signatureByteSize
            );
            address recoveredAddress = signedMessageHash.recover(
                checkedSignature
            );

            require(
                groupMembers[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );

            if (!senderSignatureExists && msg.sender == recoveredAddress) {
                senderSignatureExists = true;
            }
        }

        require(senderSignatureExists, "Sender must be claim signer");

        failedMembers = new uint32[](claim.failedMembersIndices.length);
        for (uint256 i = 0; i < claim.failedMembersIndices.length; i++) {
            uint256 memberIndex = claim.failedMembersIndices[i];
            failedMembers[i] = members[memberIndex - 1];
        }

        return failedMembers;
    }

    /// @notice Validates members indices array. Array is considered valid
    ///         if its size and each single index are in [1, groupSize] range,
    ///         indexes are unique, and sorted in an ascending order.
    ///         Reverts if validation fails.
    /// @param indices Array to validate.
    /// @param groupSize Group size used as reference.
    function validateMembersIndices(
        uint256[] calldata indices,
        uint256 groupSize
    ) internal view {
        require(
            indices.length > 0 && indices.length <= groupSize,
            "Corrupted members indices"
        );

        // Check if first and last indices are in range [1, groupSize].
        // This check combined with the loop below makes sure every single
        // index is in the correct range.
        require(
            indices[0] > 0 && indices[indices.length - 1] <= groupSize,
            "Corrupted members indices"
        );

        for (uint256 i = 0; i < indices.length - 1; i++) {
            // Check whether given index is smaller than the next one. This
            // way we are sure indexes are ordered in the ascending order
            // and there are no duplicates.
            require(indices[i] < indices[i + 1], "Corrupted members indices");
        }
    }
}
