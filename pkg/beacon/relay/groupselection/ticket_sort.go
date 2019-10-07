package groupselection

// byValue implements sort.Interface sorting tickets by their value.
type byValue []*ticket

// Len is the sort.Interface requirement for ticket ordering.
func (bv byValue) Len() int {
	return len(bv)
}

// Swap is the sort.Interface requirement for ticket ordering.
func (bv byValue) Swap(i, j int) {
	bv[i], bv[j] = bv[j], bv[i]
}

// Less is the sort.Interface requirement for ticket ordering.
func (bv byValue) Less(i, j int) bool {
	return bv[i].intValue().Cmp(bv[j].intValue()) < 1
}
