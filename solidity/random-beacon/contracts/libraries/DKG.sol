// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "./BytesLib.sol";

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
        // Block number from the moment of the DKG result submission.
        uint256 submittedResultBlock;
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

    /// @dev The minimum number of group members needed to interact according to
    ///      the protocol to produce a relay entry. The adversary can not learn
    ///      anything about the key as long as it does not break into
    ///      groupThreshold+1 of members.
    uint256 public constant groupThreshold = 33;

    /// @dev The minimum number of active and properly behaving group members
    ///      during the DKG needed to accept the result. This number is higher
    ///      than `groupThreshold` to keep a safety margin for members becoming
    ///      inactive after DKG so that the group can still produce a relay
    ///      entry.
    uint256 public constant activeThreshold = 58; // 90% of groupSize

    /// @notice Time in blocks after which DKG result is complete and ready to be
    //          published by clients.
    uint256 public constant offchainDkgTime = 5 * (1 + 5) + 2 * (1 + 10) + 20;

    event DkgStarted(uint256 indexed seed);

    event DkgResultSubmitted(
        bytes32 indexed resultHash,
        uint256 indexed seed,
        uint256 submitterMemberIndex,
        bytes indexed groupPubKey,
        uint8[] misbehavedMembersIndices,
        bytes signatures,
        uint256[] signingMembersIndices,
        uint32[] members
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

    /// @notice Determines the current state of group creation. It doesn't take
    ///         timeouts into consideration. The timeouts should be tracked and
    ///         notified separately.
    function currentState(Data storage self, SortitionPool sortitionPool)
        internal
        view
        returns (State state)
    {
        state = State.IDLE;

        if (sortitionPool.isLocked()) {
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
    function lockState(Data storage self, SortitionPool sortitionPool)
        internal
    {
        require(
            currentState(self, sortitionPool) == State.IDLE,
            "current state is not IDLE"
        );

        emit DkgStateLocked();

        sortitionPool.lock();
    }

    function start(
        Data storage self,
        SortitionPool sortitionPool,
        uint256 seed
    ) internal {
        require(
            currentState(self, sortitionPool) == State.AWAITING_SEED,
            "current state is not AWAITING_SEED"
        );

        self.startBlock = block.number;
        self.seed = seed;

        // slither-disable-next-line reentrancy-events
        emit DkgStarted(seed);
    }

    function submitResult(
        Data storage self,
        SortitionPool sortitionPool,
        Result calldata result
    ) external {
        require(
            currentState(self, sortitionPool) == State.AWAITING_RESULT,
            "current state is not AWAITING_RESULT"
        );
        require(
            !hasDkgTimedOut(self, sortitionPool),
            "dkg timeout already passed"
        );

        verify(self, sortitionPool, result);

        // TODO: Check with sortition pool that all members have minimum stake.
        // Check all members in one call or at least members that signed the result.
        // We need to know in advance that there will be something that we can
        // slash the members from.

        self.submittedResultHash = keccak256(abi.encode(result));
        self.submittedResultBlock = block.number;

        emit DkgResultSubmitted(
            self.submittedResultHash,
            self.seed,
            result.submitterMemberIndex,
            result.groupPubKey,
            result.misbehavedMembersIndices,
            result.signatures,
            result.signingMembersIndices,
            result.members
        );
    }

    /// @notice Checks if DKG timed out. The DKG timeout period includes time required
    ///         for off-chain protocol execution and time for the result publication
    ///         for all group members. After this time result cannot be submitted
    ///         and DKG can be notified about the timeout. DKG period is adjusted
    ///         by result submission offset that include blocks that were mined
    ///         while invalid result has been registered until it got challenged.
    /// @return True if DKG timed out, false otherwise.
    function hasDkgTimedOut(Data storage self, SortitionPool sortitionPool)
        internal
        view
        returns (bool)
    {
        return
            currentState(self, sortitionPool) == State.AWAITING_RESULT &&
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
    function verify(
        Data storage self,
        SortitionPool sortitionPool,
        Result calldata result
    ) internal view {
        // TODO: Verify if submitter is valid staker and signatures come from valid
        // stakers https://github.com/keep-network/keep-core/pull/2654#discussion_r728226906.

        require(result.submitterMemberIndex > 0, "Invalid submitter index");

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
            groupSize - result.misbehavedMembersIndices.length >=
                activeThreshold,
            "Too many members misbehaving during DKG"
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
        require(signaturesCount >= groupThreshold, "Too few signatures");
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
        // TODO: Revisit at optimization stage. We probably don't need to fetch
        //       addresses of all members. Addresses of signing members should
        //       be enough though this will require some additional processing
        //       here. Compare two approaches gas-wise against a real pool.
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

            require(
                membersAddresses[memberIndex - 1] == recoveredAddress,
                "Invalid signature"
            );
        }
    }

    /// @notice Notifies about DKG timeout.
    function notifyTimeout(Data storage self, SortitionPool sortitionPool)
        internal
    {
        require(hasDkgTimedOut(self, sortitionPool), "dkg has not timed out");

        emit DkgTimedOut();
    }

    /// @notice Notifies about the seed was not delivered and restores the
    ///         initial DKG state (IDLE).
    function notifySeedTimedOut(Data storage self, SortitionPool sortitionPool)
        internal
    {
        require(
            currentState(self, sortitionPool) == State.AWAITING_SEED,
            "current state is not AWAITING_SEED"
        );

        emit DkgSeedTimedOut();

        sortitionPool.unlock();
    }

    /// @notice Approves DKG result. Can be called when the challenge period for
    ///         the submitted result is finished. Considers the submitted result
    ///         as valid. For the first `resultSubmissionEligibilityDelay`
    ///         blocks after the end of the challenge period can be called only
    ///         by the DKG result submitter. After that time, can be called by
    ///         anyone.
    /// @dev Can be called after a challenge period for the submitted result.
    /// @param result Result to approve. Must match the submitted result stored
    ///        during `submitResult`.
    /// @return misbehavedMembers Identifiers of members who misbehaved during DKG.
    function approveResult(
        Data storage self,
        SortitionPool sortitionPool,
        Result calldata result
    ) external returns (uint32[] memory misbehavedMembers) {
        require(
            currentState(self, sortitionPool) == State.CHALLENGE,
            "current state is not CHALLENGE"
        );

        uint256 challengePeriodEnd = self.submittedResultBlock +
            self.parameters.resultChallengePeriodLength;

        require(
            block.number > challengePeriodEnd,
            "challenge period has not passed yet"
        );

        require(
            keccak256(abi.encode(result)) == self.submittedResultHash,
            "result under approval is different than the submitted one"
        );

        // Extract submitter member address. Submitter member index is in
        // range [1, 64] so we need to -1 when fetching identifier from members
        // array.
        address submitterMember = sortitionPool.getIDOperator(
            result.members[result.submitterMemberIndex - 1]
        );

        require(
            msg.sender == submitterMember ||
                block.number >
                challengePeriodEnd +
                    self.parameters.resultSubmissionEligibilityDelay,
            "Only the DKG result submitter can approve the result at this moment"
        );

        // Extract misbehaved members identifiers. Misbehaved members indices
        // are in range [1, 64], so we need to -1 when fetching identifiers from
        // members array.
        misbehavedMembers = new uint32[](
            result.misbehavedMembersIndices.length
        );
        for (uint256 i = 0; i < result.misbehavedMembersIndices.length; i++) {
            misbehavedMembers[i] = result.members[
                result.misbehavedMembersIndices[i] - 1
            ];
        }

        emit DkgResultApproved(self.submittedResultHash, msg.sender);

        return misbehavedMembers;
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    /// @dev Can be called during a challenge period for the submitted result.
    /// @param result Result to challenge. Must match the submitted result
    ///        stored during `submitResult`.
    /// @return maliciousResultHash Hash of the malicious result.
    /// @return maliciousMembers Identifiers of group members who signed the
    ///         malicious DKG result hash.
    function challengeResult(
        Data storage self,
        SortitionPool sortitionPool,
        Result calldata result
    )
        external
        returns (bytes32 maliciousResultHash, uint32[] memory maliciousMembers)
    {
        require(
            currentState(self, sortitionPool) == State.CHALLENGE,
            "current state is not CHALLENGE"
        );

        require(
            block.number <=
                self.submittedResultBlock +
                    self.parameters.resultChallengePeriodLength,
            "challenge period has already passed"
        );

        require(
            keccak256(abi.encode(result)) == self.submittedResultHash,
            "result under challenge is different than the submitted one"
        );

        // Compute the actual group members hash by selecting actual members IDs
        // based on seed used for current DKG execution.
        bytes32 actualGroupMembersHash = keccak256(
            abi.encodePacked(
                sortitionPool.selectGroup(groupSize, bytes32(self.seed))
            )
        );

        require(
            keccak256(abi.encodePacked(result.members)) !=
                actualGroupMembersHash,
            "unjustified challenge"
        );

        // Consider result hash as malicious.
        maliciousResultHash = self.submittedResultHash;

        // Consider all members who signed the wrong result as malicious and
        // extract their identifiers. Signing members indices are in range
        // [1, 64], so we need to -1 when fetching identifiers from members
        // array.
        maliciousMembers = new uint32[](result.signingMembersIndices.length);
        for (uint256 i = 0; i < result.signingMembersIndices.length; i++) {
            maliciousMembers[i] = result.members[
                result.signingMembersIndices[i] - 1
            ];
        }

        // Adjust DKG result submission block start, so submission eligibility
        // starts from the beginning.
        self.resultSubmissionStartBlockOffset =
            block.number -
            self.startBlock -
            offchainDkgTime;

        emit DkgResultChallenged(self.submittedResultHash, msg.sender);

        submittedResultCleanup(self);

        return (maliciousResultHash, maliciousMembers);
    }

    /// @notice Set resultChallengePeriodLength parameter.
    function setResultChallengePeriodLength(
        Data storage self,
        SortitionPool sortitionPool,
        uint256 newResultChallengePeriodLength
    ) internal {
        require(
            currentState(self, sortitionPool) == State.IDLE,
            "current state is not IDLE"
        );

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
        SortitionPool sortitionPool,
        uint256 newResultSubmissionEligibilityDelay
    ) internal {
        require(
            currentState(self, sortitionPool) == State.IDLE,
            "current state is not IDLE"
        );

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
    function complete(Data storage self, SortitionPool sortitionPool) internal {
        delete self.startBlock;
        delete self.seed;
        delete self.resultSubmissionStartBlockOffset;
        submittedResultCleanup(self);
        sortitionPool.unlock();
    }

    /// @notice Cleans up submitted result state either after DKG completion
    ///         (as part of `complete` method) or after justified challenge.
    function submittedResultCleanup(Data storage self) private {
        delete self.submittedResultHash;
        delete self.submittedResultBlock;
    }
}
