package groupselection

// tickets implements sort.Interface
type tickets []*ticket

// Len is the sort.Interface requirement for Tickets
func (ts tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for Tickets
func (ts tickets) Swap(i, j int) {
	ts[i].Proof.VirtualStakerIndex, ts[j].Proof.VirtualStakerIndex =
		ts[j].Proof.VirtualStakerIndex, ts[i].Proof.VirtualStakerIndex
}

// Less is the sort.Interface requirement for Tickets
func (ts tickets) Less(i, j int) bool {
	iVirtualStakeIndex := ts[i].Proof.VirtualStakerIndex
	jVirtualStakerIndex := ts[j].Proof.VirtualStakerIndex

	switch iVirtualStakeIndex.Cmp(jVirtualStakerIndex) {
	case -1:
		return true
	case 1:
		return false
	}

	return true
}
