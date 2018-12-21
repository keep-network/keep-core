pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";


contract KeepGroupImplV1 is Ownable {

    event GroupSelected(bytes32 groupPubKey);

    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _minStake;
    address internal _stakingProxy;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _submissionStart;

    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;
    bytes32[] internal _submissions;

    mapping (bytes32 => uint256) internal _submissionVotes;
    mapping (address => mapping (bytes32 => bool)) internal _hasVoted;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal _proofs;

    bytes32[] internal _groups;
    mapping (bytes32 => address[]) internal _groupMembers;

    mapping (string => bool) internal _initialized;

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue) public onlyOwner {
        _submissionStart = block.number;
        _randomBeaconValue = randomBeaconValue;
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    address internal _stakingContract;

    function authorizeStakingContract(address stakingContract) public onlyOwner {
        _stakingContract = stakingContract;
    }

    /**
     * @dev Submit ticket to request to participate in a new candidate group.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function submitTicket(
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public {

        if (block.number > _submissionStart + _timeoutSubmission) {
            revert("Ticket submission period is over.");
        }

        if (block.number > _submissionStart + _timeoutInitial && _tickets.length > _groupSize) {
            revert("Initial submission period is over with enough tickets received.");
        }

        // Invalid tickets are rejected and their senders penalized.
        if (!cheapCheck(msg.sender, stakerValue, virtualStakerIndex)) {
            // TODO: replace with a secure authorization protocol (addressed in RFC 4).
            TokenStaking stakingContract = TokenStaking(_stakingContract);
            stakingContract.authorizedTransferFrom(msg.sender, this, _minStake);
        } else {
            _tickets.push(ticketValue);
            _proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        }
    }

    /**
     * @dev Gets submitted tickets in ascending order.
     */
    function orderedTickets() public view returns (uint256[]) {
        return UintArrayUtils.sort(_tickets);
    }

    /**
     * @dev Gets participants ordered by their lowest-valued ticket.
     */
    function orderedParticipants() public view returns (address[]) {

        uint256[] memory ordered = orderedTickets();
        address[] memory participants = new address[](ordered.length);

        for (uint i = 0; i < ordered.length; i++) {
            Proof memory proof = _proofs[ordered[i]];
            participants[i] = proof.sender;
        }

        return participants;
    }

    /**
     * @dev Gets ticket proof.
     */
    function getTicketProof(uint256 ticketValue) public view returns (uint256, uint256) {
        return (
            _proofs[ticketValue].stakerValue,
            _proofs[ticketValue].virtualStakerIndex
        );
    }

    /**
     * @dev Performs surface-level validation of the ticket.
     * @param staker Address of the staker.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function cheapCheck(
        address staker,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public view returns(bool) {
        bool isVirtualStakerIndexValid = virtualStakerIndex > 0 && virtualStakerIndex <= stakingWeight(staker);
        bool isStakerValueValid = uint256(staker) == stakerValue;
        return isVirtualStakerIndexValid && isStakerValueValid;
    }

    /**
     * @dev Performs full verification of the ticket.
     * @param staker Address of the staker.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function costlyCheck(
        address staker,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public view returns(bool) {
        bool passedCheapCheck = cheapCheck(staker, stakerValue, virtualStakerIndex);
        uint256 expected = uint256(keccak256(abi.encodePacked(_randomBeaconValue, stakerValue, virtualStakerIndex)));
        return passedCheapCheck && ticketValue == expected;
    }

    function challenge(
        uint256 ticketValue
    ) public {

        Proof memory proof = _proofs[ticketValue];
        require(proof.sender != 0, "Ticket must be published.");

        // TODO: replace with a secure authorization protocol (addressed in RFC 4).
        // TokenStaking stakingContract = TokenStaking(_stakingContract);
        // if (costlyCheck(
        //     proof.sender,
        //     ticketValue,
        //     proof.stakerValue,
        //     proof.virtualStakerIndex
        // )) {
        //     // Slash challenger's stake balance.
        //     stakingContract.authorizedTransferFrom(msg.sender, this, _minStake);
        // } else {
        //     // Slash invalid ticket sender stake balance and reward the challenger.
        //     stakingContract.authorizedTransferFrom(proof.sender, msg.sender, _minStake);
        // }
    }

    function submitGroupPublicKey(bytes32 groupPubKey) public {

        require(
            // TODO: get participant number and implement slash/reward described in Phase 13
            AddressArrayUtils.contains(orderedParticipants(), msg.sender),
            "Sender must be in selected participants to be able to submit group pubkey"
        );

        _submissions.push(groupPubKey);
    }

    function voteForSubmission(bytes32 groupPubKey) public {

        require(
            // TODO: get participant number and implement slash/reward described in Phase 13
            AddressArrayUtils.contains(orderedParticipants(), msg.sender),
            "Sender must be in selected participants to be able to vote"
        );

        require(
            !_hasVoted[msg.sender][groupPubKey],
            "You already voted for this group pubkey"
        );

        _hasVoted[msg.sender][groupPubKey] = true;
        _submissionVotes[groupPubKey]++;
    }

    function getFinalResult() public {

        // TODO: Implement conflict resolution logic described in Phase 14
        // for now just use the first submission
        _groups.push(_submissions[0]);

        // Keep record of group members
        for (uint i = 0; i < _tickets.length; i++) {
            Proof memory proof = _proofs[_tickets[i]];
            _groupMembers[_submissions[0]].push(proof.sender);
        }

        emit GroupSelected(_submissions[0]);
    }

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Staking proxy contract.
     * @param stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param groupSize Minimum number of members in a group - to form a group.
     * @param timeoutInitial Timeout in blocks after the initial ticket submission is finished.
     * @param timeoutSubmission Timeout in blocks after the reactive ticket submission is finished.
     * @param timeoutChallenge Timeout in blocks after the period where tickets can be challenged is finished.
     */
    function initialize(
        address stakingProxy,
        uint256 minStake,
        uint256 groupSize,
        uint256 timeoutInitial,
        uint256 timeoutSubmission,
        uint256 timeoutChallenge
    ) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        _initialized["KeepGroupImplV1"] = true;
        _stakingProxy = stakingProxy;
        _minStake = minStake;
        _groupSize = groupSize;
        _timeoutInitial = timeoutInitial;
        _timeoutSubmission = timeoutSubmission;
        _timeoutChallenge = timeoutChallenge;
    }

    /**
     * @dev Gets staking weight.
     * @param staker Specifies the identity of the staker.
     * @return Number of how many virtual stakers can staker represent.
     */
    function stakingWeight(address staker) public view returns(uint256) {
        return StakingProxy(_stakingProxy).balanceOf(staker)/_minStake;
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param minStake Amount in KEEP.
     */
    function setMinimumStake(uint256 minStake) public onlyOwner {
        _minStake = minStake;
    }

    /**
     * @dev Get the minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function minimumStake() public view returns(uint256) {
        return _minStake;
    }

    /**
     * @dev Allows owner to change the groupSize.
     */
    function setGroupSize(uint256 groupSize) public onlyOwner {
        _groupSize = groupSize;
    }

    /**
     * @dev Return natural threshold, the value N virtual stakers' tickets would be expected
     * to fall below if the tokens were optimally staked, and the tickets' values were evenly 
     * distributed in the domain of the pseudorandom function.
     */
    function naturalThreshold() public view returns (uint256) {
        return _groupSize*((2**256)-1)/uint256((10**9)/_minStake);
    }
    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepGroupImplV1"];
    }

    /**
     * @dev Gets the minimum number of members in a group.
     */
    function groupSize() public view returns(uint256) {
        return _groupSize;
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return _groups.length;
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string) {
        return "V1";
    }
}
