pragma solidity ^0.5.4;


/**
 * @title Tickets contract for testing purposes. Will need to merge later.
 */

contract Tickets  {
    uint256 public groupSize = 10;
    uint256[] internal tickets;
    uint256 tail = 0;

    // Map simulates a linked list. key -> value are both indices of the tickets array.
    // 'key' index holds a higher value and points to an index that holds a next lower value
    mapping(uint256 => uint256) previousTicketsByIndex;

    constructor() public {
        previousTicketsByIndex[tail] = 0; // simulate nil, refer to an index outside tickets array.
    }

    function createOrderedTicketIndexes() internal view returns (uint256[] memory) {
        uint256[] memory ordered = new uint256[](tickets.length);
        if (ordered.length > 0) {
            ordered[tickets.length-1] = tail;
            if (ordered.length > 1) {
                for (uint i = tickets.length - 2; i > 0; i--) {
                    ordered[i] = previousTicketsByIndex[ordered[i + 1]];
                }
            }
        }

        return ordered;
    }

    // use binary search to find an index to a new ticket
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
                return mid;
            }
        }
        return mid;
    }

    function submitTicket(uint256 newTicketValue) public {
        uint256 oldTail = tail;
        uint256[] memory ordered = createOrderedTicketIndexes();
        orderedTickets = ordered;

        if (tickets.length < groupSize) {
            // add to tickets when tickets accept all
            if (tickets.length == 0 || newTicketValue > tickets[tail]) {
                tickets.push(newTicketValue);
                if (tickets.length > 1) {
                    tail = tickets.length-1;
                    previousTicketsByIndex[tail] = oldTail;
                }
            } else { // tickets[0] < newTicketValue < tickets[max]
                uint j = findIndexForNewTicket(newTicketValue, ordered);
                tickets.push(newTicketValue);
                uint temp = previousTicketsByIndex[j];
                previousTicketsByIndex[j] = tickets.length - 1;
                previousTicketsByIndex[tickets.length - 1] = temp;
                jIndex = j;
            }
        } else {
            //TODO replacing part
            //isSmallerThanCurrentHighestValue(newTicketValue)
        }
    }

    function isSmallerThanCurrentHighestValue(uint256 ticketValue) public view returns (bool) {
        return tickets[tail] > ticketValue;
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

    function getPreviousTicketsByIndex(uint index) public view returns (uint256) {
        return previousTicketsByIndex[index];
    }
    
    function getTicketLength() public view returns (uint256) {
        return tickets.length;
    }

    /* debug helper */
    uint256[] internal orderedTickets;
    function getOrdered() public view returns (uint256[] memory) {
        return orderedTickets;
    }

    /* debug helper */
    uint256 jIndex;
    function getJIndex() public view returns (uint256) {
        return jIndex;
    }

}
