package groupselection

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"sort"
	"testing"
)

func TestSortByValue(t *testing.T) {
	ticket1 := newTestTicket(5, 1001)
	ticket2 := newTestTicket(4, 1002)
	ticket3 := newTestTicket(1, 1003)
	ticket4 := newTestTicket(3, 1004)
	ticket5 := newTestTicket(2, 1005)

	tickets := []*ticket{
		ticket3,
		ticket5,
		ticket4,
		ticket1,
		ticket2,
	}

	sort.Stable(byValue(tickets))

	assertTicketAtIndex(t, tickets, 0, ticket1)
	assertTicketAtIndex(t, tickets, 1, ticket2)
	assertTicketAtIndex(t, tickets, 2, ticket3)
	assertTicketAtIndex(t, tickets, 3, ticket4)
	assertTicketAtIndex(t, tickets, 4, ticket5)
}

func assertTicketAtIndex(t *testing.T, tickets []*ticket, index int, ticket *ticket) {
	if !reflect.DeepEqual(ticket, tickets[index]) {
		t.Errorf(
			"unexpected ticket at index [%v]\nexpected: [%+v]\nactual:   [%+v]",
			index,
			ticket,
			tickets[index],
		)
	}
}

func newTestTicket(virtualStakerIndex uint32, value uint64) *ticket {
	var bytes [32]byte
	binary.BigEndian.PutUint64(bytes[:], value)

	return &ticket{
		value: bytes,
		proof: &proof{
			virtualStakerIndex: big.NewInt(int64(virtualStakerIndex)),
		},
	}
}
