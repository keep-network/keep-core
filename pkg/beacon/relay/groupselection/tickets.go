package groupselection

// tickets implements sort.Interface
type tickets []*ticket

// Len is the sort.Interface requirement for tickets
func (ts tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for tickets
func (ts tickets) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Less is the sort.Interface requirement for tickets
func (ts tickets) Less(i, j int) bool {
	return ts[i].intValue().Cmp(ts[j].intValue()) < 1
}
