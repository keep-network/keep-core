package groupselection

// tickets implements sort.Interface
type tickets []*ticket

// Len is the sort.Interface requirement for Tickets
func (ts tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for Tickets
func (ts tickets) Swap(i, j int) {
	ts[i].proof.virtualStakerIndex, ts[j].proof.virtualStakerIndex =
		ts[j].proof.virtualStakerIndex, ts[i].proof.virtualStakerIndex
}

// Less is the sort.Interface requirement for Tickets
func (ts tickets) Less(i, j int) bool {
	iVirtualStakeIndex := ts[i].proof.virtualStakerIndex
	jVirtualStakerIndex := ts[j].proof.virtualStakerIndex

	switch iVirtualStakeIndex.Cmp(jVirtualStakerIndex) {
	case -1:
		return true
	case 1:
		return false
	}

	return true
}
