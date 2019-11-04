pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

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
        // Tickets submitted by member candidates during the current group
        // selection execution and accepted by the protocol for the
        // consideration.
        uint256[] tickets;

        // Information about each accepted ticket allowing to prove its
        // validity.
        mapping(uint256 => Proof) proofs;

        // Pseudorandom seed value used as an input for the group selection.
        uint256 seed;

        // Timeout in blocks after which the ticket submission is finished.
        uint256 ticketSubmissionTimeout;

        // Number of block at which the group selection started and from which
        // ticket submissions are accepted.
        uint256 ticketSubmissionStartBlock;

        // Indicates whether a group selection is currently in progress.
        // Concurrent group selections are not allowed.
        bool inProgress;

        // Map simulates a linked list. key -> value are both indices in the
        // tickets[] array.
        // 'key' index holds a higher number and points to an index that
        // holds a next lower ticket number. That number is a ticket value.
        // Ex. tickets = [151, 42, 175, 7]
        // orderedLinkedTicketIndices[0] -> 1
        // orderedLinkedTicketIndices[1] -> 3
        // orderedLinkedTicketIndices[2] -> 0
        // orderedLinkedTicketIndices[3] -> 3 note: index that holds a smallest
        // value points to itself.
        mapping(uint256 => uint256) orderedLinkedTicketIndices;

        // Tail represents an index of a ticket in a tickets[] array which holds
        // the largest ticket value.
        uint256 tail;

        // Size of a group in the threshold relay.
        uint256 groupSize;
    }

    /**
     * @dev Starts group selection protocol.
     * @param _seed pseudorandom seed value used as an input for the group
     * selection. All submitted tickets needs to have the seed mixed-in into the
     * value.
     * @param _groupSize size of a member group that produce a relay.
     */
    function start(Storage storage self, uint256 _seed, uint256 _groupSize) public {
        cleanup(self);
        self.inProgress = true;
        self.seed = _seed;
        self.ticketSubmissionStartBlock = block.number;
        self.orderedLinkedTicketIndices[self.tail] = 0; // simulates nil
        self.groupSize = _groupSize;
    }

    /**
     * @dev Stops group selection protocol clearing up all the submitted
     * tickets.
     */
    function stop(Storage storage self) public {
        cleanup(self);
        self.inProgress = false;
    }

    /**
     * @dev Submits ticket to request to participate in a new candidate group.
     * @param ticketValue Keccak-256 hash with input values of group selection
     * seed, staker address and virtualStakerIndex.
     * @param stakerValue Staker's address as an integer.
     * @param virtualStakerIndex Index of a virtual staker - number within
     * a range of 1 to staker's weight.
     */
    function submitTicket(
        Storage storage self,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight
    ) public {
        if (block.number > self.ticketSubmissionStartBlock.add(self.ticketSubmissionTimeout)) {
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
            addTicket(self, ticketValue);
            self.proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        } else {
            // TODO: should we slash instead of reverting?
            revert("Invalid ticket");
        }
    }

    /**
     * @dev Performs full verification of the ticket.
     */
    function isTicketValid(
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex,
        uint256 stakingWeight,
        uint256 groupSelectionSeed
    ) internal view returns(bool) {
        bool isVirtualStakerIndexValid = virtualStakerIndex > 0 && virtualStakerIndex <= stakingWeight;
        bool isStakerValueValid = stakerValue == uint256(msg.sender);
        bool isTicketValueValid = ticketValue == uint256(keccak256(abi.encodePacked(groupSelectionSeed, stakerValue, virtualStakerIndex)));

        return isVirtualStakerIndexValid && isStakerValueValid && isTicketValueValid;
    }

    /**
     * @dev Add a new verified ticket to the tickets[] array.
     */
    function addTicket(Storage storage self, uint256 newTicketValue) internal {
        uint256 oldTail = self.tail;
        uint256[] memory ordered = createOrderedLinkedTicketIndices(self);

        if (self.tickets.length < self.groupSize) {
            // larger than the existing largest
            if (self.tickets.length == 0 || newTicketValue > self.tickets[self.tail]) {
                self.tickets.push(newTicketValue);
                if (self.tickets.length > 1) {
                    self.tail = self.tickets.length-1;
                    self.orderedLinkedTicketIndices[self.tail] = oldTail;
                }
            // smaller than the existing smallest
            } else if (newTicketValue < self.tickets[ordered[0]]) {
                self.tickets.push(newTicketValue);
                // last element points to itself
                self.orderedLinkedTicketIndices[self.tickets.length - 1] = self.tickets.length - 1;
                self.orderedLinkedTicketIndices[ordered[0]] = self.tickets.length - 1;
            // self.tickets[smallest] < newTicketValue < self.tickets[max]
            } else {
                self.tickets.push(newTicketValue);
                uint j = findIndexForNewTicket(self, newTicketValue, ordered);
                self.orderedLinkedTicketIndices[self.tickets.length - 1] = self.orderedLinkedTicketIndices[j];
                self.orderedLinkedTicketIndices[j] = self.tickets.length - 1;
            }
        } else if (newTicketValue < self.tickets[self.tail]) {
            // replacing existing smallest with a smaller
            if (newTicketValue < ordered[0]) {
                self.tickets[ordered[0]] = newTicketValue;
            } else {
                uint j = findIndexForNewTicket(self, newTicketValue, ordered);
                self.tickets[self.tail] = newTicketValue;
                // do not change the order if a new ticket is still largest
                if (j != self.tail) {
                    uint newTail = self.orderedLinkedTicketIndices[self.tail];
                    self.orderedLinkedTicketIndices[j] = self.tail;
                    self.orderedLinkedTicketIndices[self.tail] = self.tickets.length - 1;
                    self.tail = newTail;
                }
            }
        }
    }

    // use binary search to find an index for a new ticket in the tickets[] array
    function findIndexForNewTicket(
        Storage storage self,
        uint256 newTicketValue,
        uint256[] memory ordered
    ) internal view returns (uint256) {
        uint lo = 0;
        uint hi = ordered.length - 1;
        uint mid = 0;
        while (lo <= hi) {
            mid = lo + (hi - lo) / 2;
            if (newTicketValue < self.tickets[ordered[mid]]) {
                hi = mid - 1;
            } else if (newTicketValue > self.tickets[ordered[mid]]) {
                lo = mid + 1;
            } else {
                return ordered[mid];
            }
        }

        return ordered[lo];
    }

    // Creates an array of ticket indexes based on their values in the ascending
    // order.
    function createOrderedLinkedTicketIndices(Storage storage self) internal view returns (uint256[] memory) {
        uint256[] memory ordered = new uint256[](self.tickets.length);
        if (ordered.length > 0) {
            ordered[self.tickets.length-1] = self.tail;
            if (ordered.length > 1) {
                for (int i = int(self.tickets.length - 2); i >= 0; i--) {
                    ordered[uint(i)] = self.orderedLinkedTicketIndices[ordered[uint(i) + 1]];
                }
            }
        }

        return ordered;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants(
        Storage storage self,
        uint256 groupSize
    ) public view returns (address[] memory) {
        require(
            block.number >= self.ticketSubmissionStartBlock.add(self.ticketSubmissionTimeout),
            "Ticket submission in progress"
        );

        require(self.tickets.length >= groupSize, "Not enough tickets submitted");

        uint256[] memory ordered = createOrderedLinkedTicketIndices(self);
        address[] memory selected = new address[](groupSize);
        for (uint i = 0; i < groupSize; i++) {
            Proof memory proof = self.proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
    }

    /**
     * @dev Clears up data of the group selection.
     */
    function cleanup(Storage storage self) internal {
        for (uint i = 0; i < self.tickets.length; i++) {
            delete self.proofs[self.tickets[i]];
        }
        delete self.tickets;
    }
}