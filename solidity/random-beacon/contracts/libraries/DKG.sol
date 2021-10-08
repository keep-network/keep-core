// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "../utils/BytesLib.sol";

// TODO: Consider which functions can be internal or internal. What are implications
// to security.

/// @title DKG library
library DKG {
    using BytesLib for bytes;
    using ECDSA for bytes32;

    struct Data {
        uint256 seed;
        // Size of a group in the threshold relay.
        uint256 groupSize;
        // Time in blocks after which the next group member is eligible
        // to submit DKG result.
        uint256 dkgResultSubmissionEligibilityDelay;
        // Time in blocks at which DKG started.
        uint256 startBlock;
        // Time in blocks after which DKG result is complete and ready to be
        // published by clients.
        uint256 timeDKG;
        // The minimum number of signatures required to support DKG result.
        // This number needs to be at least the same as the signing threshold
        // and it is recommended to make it higher than the signing threshold
        // to keep a safety margin for misbehaving members.
        uint256 signatureThreshold;
    }

    /// DKG is in an invalid state. Expected in progress: `expectedInProgress`,
    /// but actual in progress: `actualInProgress`.
    error InvalidInProgressState(
        bool expectedInProgress,
        bool actualInProgress
    );

    modifier assertInProgress(Data storage self, bool expectedValue) {
        if (isInProgress(self) != expectedValue)
            revert InvalidInProgressState(expectedValue, isInProgress(self));
        _;
    }

    modifier cleanup(Data storage self) {
        _;
        delete self.seed;
        delete self.groupSize;
        delete self.signatureThreshold;
        delete self.dkgResultSubmissionEligibilityDelay;
        delete self.timeDKG;
        delete self.startBlock;
    }

    function isInProgress(Data storage self) public view returns (bool) {
        return self.startBlock > 0;
    }

    function dkgTimeout(Data storage self) public view returns (uint256) {
        return
            self.timeDKG + self.groupSize * self.dkgResultSubmissionEligibilityDelay;
    }

    function start(
        Data storage self,
        uint256 seed,
        uint256 groupSize,
        uint256 signatureThreshold,
        uint256 dkgResultSubmissionEligibilityDelay,
        uint256 timeDKG
    ) internal assertInProgress(self, false) {
        self.seed = seed;
        self.groupSize = groupSize;
        self.signatureThreshold = signatureThreshold;
        self.dkgResultSubmissionEligibilityDelay = dkgResultSubmissionEligibilityDelay;
        self.timeDKG = timeDKG;

        self.startBlock = block.number;
    }

    /// @notice Verifies the submitted DKG result against supporting member
    /// signatures and if the submitter is eligible to submit at the current
    /// block. Every signature supporting the result has to be from a unique
    /// group member.
    /// @param submitterMemberIndex Claimed submitter candidate group member index
    /// @param groupPubKey Generated candidate group public key
    /// @param misbehaved Bytes array of misbehaved (disqualified or inactive)
    /// group members indexes; Indexes reflect positions of members in the group,
    /// as outputted by the group selection protocol.
    /// @param signatures Concatenation of signatures from members supporting the
    /// result.
    /// @param signingMemberIndices Indices of members corresponding to each
    /// signature. Indices have to be unique.
    /// @param members Addresses of candidate group members as outputted by the
    /// group selection protocol.
    function verify(
        Data storage self,
        uint256 submitterMemberIndex,
        bytes memory groupPubKey,
        bytes memory misbehaved,
        bytes memory signatures,
        uint256[] memory signingMemberIndices,
        address[] memory members
    ) public view assertInProgress(self, true) {
        require(submitterMemberIndex > 0, "Invalid submitter index");
        require(
            members[submitterMemberIndex - 1] == msg.sender,
            "Unexpected submitter index"
        );

        // TODO: Revisit `timeDKG` value if it's something we need and can read
        // from governable parameters.
        uint256 T_init = self.startBlock + self.timeDKG;
        require(
            block.number >=
                (T_init +
                    (submitterMemberIndex - 1) *
                    self.dkgResultSubmissionEligibilityDelay),
            "Submitter not eligible"
        );

        require(groupPubKey.length == 128, "Malformed group public key");

        require(
            misbehaved.length <= self.groupSize - self.signatureThreshold,
            "Malformed misbehaved bytes"
        );

        uint256 signaturesCount = signatures.length / 65;
        require(signatures.length >= 65, "Too short signatures array");
        require(signatures.length % 65 == 0, "Malformed signatures array");
        require(
            signaturesCount == signingMemberIndices.length,
            "Unexpected signatures count"
        );
        require(
            signaturesCount >= self.signatureThreshold,
            "Too few signatures"
        );
        require(signaturesCount <= self.groupSize, "Too many signatures");

        bytes32 resultHash = keccak256(
            abi.encodePacked(groupPubKey, misbehaved)
        );

        bytes memory current; // Current signature to be checked.

        bool[] memory usedMemberIndices = new bool[](self.groupSize);

        for (uint256 i = 0; i < signaturesCount; i++) {
            uint256 memberIndex = signingMemberIndices[i];
            require(memberIndex > 0, "Invalid index");
            require(memberIndex <= members.length, "Index out of range");

            require(
                !usedMemberIndices[memberIndex - 1],
                "Duplicate member index"
            );
            usedMemberIndices[memberIndex - 1] = true;

            current = signatures.slice(65 * i, 65);
            address recoveredAddress = resultHash
                .toEthSignedMessageHash()
                .recover(current);

            require(
                members[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );
        }
    }

    function notifyTimeout(Data storage self)
        internal
        assertInProgress(self, true)
        cleanup(self)
    {
        if (block.number <= self.startBlock + dkgTimeout(self))
            revert NotTimedOut(
                self.startBlock + dkgTimeout(self) + 1,
                block.number
            );
    }

    function finish(Data storage self)
        internal
        assertInProgress(self, true)
        cleanup(self)
    {}
}
