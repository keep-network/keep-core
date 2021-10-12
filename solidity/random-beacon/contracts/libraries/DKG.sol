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
        uint256 dkgResultChallengePeriodLength;
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
        // List of results submitted during DKG.
        RegisteredDkgResult[] registeredDkgResults;
    }

    struct DkgResult {
        // Claimed submitter candidate group member index
        uint256 submitterMemberIndex;
        // Generated candidate group public key
        bytes groupPubKey;
        // Bytes array of misbehaved (disqualified or inactive)
        bytes misbehaved;
        // Concatenation of signatures from members supporting the result.
        bytes signatures;
        // Indices of members corresponding to each signature. Indices have to be unique.
        uint256[] signingMembersIndexes;
        // Addresses of candidate group members as outputted by the group selection protocol.
        address[] members;
    }

    struct RegisteredDkgResult {
        uint256 resultSubmittedTimestamp;
        bytes32 dkgResultHash;
    }

    modifier cleanup(Data storage self) {
        _;
        delete self.seed;
        delete self.groupSize;
        delete self.signatureThreshold;
        delete self.dkgResultSubmissionEligibilityDelay;
        delete self.dkgResultChallengePeriodLength;
        delete self.timeDKG;
        delete self.startBlock;
        delete self.registeredDkgResults;
    }

    function isInProgress(Data storage self) public view returns (bool) {
        return self.startBlock > 0;
    }

    function dkgTimeout(Data storage self) public view returns (uint256) {
        return
            self.timeDKG +
            (self.groupSize * self.dkgResultSubmissionEligibilityDelay) +
            self.dkgResultChallengePeriodLength;
    }

    function start(
        Data storage self,
        uint256 seed,
        uint256 groupSize,
        uint256 signatureThreshold,
        uint256 dkgResultSubmissionEligibilityDelay,
        uint256 dkgResultChallengePeriodLength,
        uint256 timeDKG
    ) internal {
        require(!isInProgress(self), "dkg is currently in progress");

        require(groupSize > 0, "groupSize not set");
        require(signatureThreshold > 0, "signatureThreshold not set");
        require(
            dkgResultSubmissionEligibilityDelay > 0,
            "dkgResultSubmissionEligibilityDelay not set"
        );
        require(
            dkgResultChallengePeriodLength > 0,
            "dkgResultChallengePeriodLength not set"
        );
        require(timeDKG > 0, "timeDKG not set");

        self.seed = seed;
        self.groupSize = groupSize;
        self.signatureThreshold = signatureThreshold;
        self
            .dkgResultSubmissionEligibilityDelay = dkgResultSubmissionEligibilityDelay;
        self.dkgResultChallengePeriodLength = dkgResultChallengePeriodLength;
        self.timeDKG = timeDKG;

        self.startBlock = block.number;
    }

    /// @notice Verifies the submitted DKG result against supporting member
    /// signatures and if the submitter is eligible to submit at the current
    /// block. Every signature supporting the result has to be from a unique
    /// group member.
    /// @param dkgResult DKG result.
    function verify(Data storage self, DkgResult calldata dkgResult)
        public
        view
    {
        require(isInProgress(self), "dkg is currently not in progress");

        assert(self.startBlock > 0);
        assert(self.timeDKG > 0);
        assert(self.dkgResultSubmissionEligibilityDelay > 0);
        assert(self.groupSize > 0);
        assert(self.signatureThreshold > 0);

        require(dkgResult.submitterMemberIndex > 0, "Invalid submitter index");
        require(
            dkgResult.members[dkgResult.submitterMemberIndex - 1] == msg.sender,
            "Sender address doesn't match member with submitter index"
        );

        uint256 T_init = self.startBlock + self.timeDKG;
        require(
            block.number >=
                (T_init +
                    (dkgResult.submitterMemberIndex - 1) *
                    self.dkgResultSubmissionEligibilityDelay),
            "Submitter not eligible"
        );

        require(
            dkgResult.groupPubKey.length == 128,
            "Malformed group public key"
        );

        require(
            dkgResult.misbehaved.length <=
                self.groupSize - self.signatureThreshold,
            "Malformed misbehaved bytes"
        );

        uint256 signaturesCount = dkgResult.signatures.length / 65;
        require(
            dkgResult.signatures.length >= 65,
            "Too short signatures array"
        );
        require(
            dkgResult.signatures.length % 65 == 0,
            "Malformed signatures array"
        );
        require(
            signaturesCount == dkgResult.signingMembersIndexes.length,
            "Unexpected signatures count"
        );
        require(
            signaturesCount >= self.signatureThreshold,
            "Too few signatures"
        );
        require(signaturesCount <= self.groupSize, "Too many signatures");

        bytes32 resultHash = keccak256(
            abi.encodePacked(dkgResult.groupPubKey, dkgResult.misbehaved)
        );

        bytes memory current; // Current signature to be checked.

        bool[] memory usedMemberIndices = new bool[](self.groupSize);

        for (uint256 i = 0; i < signaturesCount; i++) {
            uint256 memberIndex = dkgResult.signingMembersIndexes[i];
            require(memberIndex > 0, "Invalid index");
            require(
                memberIndex <= dkgResult.members.length,
                "Index out of range"
            );

            require(
                !usedMemberIndices[memberIndex - 1],
                "Duplicate member index"
            );
            usedMemberIndices[memberIndex - 1] = true;

            current = dkgResult.signatures.slice(65 * i, 65);
            address recoveredAddress = resultHash
                .toEthSignedMessageHash()
                .recover(current);

            require(
                dkgResult.members[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );
        }
    }

    function notifyTimeout(Data storage self) public cleanup(self) {
        require(isInProgress(self), "dkg is currently not in progress");

        require(
            block.number > self.startBlock + dkgTimeout(self),
            "timeout not passed yet"
        );

        revert("TODO: Implement");
    }

    function submitDkgResult(Data storage self, DkgResult calldata dkgResult)
        external
        returns (uint256)
    {
        require(isInProgress(self), "dkg is currently not in progress");

        verify(self, dkgResult);

        bytes32 dkgResultHash = keccak256(abi.encode(dkgResult));

        self.registeredDkgResults.push(
            RegisteredDkgResult(block.timestamp, dkgResultHash)
        );

        uint256 resultIndex = self.registeredDkgResults.length;

        return resultIndex;
    }

    function challengeResult(
        Data storage self,
        uint256 resultIndex,
        DkgResult calldata dkgResult
    ) external {
        assert(self.dkgResultChallengePeriodLength > 0);

        require(isInProgress(self), "dkg is currently not in progress");

        RegisteredDkgResult memory submittedDkgResult = self
            .registeredDkgResults[resultIndex];

        require(
            block.timestamp <
                submittedDkgResult.resultSubmittedTimestamp +
                    self.dkgResultChallengePeriodLength,
            "Challenge period has already passed"
        );

        bytes32 dkgResultHash = keccak256(abi.encode(dkgResult));

        require(
            dkgResultHash == submittedDkgResult.dkgResultHash,
            "invalid result"
        );

        // TODO: Verify members with sortition pool
    }

    function acceptResult(
        Data storage self,
        uint256 resultIndex,
        DkgResult calldata dkgResult
    ) external cleanup(self) {
        assert(self.dkgResultChallengePeriodLength > 0);

        require(isInProgress(self), "dkg is currently not in progress");

        require(
            resultIndex < self.registeredDkgResults.length,
            "invalid result index"
        );

        RegisteredDkgResult memory registeredDkgResult = self
            .registeredDkgResults[resultIndex];

        require(
            block.timestamp >=
                registeredDkgResult.resultSubmittedTimestamp +
                    self.dkgResultChallengePeriodLength,
            "Challenge period has not passed yet"
        );

        bytes32 dkgResultHash = keccak256(abi.encode(dkgResult));

        require(
            dkgResultHash == registeredDkgResult.dkgResultHash,
            "invalid result"
        );
    }
}
