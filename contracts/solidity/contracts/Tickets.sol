pragma solidity ^0.5.4;


/**
 * @title Tickets contract for testing purposes. Will need to merge later.
 */

contract Tickets  {
    uint256 public groupSize = 10;
    uint256[] internal tickets;
    uint256 internal tail = 0;

    // Map simulates a linked list. key -> value are both indices of the tickets array.
    // 'key' index holds a higher value and points to an index that holds a next lower value
    // Ex. tickets = [151, 42, 175, 7]
    // orderedLinkedTicketIndices[0] -> 1
    // orderedLinkedTicketIndices[1] -> 3
    // orderedLinkedTicketIndices[2] -> 0
    // orderedLinkedTicketIndices[3] -> 3 note: index that holds a smallest value points to itself.
    mapping(uint256 => uint256) public orderedLinkedTicketIndices;

    constructor() public {
        orderedLinkedTicketIndices[tail] = 0; // simulates nil
    }

    function submitTicket(uint256 newTicketValue) public {
        uint256 oldTail = tail;
        uint256[] memory ordered = createOrderedLinkedTicketIndices();
        orderedTickets = ordered;

        if (tickets.length < groupSize) {
            // bigger than the existing biggest
            if (tickets.length == 0 || newTicketValue > tickets[tail]) {
                tickets.push(newTicketValue);
                if (tickets.length > 1) {
                    tail = tickets.length-1;
                    orderedLinkedTicketIndices[tail] = oldTail;
                }
            // smaller than the existing smallest
            } else if (newTicketValue < tickets[ordered[0]]) {
                tickets.push(newTicketValue);
                orderedLinkedTicketIndices[tickets.length - 1] = tickets.length - 1; // last element point to itself
                orderedLinkedTicketIndices[ordered[0]] = tickets.length - 1;
            // tickets[smallest] < newTicketValue < tickets[max]
            } else {
                tickets.push(newTicketValue);
                uint j = findIndexForNewTicket(newTicketValue, ordered);
                orderedLinkedTicketIndices[tickets.length - 1] = orderedLinkedTicketIndices[j];
                orderedLinkedTicketIndices[j] = tickets.length - 1;

                jIndex = j;
            }
        } else if (newTicketValue < tickets[tail]) { // tickets[groupSize]
            // replacing existing smallest with a smaller
            if (newTicketValue < ordered[0]) {
                tickets[ordered[0]] = newTicketValue;
            } else {
                uint j = findIndexForNewTicket(newTicketValue, ordered);
                tickets[tail] = newTicketValue;
                // do not change the order if a new ticket is still highest
                if (j != tail) {
                    uint newTail = orderedLinkedTicketIndices[tail];
                    orderedLinkedTicketIndices[j] = tail;
                    orderedLinkedTicketIndices[tail] = tickets.length - 1;
                    tail = newTail;
                }
                jIndex = j;
            }
        }
    }

    // use binary search to find an index for a new ticket
    function findIndexForNewTicket(uint256 newTicketValue, uint256[] memory ordered) internal view returns (uint256) {
        uint lo = 0;
        uint hi = ordered.length - 1;
        uint mid = 0;
        while (lo <= hi) {
            mid = lo + (hi - lo) / 2;
            if (newTicketValue < tickets[ordered[mid]]) {
                hi = mid - 1;
            } else if (newTicketValue > tickets[ordered[mid]]) {
                lo = mid + 1;
            } else {
                return ordered[mid];
            }
        }
        return ordered[lo];
    }

    function createOrderedLinkedTicketIndices() internal view returns (uint256[] memory) {
        uint256[] memory ordered = new uint256[](tickets.length);
        if (ordered.length > 0) {
            ordered[tickets.length-1] = tail;
            if (ordered.length > 1) {
                for (int i = int(tickets.length - 2); i >= 0; i--) {
                    ordered[uint(i)] = orderedLinkedTicketIndices[ordered[uint(i) + 1]];
                }
            }
        }

        return ordered;
    }

    function getTail() public view returns (uint256) {
        return tail;
    }

    function getTicketMaxValue() public view returns (uint256) {
        return tickets[tail];
    }

    function cleanup() public {
        delete tickets;
    }

    function getTickets() public view returns (uint256[] memory) {
        return tickets;
    }

    function getOrderedLinkedTicketIndices(uint index) public view returns (uint256) {
        return orderedLinkedTicketIndices[index];
    }
    
    function getTicketLength() public view returns (uint256) {
        return tickets.length;
    }

    /* debug helper - can be removed later*/
    uint256[] internal orderedTickets;
    function getOrdered() public view returns (uint256[] memory) {
        return orderedTickets;
    }

    /* debug helper - can be removed later*/
    uint256 jIndex;
    function getJIndex() public view returns (uint256) {
        return jIndex;
    }

}
