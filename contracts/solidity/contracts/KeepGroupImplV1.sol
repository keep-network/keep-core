pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";


contract KeepGroupImplV1 is Ownable {

    using SafeMath for uint256;

    struct DkgResult {
        bool success;
        bytes groupPubKey;
        bytes disqualified;
        bytes inactive;
    }

    event DkgResultSubmittedEvent(uint256 requestId, uint256 memberIndex, bytes groupPubKey);     
    event DkgResultVoteEvent(uint256 requestId, uint256 memberIndex, bytes32 resultHash);
    event DkgResultElectedEvent(uint256 requestId, bool success);
    event GroupRegisteredEvent(uint256 requestId, bytes groupPubKey);

    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _minStake;
    address internal _stakingProxy;
    address internal _randomBeacon;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _submissionStart;

    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;
    bytes[] internal _submissions;
    bytes32[] internal _dkgResultHashes;

    mapping (uint256 => DkgResult) internal _requestIdToDkgResult;
    mapping (uint256 => bool) internal _dkgResultPublished; 
    mapping (bytes32 => DkgResult) internal _receivedSubmissions;
    mapping (bytes32 => bool) internal _votedDkg;
    mapping (bytes32 => uint256) internal _submissionVotes;
    mapping (bytes32 => bool) internal _resultPublished;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal _proofs;

    bytes[] internal _groups;
    mapping (bytes => address[]) internal _groupMembers;

    mapping (string => bool) internal _initialized;

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue) public {
        require(msg.sender == _randomBeacon);
        cleanup();
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

        // Invalid tickets are rejected and their senders penalized.
        if (!cheapCheck(msg.sender, stakerValue, virtualStakerIndex)) {
            // TODO: replace with a secure authorization protocol (addressed in RFC 4).
            TokenStaking stakingContract = TokenStaking(_stakingContract);
            stakingContract.authorizedTransferFrom(msg.sender, address(this), _minStake);
        } else {
            _tickets.push(ticketValue);
            _proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        }
    }

    /**
     * @dev Gets submitted tickets in ascending order.
     */
    function orderedTickets() public view returns (uint256[] memory) {
        return UintArrayUtils.sort(_tickets);
    }
    
    /**
     * @dev Gets selected tickets in ascending order.
     */
    function selectedTickets() public view returns (uint256[] memory) {

        require(
            block.number > _submissionStart + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );

        uint256[] memory ordered = orderedTickets();
        uint256[] memory selected = new uint256[](_groupSize);

        for (uint i = 0; i < _groupSize; i++) {
            selected[i] = ordered[i];
        }

        return selected;
    }

    /**
     * @dev Gets participants ordered by their lowest-valued ticket.
     */
    function orderedParticipants() public view returns (address[] memory) {

        uint256[] memory ordered = orderedTickets();
        address[] memory participants = new address[](ordered.length);

        for (uint i = 0; i < ordered.length; i++) {
            Proof memory proof = _proofs[ordered[i]];
            participants[i] = proof.sender;
        }

        return participants;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants() public view returns (address[] memory) {

        require(
            block.number > _submissionStart + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );

        uint256[] memory ordered = orderedTickets();
        address[] memory selected = new address[](_groupSize);

        for (uint i = 0; i < _groupSize; i++) {
            Proof memory proof = _proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
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

    /**
     * @dev Check if member is inactive.
     * @param dqBytes bytes representing disqualified members.
     * @param memberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function _isDisqualified(bytes memory dqBytes, uint256 memberIndex) internal pure returns (bool){
        return dqBytes[memberIndex] != 0x00;
    }

     /** 
     * @dev Check if member is inactive.
     * @param iaBytes bytes representing inactive members.
     * @param memberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function _isInactive(bytes memory iaBytes, uint256 memberIndex) internal pure returns (bool){
        return iaBytes[memberIndex] != 0x00;
    }

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 13 of the protocol.
     * @param requestId Relay request ID assosciated with DKG protocol execution.
     * @param memberIndex position the user claims to be.
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
        uint256 memberIndex,
        bool success,
        bytes memory groupPubKey,
        bytes memory disqualified,
        bytes memory inactive
    ) public {
        require(validateIndex(memberIndex));

        bytes32 resultHash = keccak256(abi.encode(success, groupPubKey, disqualified, inactive));
        bytes32 submitterID = keccak256(abi.encodePacked(msg.sender, memberIndex, _randomBeaconValue));
    
        require(eligibleSubmitter(memberIndex), "not an eligible submitter");
        require(!_votedDkg[submitterID], "already voted for or submitted a result");
        
        //check empty for first submitter incentives. Should not re enter. voting begins after first submission
        if(!_resultPublished[resultHash]){
            if(_dkgResultHashes.length == 0){
                //TODO: punish/reward
                //First submitter incentive logic.kl
            }
            _receivedSubmissions[resultHash] = DkgResult(success, groupPubKey, disqualified, inactive);
            _dkgResultHashes.push(resultHash);
            _submissionVotes[resultHash] = 1;
            _votedDkg[submitterID] = true;//cannot vote after submiting DKG result
            _resultPublished[resultHash] = true;
            emit DkgResultSubmittedEvent(requestId, memberIndex, groupPubKey);
        }
        else{
            _addVote(resultHash, submitterID, requestId, memberIndex);
        }  
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return _dkgResultPublished[requestId];
    }

    /** 
     * @dev Check if submitter is eligible to submit.
     * @param memberIndex the claimed index of the submitter.
     * @return true if the submitter is eligible. False otherwise.
     */
    function eligibleSubmitter(uint memberIndex) internal pure returns (bool) {
        // TODO: implement
        require(memberIndex != 0);
        return true;
    }

    /**
     * @dev Check if provided index belongs to staker owner.
     * @param memberIndex the claimed index of the user.
     * @return true if the ticket at the given index is issued by msg.sender. False otherwise.
     */
    function validateIndex(uint memberIndex) internal pure returns (bool) {
        // TODO: implement
        require(memberIndex != 0);
        return true;
    }

    /** 
    * @dev returns the final agreed upn DKG result or error result
    */
    function getFinalResult(uint256 requestId)public returns (bytes memory) {
        bytes32 leadingResult;
        uint highestVoteN;

        for(uint i = 0; i < _dkgResultHashes.length; i++){
            if(_submissionVotes[_dkgResultHashes[i]] > highestVoteN){
                highestVoteN = _submissionVotes[_dkgResultHashes[i]];
                leadingResult = _dkgResultHashes[i];
            }
        }
        _dkgResultPublished[requestId] = true;
        
        address[] memory members = orderedParticipants();
        bytes memory groupPublicKey = _receivedSubmissions[leadingResult].groupPubKey;
        for (uint i = 0; i < _groupSize; i++) {
            if(!_isInactive(_receivedSubmissions[leadingResult].inactive, i) &&
                !_isDisqualified(_receivedSubmissions[leadingResult].disqualified, i)){
                _groupMembers[groupPublicKey].push(members[i]);
            }
        }
        _groups.push(groupPublicKey);
        
        emit DkgResultElectedEvent(requestId, true);
        emit GroupRegisteredEvent(requestId, _receivedSubmissions[leadingResult].groupPubKey);
        cleanup();
        return _receivedSubmissions[leadingResult].groupPubKey;
    }

    /** 
     * @dev receives vote for provided resultHash.
     * @param memberIndex the claimed index of the user.
     * @param resultHash Hash of DKG result to vote for
     */
    function voteOnDkgResult(
        uint256 requestId,
        uint256 memberIndex,
        bytes32 resultHash
    ) public {
        bytes32 submitterID = keccak256(abi.encodePacked(msg.sender, memberIndex, _randomBeaconValue));
        _addVote(resultHash, submitterID, requestId, memberIndex);
    }
 
     /** 
     * @dev adds vote for provided resultHash.
     * @param resultHash the hash of the DKG result the submitter claims is correct.
     * @param submitterID Hash of the submitterID index and address
     */
    function _addVote(bytes32 resultHash, bytes32 submitterID, uint256 requestId, uint256 memberIndex) internal{
        _votedDkg[submitterID] = true;
        _submissionVotes[resultHash] += 1;
        emit DkgResultVoteEvent(requestId, memberIndex, resultHash);
    }

     /** 
     * @dev get resultHashes and corresponding votes
     * @return  bytes32 array of result Hashes and uint256 array of corresponding result votes 
     */
    function getDkgResultSubmissions() public view returns (bytes32[] memory, uint256[] memory) {
        uint256[] memory votes = new uint256[](_dkgResultHashes.length);
        for(uint i = 0; i < _dkgResultHashes.length; i++){
            votes[i] = _submissionVotes[_dkgResultHashes[i]];
        }
        return (_dkgResultHashes, votes);
    }  


    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Staking proxy contract.
     * @param stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param randomBeacon Address of a random beacon contract that will be linked to this contract.
     * @param minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param groupSize Size of a group in the threshold relay.
     * @param groupThreshold Minimum number of interacting group members needed to produce a relay entry.
     * @param timeoutInitial Timeout in blocks after the initial ticket submission is finished.
     * @param timeoutSubmission Timeout in blocks after the reactive ticket submission is finished.
     * @param timeoutChallenge Timeout in blocks after the period where tickets can be challenged is finished.
     */
    function initialize(
        address stakingProxy,
        address randomBeacon,
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
        _randomBeacon = randomBeacon;
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
     * @dev ticketInitialSubmissionTimeout is the duration (in blocks) the
     * staker has to submit tickets that fall under the natural threshold
     * to satisfy the initial ticket timeout (see group selection, phase 2a).
     */
    function ticketInitialSubmissionTimeout() public view returns (uint256) {
        return _timeoutInitial;
    }

    /**
     * @dev ticketReactiveSubmissionTimeout is the duration (in blocks) the
     * staker has to submit any tickets that did not fall under the natural
     * threshold. This final chance to submit tickets is called reactive
     * ticket submission (defined in the group selection algorithm, 2b).
     */
    function ticketReactiveSubmissionTimeout() public view returns (uint256) {
        return _timeoutSubmission;
    }

    /**
     * @dev ticketChallengeTimeout is the duration (in blocks) the staker
     * has to submit any challenges for tickets that fail any checks.
     */
    function ticketChallengeTimeout() public view returns (uint256) {
        return _timeoutChallenge;
    }

    /**
     * @dev ticketSubmissionStartBlock block number at which current group
     * selection started.
     */
    function ticketSubmissionStartBlock() public view returns (uint256) {
        return _submissionStart;
    }

    /**
     * @dev Return total number of all tokens issued.
     */
    function tokenSupply() public pure returns (uint256) {
        return (10**9) * (10**18);
    }

    /**
     * @dev Return natural threshold, the value N virtual stakers' tickets would be expected
     * to fall below if the tokens were optimally staked, and the tickets' values were evenly 
     * distributed in the domain of the pseudorandom function.
     */
    function naturalThreshold() public view returns (uint256) {
        uint256 space = 2**256-1; // Space consisting of all possible tickets.
        return _groupSize.mul(space.div(tokenSupply().div(_minStake)));
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
    function version() public pure returns (string memory) {
        return "V1";
    }

    /**
     * @dev Cleanup data of previous group selection.
     */
    function cleanup() private {

        for (uint i = 0; i < _tickets.length; i++) {
            delete _proofs[_tickets[i]];
        }

        delete _tickets;
        delete _dkgResultHashes;
    }

}