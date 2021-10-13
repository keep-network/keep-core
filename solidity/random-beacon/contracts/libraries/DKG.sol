// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "./BytesLib.sol";

library DKG {
  using BytesLib for bytes;
  using ECDSA for bytes32;

  /// @dev Size of a group in the threshold relay.
  uint256 public constant groupSize = 64;

  /// @dev Minimum number of group members needed to interact according to the
  /// protocol to produce a relay entry.
  uint256 public constant groupThreshold = 33;

  /// @dev The minimum number of signatures required to support DKG result.
  /// This number needs to be at least the same as the signing threshold
  /// and it is recommended to make it higher than the signing threshold
  /// to keep a safety margin for misbehaving members.
  uint256 public constant signatureThreshold =
    groupThreshold + (groupSize - groupThreshold) / 2;

  /// @notice Time in blocks after which DKG result is complete and ready to be
  // published by clients.
  uint256 public constant offchainDkgTime = 5 * (1 + 5) + 2 * (1 + 10) + 20;

  struct Data {
    // Time in blocks after which the next group member is eligible
    // to submit DKG result.
    uint256 resultSubmissionEligibilityDelay;
    // Time in blocks at which DKG started.
    uint256 startBlock;
    // Hash of submitted DKG result.
    bytes32 submittedResultHash;
    // Block number from the moment of the DKG result submission.
    uint256 submittedResultBlock;
  }

  /// @notice DKG result.
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
    uint256[] signingMemberIndices;
    // Addresses of candidate group members as outputted by the group selection protocol.
    address[] members;
  }

  function start(Data storage self, uint256 resultSubmissionEligibilityDelay)
    internal
  {
    require(!isInProgress(self), "dkg is currently in progress");

    require(
      resultSubmissionEligibilityDelay > 0,
      "resultSubmissionEligibilityDelay not set"
    );

    self.resultSubmissionEligibilityDelay = resultSubmissionEligibilityDelay;

    self.startBlock = block.number;
  }

  function submitDkgResult(Data storage self, DkgResult calldata dkgResult)
    external
  {
    require(isInProgress(self), "dkg is currently not in progress");

    bytes32 dkgResultHash = keccak256(abi.encode(dkgResult));

    require(
      self.submittedResultHash == 0,
      "result was already submitted in the current dkg"
    );

    assert(self.startBlock > 0);
    assert(self.resultSubmissionEligibilityDelay > 0);

    verify(
      self,
      dkgResult.submitterMemberIndex,
      dkgResult.groupPubKey,
      dkgResult.misbehaved,
      dkgResult.signatures,
      dkgResult.signingMemberIndices,
      dkgResult.members,
      self.startBlock
    );

    self.submittedResultHash = dkgResultHash;
    self.submittedResultBlock = block.number;
  }

  /// @notice Checks if DKG is currently in progress.
  /// @return True if DKG is in progress, false otherwise.
  function isInProgress(Data storage self) public view returns (bool) {
    return self.startBlock > 0;
  }

  /// @notice Verifies the submitted DKG result against supporting member
  /// signatures and if the submitter is eligible to submit at the current
  /// block. Every signature supporting the result has to be from a unique
  /// group member.
  ///
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
  /// @param groupSelectionEndBlock Block height at which the group selection
  /// protocol ended.
  function verify(
    Data storage self,
    uint256 submitterMemberIndex,
    bytes memory groupPubKey,
    bytes memory misbehaved,
    bytes memory signatures,
    uint256[] memory signingMemberIndices,
    address[] memory members,
    uint256 groupSelectionEndBlock
  ) public view {
    require(submitterMemberIndex > 0, "Invalid submitter index");
    require(
      members[submitterMemberIndex - 1] == msg.sender,
      "Unexpected submitter index"
    );

    uint256 T_init = groupSelectionEndBlock + offchainDkgTime;
    require(
      block.number >=
        (T_init +
          (submitterMemberIndex - 1) *
          self.resultSubmissionEligibilityDelay),
      "Submitter not eligible"
    );

    require(groupPubKey.length == 128, "Malformed group public key");

    require(
      misbehaved.length <= groupSize - signatureThreshold,
      "Malformed misbehaved bytes"
    );

    uint256 signaturesCount = signatures.length / 65;
    require(signatures.length >= 65, "Too short signatures array");
    require(signatures.length % 65 == 0, "Malformed signatures array");
    require(
      signaturesCount == signingMemberIndices.length,
      "Unexpected signatures count"
    );
    require(signaturesCount >= signatureThreshold, "Too few signatures");
    require(signaturesCount <= groupSize, "Too many signatures");

    bytes32 resultHash = keccak256(abi.encodePacked(groupPubKey, misbehaved));

    bytes memory current; // Current signature to be checked.

    bool[] memory usedMemberIndices = new bool[](groupSize);

    for (uint256 i = 0; i < signaturesCount; i++) {
      uint256 memberIndex = signingMemberIndices[i];
      require(memberIndex > 0, "Invalid index");
      require(memberIndex <= members.length, "Index out of range");

      require(!usedMemberIndices[memberIndex - 1], "Duplicate member index");
      usedMemberIndices[memberIndex - 1] = true;

      current = signatures.slice(65 * i, 65);
      address recoveredAddress = resultHash.toEthSignedMessageHash().recover(
        current
      );
      require(
        members[memberIndex - 1] == recoveredAddress,
        "Invalid signature"
      );
    }
  }
}
