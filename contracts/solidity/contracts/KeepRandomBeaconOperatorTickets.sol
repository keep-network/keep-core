pragma solidity ^0.5.4;

import "./utils/UintArrayUtils.sol";
import "./KeepRandomBeaconOperatorLinkedContract.sol";

/**
 * @title KeepRandomBeaconOperatorTickets
 * @dev A helper contract for operator contract to store tickets, proofs and
 * perform ticket sortition.
 */
contract KeepRandomBeaconOperatorTickets is KeepRandomBeaconOperatorLinkedContract {

    // Timeout in blocks after the initial ticket submission is finished.
    uint256 public ticketInitialSubmissionTimeout = 3;

    // Timeout in blocks after the reactive ticket submission is finished.
    uint256 public ticketReactiveSubmissionTimeout = 6;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal proofs;

    uint256 public ticketSubmissionStartBlock;
    uint256[] internal tickets;

    /**
     * @dev Reverts if ticket submission period is not over.
     */
    modifier whenTicketSubmissionIsOver() {
        require(
            block.number >= ticketSubmissionStartBlock + ticketReactiveSubmissionTimeout,
            "Ticket submission submission period must be over."
        );
        _;
    }

    /**
     * @dev Submits ticket to request to participate in a new candidate group.
     * @param staker Staker's address.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     * @param stakingWeight Staker's weight.
     * @param groupSelectionRelayEntry Group selection relay entry.
     */
    function submitTicket(
        address staker,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight,
        uint256 groupSelectionRelayEntry
    ) public onlyOperatorContract {

        if (block.number > ticketSubmissionStartBlock + ticketReactiveSubmissionTimeout) {
            revert("Ticket submission period is over.");
        }

        if (proofs[ticketValue].sender != address(0)) {
            revert("Ticket with the given value has already been submitted.");
        }

        // Invalid tickets are rejected and their senders penalized.
        if (isTicketValid(staker, ticketValue, stakerValue, virtualStakerIndex, stakingWeight, groupSelectionRelayEntry)) {
            tickets.push(ticketValue);
            proofs[ticketValue] = Proof(staker, stakerValue, virtualStakerIndex);
        } else {
            // TODO: should we slash instead of reverting?
            revert("Invalid ticket");
        }
    }

    /**
     * @dev Gets submitted tickets in ascending order.
     */
    function orderedTickets() public view returns (uint256[] memory) {
        return UintArrayUtils.sort(tickets);
    }

    /**
     * @dev Gets the number of submitted group candidate tickets so far.
     */
    function submittedTicketsCount() public view returns (uint256) {
        return tickets.length;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     * @param groupSize Size of a group in the threshold relay.
     */
    function selectedParticipants(uint256 groupSize) public view whenTicketSubmissionIsOver returns (address[] memory) {

        uint256[] memory ordered = orderedTickets();

        require(
            ordered.length >= groupSize,
            "The number of submitted tickets is less than specified group size."
        );

        address[] memory selected = new address[](groupSize);

        for (uint i = 0; i < groupSize; i++) {
            Proof memory proof = proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
    }

    /**
     * @dev Performs full verification of the ticket.
     * @param staker Address of the staker.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     * @param stakingWeight Staker's weight.
     * @param groupSelectionRelayEntry Group selection relay entry.
     */
    function isTicketValid(
        address staker,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight,
        uint256 groupSelectionRelayEntry
    ) public pure returns(bool) {
        bool isVirtualStakerIndexValid = virtualStakerIndex > 0 && virtualStakerIndex <= stakingWeight;
        bool isStakerValueValid = uint256(staker) == stakerValue;
        bool isTicketValueValid = uint256(keccak256(abi.encodePacked(groupSelectionRelayEntry, stakerValue, virtualStakerIndex))) == ticketValue;

        return isVirtualStakerIndexValid && isStakerValueValid && isTicketValueValid;
    }

    /**
     * @dev Sets ticket submission start block.
     */
    function setTicketSubmissionStartBlock(uint256 value) public onlyOperatorContract {
        ticketSubmissionStartBlock = value;
    }

    /**
     * @dev Cleanup data of previous group selection.
     */
    function cleanup() public onlyOperatorContract {
        for (uint i = 0; i < tickets.length; i++) {
            delete proofs[tickets[i]];
        }

        delete tickets;
    }

}
