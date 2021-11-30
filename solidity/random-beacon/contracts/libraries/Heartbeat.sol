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

pragma solidity ^0.8.6;

import "./BytesLib.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";

library Heartbeat {
    using BytesLib for bytes;
    using ECDSA for bytes32;

    struct FailureClaim {
        // ID of the group raising the claim.
        uint64 groupId;
        // Public key of the group.
        bytes groupPubKey;
        // Identifiers of all group members.
        uint32[] groupMembers;
        // Indices of members accused of failed heartbeat. Indices must be in
        // range <1, groupMembers.length>, unique, and sorted in ascending order.
        uint256[] failedMembersIndices;
        // Concatenation of signatures from members supporting the claim.
        // The message to be signed by each member is failed heartbeat nonce
        // for given group, keccak256 hash of the group public key, and failed
        // members indices. The calculated hash should be prefixed with prefixed
        // with `\x19Ethereum signed message:\n` before signing, so the message
        // to sign is:
        // `\x19Ethereum signed message:\n${keccak256(
        //    nonce, groupPubKey, failedMembersIndices
        // )}`
        bytes signatures;
        // Indices of members corresponding to each signature. Indices must be
        // in range <1, groupMembers.length>, unique, and sorted in ascending
        // order.
        uint256[] signingMembersIndices;
    }

    uint256 public constant signatureByteSize = 65;

    /// @notice Verifies the failure claim according to rules mentioned in
    ///         `FailureClaim` struct documentation. Reverts if verification
    ///         fails.
    /// @param claim Failure claim. Group data passed in the claim must be
    ///        validated by the calling code. This function assumes they are
    ///        all correct.
    /// @param sortitionPool Sortition pool used by the application performing
    ///        claim verification.
    /// @return failedMembers Identifiers of members who failed the heartbeat.
    function verifyFailureClaim(
        FailureClaim calldata claim,
        SortitionPool sortitionPool,
        uint256 nonce
    ) internal view returns (uint32[] memory failedMembers) {
        uint256 groupSize = claim.groupMembers.length;
        // At least half of the members plus one must vote for the claim.
        uint256 groupThreshold = (groupSize / 2) + 1;

        // Validate failed members indices. Maximum indices count is equal to
        // the group size and is not limited deliberately to leave a theoretical
        // possibility to accuse more members than `groupSize - groupThreshold`.
        validateMembersIndices(claim.failedMembersIndices, groupSize);

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
        require(signaturesCount <= groupSize, "Too many signatures");

        // Validate signing members indices. Note that `signingMembersIndices`
        // were already partially validated during signature validation.
        validateMembersIndices(claim.signingMembersIndices, groupSize);

        // Each signing member needs to sign the hash of packed `groupPubKey`
        // and `failedMembersIndices` parameters.
        bytes32 signedMessageHash = keccak256(
            abi.encodePacked(
                nonce,
                claim.groupPubKey,
                claim.failedMembersIndices
            )
        );

        address[] memory groupMembersAddresses = sortitionPool.getIDOperators(
            claim.groupMembers
        );

        // Verify each signature.
        bytes memory checkedSignature;
        for (uint256 i = 0; i < signaturesCount; i++) {
            uint256 memberIndex = claim.signingMembersIndices[i];
            checkedSignature = claim.signatures.slice(
                signatureByteSize * i,
                signatureByteSize
            );
            address recoveredAddress = signedMessageHash
                .toEthSignedMessageHash()
                .recover(checkedSignature);

            require(
                groupMembersAddresses[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );
        }

        failedMembers = new uint32[](claim.failedMembersIndices.length);
        for (uint256 i = 0; i < claim.failedMembersIndices.length; i++) {
            uint256 memberIndex = claim.failedMembersIndices[i];
            failedMembers[i] = claim.groupMembers[memberIndex - 1];
        }

        return failedMembers;
    }

    /// @notice Validates members indices array. Array is considered valid
    ///         if its size and each single index are in <1, groupSize> range,
    ///         indexes are unique, and sorted in an ascending order.
    ///         Reverts if validation fails..
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

        // Check if first and last indices are in range <1, groupSize>.
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
