pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";


contract KeepGroupImplV1 is Ownable {

    using SafeMath for uint256;

    event GroupSelected(bytes32 groupPubKey);

    struct DkgResult {
        bool success;
        bytes32 groupPubKey;
        bytes disqualified;
        bytes inactive;
    }

    event DkgResultPublishedEvent(uint256 requestId);
    
    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _minStake;
    address internal _stakingProxy;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _submissionStart;
    uint256 internal _requestId;

    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;
    bytes32[] internal _submissions;

    mapping (uint256 => DkgResult) internal _requestIdToDkgResult;
    mapping (uint256 => bool) internal _dkgResultPublished;
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
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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

        if (block.number > _submissionStart + _timeoutInitial && _tickets.length >= _groupSize) {
            revert("Initial submission period is over with enough tickets received.");
        }

        if (block.number < _submissionStart + _timeoutInitial && ticketValue > naturalThreshold()) {
            revert("Ticket must be below natural threshold during initial submission period.");
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
    function getTicketProof(uint256 ticketValue) public view returns (address, uint256, uint256) {
        return (
            _proofs[ticketValue].sender,
            _proofs[ticketValue].stakerValue,
            _proofs[ticketValue].virtualStakerIndex
        );
    }

    /**
     * @dev Performs surface-level validation of the ticket.
     * @param staker Address of the staker.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 13 of the protocol.
     * @param requestId Relay request ID assosciated with DKG protocol execution.
     * @param success Result of DKG protocol execution; true if success, false otherwise.
     * @param groupPubKey Group public key generated as a result of protocol execution.
     * @param disqualified bytes representing disqualified group members; 1 at the specific index 
     * means that the member has been disqualified. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     * @param inactive bytes representing inactive group members; 1 at the specific index means
     * that the member has been marked as inactive. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     */
    function submitDkgResult(
        uint256 requestId,
        bool success, 
        bytes32 groupPubKey,
        bytes disqualified,
        bytes inactive
    ) public {
        _requestIdToDkgResult[requestId] = DkgResult(success, groupPubKey, disqualified, inactive);
        _dkgResultPublished[requestId] = true;
        _requestId = requestId;
  
        emit DkgResultPublishedEvent(requestId);
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return _dkgResultPublished[requestId];
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
     * @param groupSize Size of a group in the threshold relay.
     * @param groupThreshold Minimum number of interacting group members needed to produce a relay entry.
     * @param timeoutInitial Timeout in blocks after the initial ticket submission is finished.
     * @param timeoutSubmission Timeout in blocks after the reactive ticket submission is finished.
     * @param timeoutChallenge Timeout in blocks after the period where tickets can be challenged is finished.
     */
    function initialize(
        address stakingProxy,
        uint256 minStake,
        uint256 groupThreshold,
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
        _groupThreshold = groupThreshold;
        _timeoutInitial = timeoutInitial;
        _timeoutSubmission = timeoutSubmission;
        _timeoutChallenge = timeoutChallenge;
    }

    /**
     * @dev Checks that the specified user has enough stake.
     * @param staker Specifies the identity of the staker.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address staker) public view returns(bool) {
        return StakingProxy(_stakingProxy).balanceOf(staker) >= _minStake;
    }

    /**
     * @dev Gets staking weight.
     * @param staker Specifies the identity of the staker.
     * @return Number of how many virtual stakers can staker represent.
     */
    function stakingWeight(address staker) public view returns(uint256) {
        return StakingProxy(_stakingProxy).balanceOf(staker)/_minStake;
    }

     /* 
     * @dev Check if member is disqualified.
     * @param groupPubKey public key of the specified group.
     * @param gmemberIndex position of the member to check.
     * @return true if staker is disqualified, false otherwise.
     */
    function isDisqualified(bytes32 groupPubKey, uint256 memberIndex) public view returns (bool){
        //NOTE: variable _requestId updated on submitDkgResult()
        //better way to get DkgResult using groupPubKey?
        require(_dkgResultPublished[_requestId] == true, 
            "DKG Result is not currently submitted");
        require(_requestIdToDkgResult[_requestId].groupPubKey == groupPubKey,
            "recent DGK submission does not match Group Public Key");
        return _requestIdToDkgResult[_requestId].disqualified[memberIndex] != 0x00;
    }

     /*
     * @dev Check if member is inactive.
     * @param groupPubKey public key of the specified group.
     * @param gmemberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function isInactive(bytes32 groupPubKey, uint256 memberIndex) public view returns (bool){
        //NOTE: variable _requestId updated on submitDkgResult()
        //better way to get DkgResult using groupPubKey?
        require(_dkgResultPublished[_requestId] == true, 
            "DKG Result is not currently submitted");
        require(_requestIdToDkgResult[_requestId].groupPubKey == groupPubKey,
            "recent DGK submission does not match Group Public Key");
        return _requestIdToDkgResult[_requestId].inactive[memberIndex] != 0x00;
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
        uint256 space = 2**256-1; // Space consisting of all possible tickets.
        uint256 tokens = 10**9; // Total number of all tokens issued.
        return _groupSize.mul(space.div(tokens.div(_minStake)));
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepGroupImplV1"];
    }

    /**
     * @dev Gets size of a group in the threshold relay.
     */
    function groupSize() public view returns(uint256) {
        return _groupSize;
    }

    /**
     * @dev Gets number of interacting group members needed to produce a relay entry.
    */
    function groupThreshold() public view returns(uint256) {
        return _groupThreshold;
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
