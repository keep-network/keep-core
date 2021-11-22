// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "./BytesLib.sol";
import {ISortitionPool} from "../RandomBeacon.sol";

library DKG {
    using BytesLib for bytes;
    using ECDSA for bytes32;

    struct Parameters {
        // Time in blocks during which a submitted result can be challenged.
        uint256 resultChallengePeriodLength;
        // Time in blocks after which the next group member is eligible
        // to submit DKG result.
        uint256 resultSubmissionEligibilityDelay;
    }

    struct Data {
        // Address of the Sortition Pool contract.
        ISortitionPool sortitionPool;
        // DKG parameters. The parameters should persist between DKG executions.
        // They should be updated with dedicated set functions only when DKG is not
        // in progress.
        Parameters parameters;
        // Time in blocks at which DKG started.
        uint256 startBlock;
        // Seed used to start DKG.
        uint256 seed;
        // Time in blocks that should be added to result submission eligibility
        // delay calculation. It is used in case of a challenge to adjust
        // block calculation for members submission eligibility.
        uint256 resultSubmissionStartBlockOffset;
        // Hash of submitted DKG result.
        bytes32 submittedResultHash;
        // Hash of group members array submitted in the result.
        bytes32 submittedResultGroupMembersHash;
        // Identifiers of group members who signed the submitted result.
        uint32[] submittedResultSigningMembers;
        // Block number from the moment of the DKG result submission.
        uint256 submittedResultBlock;
        // Misbehaved (inactive or disqualified) members from the DKG result.
        uint32[] submittedResultMisbehavedMembers;
        // Address of the DKG result submitter
        address resultSubmitter;
    }

    /// @notice DKG result.
    struct Result {
        // Claimed submitter candidate group member index.
        // Must be in range [1, 64].
        uint256 submitterMemberIndex;
        // Generated candidate group public key
        bytes groupPubKey;
        // Array of misbehaved members indices (disqualified or inactive).
        // Must be in range [1, 64] and unique.
        uint8[] misbehavedMembersIndices;
        // Concatenation of signatures from members supporting the result.
        // The message to be signed by each member is keccak256 hash of the
        // calculated group public key, misbehaved members as bytes and DKG
        // start block. The calculated hash should be prefixed with prefixed with
        // `\x19Ethereum signed message:\n` before signing, so the message to
        // sign is:
        // `\x19Ethereum signed message:\n${keccak256(groupPubKey,misbehaved,startBlock)}`
        bytes signatures;
        // Indices of members corresponding to each signature. Must be in
        // range [1, 64] and unique.
        uint256[] signingMembersIndices;
        // Identifiers of candidate group members as outputted by the group
        // selection protocol.
        uint32[] members;
    }

    /// @notice States for phases of group creation. The states doesn't include
    ///         timeouts which should be tracked and notified individually.
    enum State {
        // Group creation is not in progress. It is a state set after group creation
        // completion either by timeout or by a result approval.
        IDLE,
        // Group creation is awaiting the seed and sortition pool is locked.
        AWAITING_SEED,
        // Off-chain DKG protocol execution is in progress. A result is being calculated
        // by the clients in this state. It's not yet possible to submit the result.
        KEY_GENERATION,
        // After off-chain DKG protocol execution the contract awaits result submission.
        // This is a state to which group creation returns in case of a result
        // challenge notification.
        AWAITING_RESULT,
        // DKG result was submitted and awaits an approval or a challenge. If a result
        // gets challenge the state returns to `AWAITING_RESULT`. If a result gets
        // approval the state changes to `IDLE`.
        CHALLENGE
    }

    /// @dev Size of a group in the threshold relay.
    uint256 public constant groupSize = 64;

    /// @dev Minimum number of group members needed to interact according to the
    ///      protocol to produce a relay entry.
    uint256 public constant groupThreshold = 33;

    /// @dev The minimum number of signatures required to support DKG result.
    ///      This number needs to be at least the same as the signing threshold
    ///      and it is recommended to make it higher than the signing threshold
    ///      to keep a safety margin for in case some members become inactive.
    uint256 public constant signatureThreshold =
        groupThreshold + (groupSize - groupThreshold) / 2;

    /// @notice Time in blocks after which DKG result is complete and ready to be
    //          published by clients.
    uint256 public constant offchainDkgTime = 5 * (1 + 5) + 2 * (1 + 10) + 20;

    event DkgStarted(uint256 indexed seed);

    // TODO: Revisit properties returned in this event when working on result
    // challenges and the client.
    //  TODO: Should it also return seed to link the result with a DKG run?
    event DkgResultSubmitted(
        bytes32 indexed resultHash,
        bytes indexed groupPubKey,
        address indexed submitter
    );

    event DkgTimedOut();

    event DkgResultApproved(
        bytes32 indexed resultHash,
        address indexed approver
    );

    event DkgResultChallenged(
        bytes32 indexed resultHash,
        address indexed challenger
    );

    event DkgStateLocked();

    event DkgSeedTimedOut();

    /// @notice Initializes the sortitionPool parameter. Can be performed only once.
    /// @param _sortitionPool Value of the parameter.
    function initSortitionPool(Data storage self, ISortitionPool _sortitionPool)
        internal
    {
        require(
            address(self.sortitionPool) == address(0),
            "Sortition Pool address already set"
        );

        self.sortitionPool = _sortitionPool;
    }

    /// @notice Determines the current state of group creation. It doesn't take
    ///         timeouts into consideration. The timeouts should be tracked and
    ///         notified separately.
    function currentState(Data storage self)
        internal
        view
        returns (State state)
    {
        state = State.IDLE;

        if (self.sortitionPool.isLocked()) {
            state = State.AWAITING_SEED;

            if (self.startBlock > 0) {
                state = State.KEY_GENERATION;

                if (block.number > self.startBlock + offchainDkgTime) {
                    state = State.AWAITING_RESULT;

                    if (self.submittedResultBlock > 0) {
                        state = State.CHALLENGE;
                    }
                }
            }
        }
    }

    /// @notice Locks the sortition pool and starts awaiting for the
    ///         group creation seed.
    function lockState(Data storage self) internal {
        require(currentState(self) == State.IDLE, "current state is not IDLE");

        emit DkgStateLocked();

        self.sortitionPool.lock();
    }

    function start(Data storage self, uint256 seed) internal {
        require(
            currentState(self) == State.AWAITING_SEED,
            "current state is not AWAITING_SEED"
        );

        self.startBlock = block.number;
        self.seed = seed;

        // slither-disable-next-line reentrancy-events
        emit DkgStarted(seed);
    }

    function submitResult(Data storage self, Result calldata result) internal {
        require(
            currentState(self) == State.AWAITING_RESULT,
            "current state is not AWAITING_RESULT"
        );
        require(!hasDkgTimedOut(self), "dkg timeout already passed");

        uint32[] memory signingMembers = verify(self, result);

        // TODO: Check with sortition pool that all members have minimum stake.
        // Check all members in one call or at least members that signed the result.
        // We need to know in advance that there will be something that we can
        // slash the members from.

        for (uint256 i = 0; i < result.misbehavedMembersIndices.length; i++) {
            // group member indices start from 1, so we need to -1 on misbehaved
            uint32 memberArrayPosition = result.misbehavedMembersIndices[i] - 1;
            self.submittedResultMisbehavedMembers.push(
                result.members[memberArrayPosition]
            );
        }
        self.submittedResultHash = keccak256(abi.encode(result));
        self.submittedResultGroupMembersHash = keccak256(
            abi.encodePacked(result.members)
        );
        self.submittedResultSigningMembers = signingMembers;
        self.submittedResultBlock = block.number;
        self.resultSubmitter = msg.sender;

        emit DkgResultSubmitted(
            self.submittedResultHash,
            result.groupPubKey,
            msg.sender
        );
    }

    /// @notice Checks if DKG timed out. The DKG timeout period includes time required
    ///         for off-chain protocol execution and time for the result publication
    ///         for all group members. After this time result cannot be submitted
    ///         and DKG can be notified about the timeout. DKG period is adjusted
    ///         by result submission offset that include blocks that were mined
    ///         while invalid result has been registered until it got challenged.
    /// @return True if DKG timed out, false otherwise.
    function hasDkgTimedOut(Data storage self) internal view returns (bool) {
        return
            currentState(self) == State.AWAITING_RESULT &&
            block.number >
            (self.startBlock +
                offchainDkgTime +
                self.resultSubmissionStartBlockOffset +
                groupSize *
                self.parameters.resultSubmissionEligibilityDelay);
    }

    /// @notice Verifies the submitted DKG result against supporting member
    ///         signatures and if the submitter is eligible to submit at the current
    ///         block. Every signature supporting the result has to be from a unique
    ///         group member.
    /// @dev The message to be signed by each member is keccak256 hash of the
    ///      calculated group public key, misbehaved members as bytes and DKG
    ///      start block. The calculated hash should be prefixed with prefixed with
    ///      `\x19Ethereum signed message:\n` before signing, so the message to
    ///      sign is:
    ///      `\x19Ethereum signed message:\n${keccak256(groupPubKey,misbehaved,startBlock)}`
    ///      Members indexing in the group starts with 1.
    /// @param result DKG result which will be verified.
    /// @return signingMembers Identifiers of group members whose signatures of
    ///         the DKG result hash were proven to be valid.
    function verify(Data storage self, Result calldata result)
        internal
        view
        returns (uint32[] memory signingMembers)
    {
        // TODO: Verify if submitter is valid staker and signatures come from valid
        // stakers https://github.com/keep-network/keep-core/pull/2654#discussion_r728226906.

        require(result.submitterMemberIndex > 0, "Invalid submitter index");

        ISortitionPool sortitionPool = self.sortitionPool;

        require(
            sortitionPool.getIDOperator(
                result.members[result.submitterMemberIndex - 1]
            ) == msg.sender,
            "Unexpected submitter index"
        );

        uint256 T_init = self.startBlock +
            offchainDkgTime +
            self.resultSubmissionStartBlockOffset;
        require(
            block.number >=
                (T_init +
                    (result.submitterMemberIndex - 1) *
                    self.parameters.resultSubmissionEligibilityDelay),
            "Submitter not eligible"
        );

        require(result.groupPubKey.length == 128, "Malformed group public key");

        require(
            result.misbehavedMembersIndices.length <=
                groupSize - signatureThreshold,
            "Unexpected misbehaved members count"
        );

        if (result.misbehavedMembersIndices.length > 1) {
            for (
                uint256 i = 1;
                i < result.misbehavedMembersIndices.length;
                i++
            ) {
                require(
                    result.misbehavedMembersIndices[i - 1] <
                        result.misbehavedMembersIndices[i],
                    "Corrupted misbehaved members indices"
                );
            }
        }

        uint256 signaturesCount = result.signatures.length / 65;
        require(result.signatures.length >= 65, "Too short signatures array");
        require(
            result.signatures.length % 65 == 0,
            "Malformed signatures array"
        );
        require(
            signaturesCount == result.signingMembersIndices.length,
            "Unexpected signatures count"
        );
        require(signaturesCount >= signatureThreshold, "Too few signatures");
        require(signaturesCount <= groupSize, "Too many signatures");

        bytes32 resultHash = keccak256(
            abi.encodePacked(
                result.groupPubKey,
                result.misbehavedMembersIndices,
                self.startBlock
            )
        );

        bytes memory current; // Current signature to be checked.
        bool[] memory usedMemberIndices = new bool[](groupSize);
        signingMembers = new uint32[](signaturesCount);
        address[] memory membersAddresses = sortitionPool.getIDOperators(
            result.members
        );

        for (uint256 i = 0; i < signaturesCount; i++) {
            uint256 memberIndex = result.signingMembersIndices[i];
            require(memberIndex > 0, "Invalid index");
            require(memberIndex <= result.members.length, "Index out of range");

            require(
                !usedMemberIndices[memberIndex - 1],
                "Duplicate member index"
            );
            usedMemberIndices[memberIndex - 1] = true;

            current = result.signatures.slice(65 * i, 65);
            address recoveredAddress = resultHash
                .toEthSignedMessageHash()
                .recover(current);

            signingMembers[i] = result.members[memberIndex - 1];

            require(
                membersAddresses[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );
        }

        return signingMembers;
    }

    /// @notice Notifies about DKG timeout.
    function notifyTimeout(Data storage self) internal {
        require(hasDkgTimedOut(self), "dkg has not timed out");

        emit DkgTimedOut();
    }

    /// @notice Notifies about the seed was not delivered and restores the
    ///         initial DKG state (IDLE).
    function notifySeedTimedOut(Data storage self) internal {
        require(
            currentState(self) == State.AWAITING_SEED,
            "current state is not AWAITING_SEED"
        );

        emit DkgSeedTimedOut();

        self.sortitionPool.unlock();
    }

    /// @notice Approves DKG result. Can be called after challenge period for the
    ///         submitted result is finished. Considers the submitted result as
    ///         valid and completes the group creation.
    function approveResult(Data storage self) internal {
        require(
            currentState(self) == State.CHALLENGE,
            "current state is not CHALLENGE"
        );

        require(
            block.number >
                self.submittedResultBlock +
                    self.parameters.resultChallengePeriodLength,
            "challenge period has not passed yet"
        );

        emit DkgResultApproved(self.submittedResultHash, msg.sender);
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    /// @dev Can be called during a challenge period for the submitted result.
    /// @return maliciousMembers Identifiers of group members who signed the
    ///         malicious DKG result hash.
    function challengeResult(Data storage self)
        internal
        returns (uint32[] memory maliciousMembers)
    {
        require(
            currentState(self) == State.CHALLENGE,
            "current state is not CHALLENGE"
        );

        require(
            block.number <=
                self.submittedResultBlock +
                    self.parameters.resultChallengePeriodLength,
            "challenge period has already passed"
        );

        // Compute the actual group members hash by selecting actual members IDs
        // based on seed used for current DKG execution.
        bytes32 actualGroupMembersHash = keccak256(
            abi.encodePacked(
                self.sortitionPool.selectGroup(groupSize, bytes32(self.seed))
            )
        );

        require(
            self.submittedResultGroupMembersHash != actualGroupMembersHash,
            "unjustified challenge"
        );

        // Consider all members who signed the wrong result as malicious.
        maliciousMembers = self.submittedResultSigningMembers;

        // Adjust DKG result submission block start, so submission eligibility
        // starts from the beginning.
        self.resultSubmissionStartBlockOffset =
            block.number -
            self.startBlock -
            offchainDkgTime;

        emit DkgResultChallenged(self.submittedResultHash, msg.sender);

        submittedResultCleanup(self);

        return maliciousMembers;
    }

    /// @notice Set resultChallengePeriodLength parameter.
    function setResultChallengePeriodLength(
        Data storage self,
        uint256 newResultChallengePeriodLength
    ) internal {
        require(currentState(self) == State.IDLE, "current state is not IDLE");

        require(
            newResultChallengePeriodLength > 0,
            "new value should be greater than zero"
        );

        self
            .parameters
            .resultChallengePeriodLength = newResultChallengePeriodLength;
    }

    /// @notice Set resultSubmissionEligibilityDelay parameter.
    function setResultSubmissionEligibilityDelay(
        Data storage self,
        uint256 newResultSubmissionEligibilityDelay
    ) internal {
        require(currentState(self) == State.IDLE, "current state is not IDLE");

        require(
            newResultSubmissionEligibilityDelay > 0,
            "new value should be greater than zero"
        );

        self
            .parameters
            .resultSubmissionEligibilityDelay = newResultSubmissionEligibilityDelay;
    }

    /// @notice Completes DKG by cleaning up state.
    /// @dev Should be called after DKG times out or a result is approved.
    function complete(Data storage self) internal {
        delete self.startBlock;
        delete self.seed;
        delete self.resultSubmissionStartBlockOffset;
        submittedResultCleanup(self);
        self.sortitionPool.unlock();
    }

    /// @notice Cleans up submitted result state either after DKG completion
    ///         (as part of `cleanup` modifier) or after justified challenge.
    function submittedResultCleanup(Data storage self) private {
        delete self.submittedResultHash;
        delete self.submittedResultGroupMembersHash;
        delete self.submittedResultSigningMembers;
        delete self.submittedResultBlock;
        delete self.submittedResultMisbehavedMembers;
        delete self.resultSubmitter;
    }
}
