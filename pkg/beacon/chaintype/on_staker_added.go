package chaintype

// OnStakerAdded is the data for the
// OnStakerAdded event.  This type may
// only be needed in Milestone 1 - it may
// change at Milestone 2.
type OnStakerAdded struct {
	Index         int
	GroupMemberID string
}
