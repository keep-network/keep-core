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

        // Map simulates a sorted linked list of ticket values by their indexes.
        // key -> value represent indices from the tickets[] array.
        // 'key' index holds an index of a ticket and 'value' holds an index
        // of the next ticket. Tickets are sorted by their value in
        // descending order starting from the tail.
        // Ex. tickets = [151, 42, 175, 7]
        // tail: 2 because tickets[2] = 175
        // previousTicketIndex[0] -> 1
        // previousTicketIndex[1] -> 3
        // previousTicketIndex[2] -> 0
        // previousTicketIndex[3] -> 3 note: index that holds a lowest
        // value points to itself because there is no `nil` in Solidity.
        // Traversing from tail: [2]->[0]->[1]->[3] result in 175->151->42->7
        mapping(uint256 => uint256) previousTicketIndex;

        // Tail represents an index of a ticket in a tickets[] array which holds
        // the largest ticket value and is used by `previousTicketIndex`. It is
        // a tail of the linked list defined by `previousTicketIndex`.
        uint256 tail;

        // Size of a group in the threshold relay.
        uint256 groupSize;
    }

    /**
     * @dev Starts group selection protocol.
     * @param _seed pseudorandom seed value used as an input for the group
     * selection. All submitted tickets needs to have the seed mixed-in into the
     * value.
     */
    function start(Storage storage self, uint256 _seed) public {
        cleanup(self);
        self.inProgress = true;
        self.seed = _seed;
        self.ticketSubmissionStartBlock = block.number;
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
        uint256[] memory ordered = createTicketValueOrderedIndices(self);

        // any ticket goes when the tickets array size is lower than the group size
        if (self.tickets.length < self.groupSize) {
            // no tickets or larger than the current largest
            if (self.tickets.length == 0 || newTicketValue > self.tickets[self.tail]) {
                self.tickets.push(newTicketValue);
                if (self.tickets.length > 1) {
                    self.tail = self.tickets.length-1;
                    self.previousTicketIndex[self.tail] = oldTail;
                }
            // lower than the current lowest
            } else if (newTicketValue < self.tickets[ordered[0]]) {
                self.tickets.push(newTicketValue);
                // last element points to itself
                self.previousTicketIndex[self.tickets.length - 1] = self.tickets.length - 1;
                // prev index of a lowest ticket value points to a new lowest
                self.previousTicketIndex[ordered[0]] = self.tickets.length - 1;
            // larger than the lowest ticket value and lower than the largest ticket value,
            } else {
                self.tickets.push(newTicketValue);
                uint j = findReplacementIndex(self, newTicketValue, ordered);
                self.previousTicketIndex[self.tickets.length - 1] = self.previousTicketIndex[j];
                self.previousTicketIndex[j] = self.tickets.length - 1;
            }
        } else if (newTicketValue < self.tickets[self.tail]) {
            // new ticket is lower than currently lowest
            if (newTicketValue < ordered[0]) {
                // replacing largest ticket with a lowest
                self.tickets[self.tail] = newTicketValue;
                // updating the previousTicketIndex map
                uint newTail = self.previousTicketIndex[self.tail];
                self.previousTicketIndex[ordered[0]] = self.tail;
                self.previousTicketIndex[self.tail] = self.tail;
                self.tail = newTail;
            } else { // new ticket is between lowest and largest
                uint j = findReplacementIndex(self, newTicketValue, ordered);
                self.tickets[self.tail] = newTicketValue;
                // do not change the order if a new ticket is still largest
                if (j != self.tail) {
                    uint newTail = self.previousTicketIndex[self.tail];
                    self.previousTicketIndex[j] = self.tail;
                    self.previousTicketIndex[self.tail] = self.tickets.length - 1;
                    self.tail = newTail;
                }
            }
        }
    }

    /**
     * @dev Use binary search to find an index for a new ticket in the tickets[] array
     */
    function findReplacementIndex(
        Storage storage self,
        uint256 newTicketValue,
        uint256[] memory ordered
    ) internal view returns (uint256) {
        uint256 lo = 0;
        uint256 hi = ordered.length - 1;
        uint256 mid = 0;
        while (lo <= hi) {
            mid = (lo + hi) >> 1;
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

    /**
     * @dev Creates an array of ticket indexes based on their values in the ascending order:
     *
     * ordered[n-1] = tail
     * ordered[n-2] = previousTicketIndex[tail]
     * ordered[n-3] = previousTicketIndex[ordered[n-2]]
     */
    function createTicketValueOrderedIndices(Storage storage self) internal view returns (uint256[] memory) {
        uint256[] memory ordered = new uint256[](self.tickets.length);
        if (ordered.length > 0) {
            ordered[self.tickets.length-1] = self.tail;
            if (ordered.length > 1) {
                for (uint256 i = self.tickets.length - 1; i > 0; i--) {
                    ordered[i-1] = self.previousTicketIndex[ordered[i]];
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

        address[] memory selected = new address[](groupSize);
        uint256 linkedIndex;
        uint256 ticketValue;

        if (selected.length > 0) {
            linkedIndex = self.tail;
            ticketValue = self.tickets[linkedIndex];
            selected[self.tickets.length-1] = self.proofs[ticketValue].sender;

            if (selected.length > 1) {
                for (uint256 i = self.tickets.length - 1; i > 0; i--) {
                    ticketValue = self.tickets[self.previousTicketIndex[linkedIndex]];
                    linkedIndex = self.previousTicketIndex[linkedIndex];
                    selected[i-1] = self.proofs[ticketValue].sender;
                }
            }
        }

        return selected;
    }

    /**
     * @dev Clears up data of the group selection.
     */
    function cleanup(Storage storage self) internal {
        for (uint i = 0; i < self.tickets.length; i++) {
            delete self.proofs[self.tickets[i]];
            delete self.previousTicketIndex[i];
        }
        delete self.tickets;
    }
}