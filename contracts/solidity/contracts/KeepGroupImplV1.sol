pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";


contract KeepGroupImplV1 is Ownable {

    using SafeMath for uint256;

    event OnGroupRegistered(bytes groupPubKey);

    struct DkgResult {
        bool success;
        bytes groupPubKey;
        bytes disqualified;
        bytes inactive;
    }

    event DkgResultPublishedEvent(address publisher, bytes groupPubKey);

    // Legacy code moved from Random Beacon contract
    // TODO: refactor according to the Phase 14
    event SubmitGroupPublicKeyEvent(bytes groupPublicKey, uint256 activationBlockHeight);//add RequestId replacement

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

    mapping (bytes32 => bool) internal _votedDkg;
    mapping (bytes32 => bool) internal _resultPublished;
    mapping (address => DkgResult) internal _publisherToDkgResult;
    mapping (bytes32 => DkgResult) internal _receivedSubmissions;
    mapping (bytes32 => uint) internal _submissionVotes;

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
     * @dev Gets selected tickets in ascending order.
     */
    function selectedTickets() public view returns (uint256[]) {

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
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants() public view returns (address[]) {

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
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
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
    
    /*
     * @dev Check if member is inactive.
     * @param dqBytes bytes representing disqualified members.
     * @param memberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function _isDisqualified(bytes dqBytes, uint256 memberIndex) internal view returns (bool){
        return dqBytes[memberIndex] != 0x00;
    }

     /*
     * @dev Check if member is inactive.
     * @param iaBytes bytes representing inactive members.
     * @param gmemberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function _isInactive(bytes iaBytes, uint256 memberIndex) internal view returns (bool){
        return iaBytes[memberIndex] != 0x00;
    }

    /*
     * @dev receives a DKG result submission, will be added if conditions are met.
     * @param index the claimed index (P_i) of the user.
     * @param success Result of DKG protocol execution; true if success, false otherwise.
     * @param groupPubKey Group public key generated as a result of protocol execution.
     * @param disqualified bytes representing disqualified group members; 1 at the specific index 
     * means that the member has been disqualified. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     * @param inactive bytes representing inactive group members; 1 at the specific index means
     * that the member has been marked as inactive. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     */
   
    function receiveSubmission(
        uint256 index, 
        bool success, 
        bytes groupPubKey,
        bytes disqualified,
        bytes inactive)public {

        require(validateIndex(index));

        bytes32 resultHash = keccak256(abi.encodePacked(success, groupPubKey, disqualified, inactive));
        bytes32 submitterID = keccak256(abi.encodePacked(msg.sender, index, _randomBeaconValue));
    
        require(eligibleSubmitter(index), "not an eligible submitter yet");
        require(!_votedDkg[submitterID], "already voted for or submitted a result");
         //Find better place for this, checking every submission seems pointless
        require(
            _tickets.length >= _groupSize,
            "There should be enough valid tickets submitted to form a group."
            );
        //check empty for first submitter incentives. Should not re enter. voting begins after first submission
        if(!_resultPublished[resultHash]){
            if(_dkgResultHashes.length == 0){
                //TODO: punish/reward
                //First submitter incentive logic.
            }
            _receivedSubmissions[resultHash] = DkgResult(success, groupPubKey, disqualified, inactive);
            _dkgResultHashes.push(resultHash);
            _submissionVotes[resultHash] = 1;
            _votedDkg[submitterID] = true;//cannot vote after submiting DKG result
            _resultPublished[resultHash] = true;
            _publisherToDkgResult[msg.sender] = _receivedSubmissions[resultHash];
            emit DkgResultPublishedEvent(msg.sender, groupPubKey);
            
        }
        else{
            _addVote(resultHash, submitterID);
        }  
    }
   
    /*
     * @dev receives vote for provided resultHash.
     * @param index the claimed index of the user.
     * @param resultHash Hash of DKG result to vote for
     */
    function receiveVote(uint256 index, bytes32 resultHash)public {
        require(validateIndex(index));
        bytes32 submitterID = keccak256(abi.encodePacked(msg.sender, index, _randomBeaconValue));
        require(!_votedDkg[submitterID], "already voted for or submitted a result");
        require(_submissionVotes[resultHash] != 0, "Result hash not published yet");
        _addVote(resultHash, submitterID);
    }

    /*
     * @dev Check if submitter is eligible to submit.
     * @param index the claimed index of the submitter.
     * @return true if the submitter is eligible. False otherwise.
     */
    function eligibleSubmitter(uint index) public returns (bool){
        uint T_init = _submissionStart + _timeoutChallenge;
        uint T_step = 2; //time between eligibility increments Placeholder
        require(block.number > T_init, "Ticket submission challenge period must be over.");
        if(index == 1) return true;
        
        //(2* (T_step)) -> time for first submitter to submit DKG Result.
        //No way to calculate DKG time on-chain, so DKG result submission opens on ticket-challenge close.
        //first submitter usable time period is really (2 * T_step - DKG time).
        else if(block.number >= ((T_init + (2 * (T_step))) + ((index-2) * T_step))){
            return true;
        }
        else return false;
    }

    /*
     * @dev Check if provided index belongs to staker owner.
     * @param index the claimed index of the user.
     * @return true if the ticket at the given index is issued by msg.sender. False otherwise.
     */
    function validateIndex(uint index)public returns(bool){   
        require(index != 0, "can't be 0 index"); 
        require(index <= _groupSize, "must be within selected range");
        uint256[] memory ordered = orderedTickets();
        return(_proofs[ordered[index - 1]].sender == msg.sender);
    }

    /*
     * @dev add vote for provided resultHash.
     * @param resultHash the hash of the DKG result the submitter claims is correct.
     * @param submitterID Hash of the submitterID index and address
     */
    function _addVote(bytes32 resultHash, bytes32 submitterID) internal{
        _votedDkg[submitterID] = true;
        _submissionVotes[resultHash] += 1;
    }

    /*
     * @dev returns the final DKG result.
     */
    function submitGroupPublicKey()public returns (bytes32) {
        bytes32 leadingResult;
        uint highestVoteN;
        uint highestVoteNtemp;
        uint totalVotes;
        uint f_max = _groupSize/2 + 1;

        //TODO:
        //method cannot be called before voting period is over or everyone has voted as it is liked with cleanup()

        for(uint i = 0; i < _dkgResultHashes.length; i++){
            highestVoteNtemp = _submissionVotes[_dkgResultHashes[i]];
            if(highestVoteNtemp > highestVoteN){
                highestVoteN = highestVoteNtemp;
                leadingResult = _dkgResultHashes[i];
            }
            totalVotes += highestVoteNtemp;
        }
        if(totalVotes - highestVoteN >= f_max){
            cleanup();
            return 0x0;
            //TODO
            //return Result.failure(disqualified = [])
        }
        else{
            address[] memory members = orderedParticipants();
            bytes memory groupPublicKey = _receivedSubmissions[leadingResult].groupPubKey;
            for (i = 0; i < _groupSize; i++) {
                if(!_isInactive(_receivedSubmissions[leadingResult].inactive, i) &&
                    !_isDisqualified(_receivedSubmissions[leadingResult].disqualified, i)){
                    _groupMembers[groupPublicKey].push(members[i]);
                }
            }
            _groups.push(groupPublicKey);
            //emist _randomBeaconValue instead of RequestID
            emit SubmitGroupPublicKeyEvent(groupPublicKey, _randomBeaconValue);
            cleanup();
            return leadingResult;
            //TODO:
            //return value as DKG result
        }

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
    function tokenSupply() public view returns (uint256) {
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
    function version() public pure returns (string) {
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

        // TODO: cleanup DkgResults
    }

    /**
     * @dev Get block height
     */
    function blockHeight() public view returns(uint256) {
        return block.number;
    }

     /*
     * @dev Helper - get block height when user at given index is eligible to submit DKG result
     * @param index the claimed index of the submitter.
     * @return the block number when a user is eligible.
     */
    function eligibleTime(uint index) public view returns(uint){
        uint T_init = _submissionStart + _timeoutChallenge;
        uint T_step = 2; //time between eligibility increments Placeholder
        if(index == 1 ){
            return T_init + 1;
        }
        return ((T_init + (2 * (T_step))) + ((index-2) * T_step));
    }
}
