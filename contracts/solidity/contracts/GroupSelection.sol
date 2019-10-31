pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./utils/UintArrayUtils.sol";

/**
 * The group selection protocol is an interactive method of selecting candidate
 * group from the set of all stakers given a pseudorandom seed value.
 *
 * The protocol produces a representative result, where each staker's profit is
 * proportional to the number of tokens they have staked. Produced candidate
 * groups are of constant size.
 *
 * Group selection protocol accepts seed as an input - a pseudorandom value
 * used to construct candidate tickets. Each candidate group member can
 * submit their tickets. The maximum number of tickets one can submit depends
 * on their staking weight - relation of the minimum stake to the candidate's
 * stake.
 *
 * There is a certain timeout, expressed in blocks, when tickets can be
 * submitted. Each ticket is a mix of staker's address, virtual staker index
 * and group selection seed. Candidate group members are selected based on
 * the best tickets submitted. There has to be a minimum number of tickets
 * submitted, equal to the candidate group size so that the protocol can
 * complete successfully.
 */
library GroupSelection {

    using SafeMath for uint256;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    struct Storage {
        // All tickets submitted by member candidates during the current group
        // selection execution and accepted by the protocol for the
        // consideration.
        uint256[] tickets;

        // Information about each accepted ticket allowing to prove its
        // validity.
        mapping(uint256 => Proof) proofs;

        // Pseudorandom seed value used as an input for the group selection.
        uint256 seed;

        // Timeout in blocks after the reactive ticket submission is finished.
        uint256 ticketReactiveSubmissionTimeout;

        // Number of block at which the group selection started and ticket
        // submissions are accepted.
        uint256 ticketSubmissionStartBlock;

        // Indicates whether a group selection is currently in progress.
        // Concurrent group selections are not allowed.
        bool inProgress;
    }

    function start(Storage storage self, uint256 _seed) public {
        cleanup(self);
        self.inProgress = true;
        self.seed = _seed;
        self.ticketSubmissionStartBlock = block.number;
    }

    function stop(Storage storage self) public {
        cleanup(self);
        self.inProgress = false;
    }

    /**
     * @dev Submits ticket to request to participate in a new candidate group.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function submitTicket(
        Storage storage self,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight
    ) public {

        if (block.number > self.ticketSubmissionStartBlock + self.ticketReactiveSubmissionTimeout) {
            revert("Ticket submission is over");
        }

        if (self.proofs[ticketValue].sender != address(0)) {
            revert("Duplicate ticket");
        }

        if (isTicketValid(
            ticketValue,
            stakerValue,
            virtualStakerIndex,
            stakingWeight,
            self.seed
        )) {
            self.tickets.push(ticketValue);
            self.proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        } else {
            // TODO: should we slash instead of reverting?
            revert("Invalid ticket");
        }
    }

    /**
     * @dev Performs full verification of the ticket.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function isTicketValid(
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight,
        uint256 groupSelectionSeed
    ) internal view returns(bool) {
        bool isVirtualStakerIndexValid = virtualStakerIndex > 0 && virtualStakerIndex <= stakingWeight;
        bool isStakerValueValid = uint256(msg.sender) == stakerValue;
        bool isTicketValueValid = uint256(keccak256(abi.encodePacked(groupSelectionSeed, stakerValue, virtualStakerIndex))) == ticketValue;

        return isVirtualStakerIndexValid && isStakerValueValid && isTicketValueValid;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants(
        Storage storage self,
        uint256 groupSize
    ) public view returns (address[] memory) {
        require(
            block.number >= self.ticketSubmissionStartBlock + self.ticketReactiveSubmissionTimeout,
            "Ticket submission in progress"
        );

        require(self.tickets.length >= groupSize, "Not enough tickets submitted");

        uint256[] memory ordered = UintArrayUtils.sort(self.tickets);

        address[] memory selected = new address[](groupSize);

        for (uint i = 0; i < groupSize; i++) {
            Proof memory proof = self.proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
    }

    /**
     * @dev Cleanup data of previous group selection.
     */
    function cleanup(Storage storage self) internal {
        for (uint i = 0; i < self.tickets.length; i++) {
            delete self.proofs[self.tickets[i]];
        }
        delete self.tickets;
    }
}