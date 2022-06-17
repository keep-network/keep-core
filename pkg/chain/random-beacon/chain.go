package beacon

// Handle for interaction with the Random Beacon module contracts.
type Handle interface {
	JoinSortitionPool() error
}
