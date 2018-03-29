package RelayContract

// RelayContractContext currently the context just contains a
// debug flag.  It will probably contain information about
// how to and where to log and if this is running in test mode
// or in production mode in the future.
type RelayContractContext struct {
	dbOnFlag bool
}

// dbOn returns true if the debugging flag is enabled.
func (ctx *RelayContractContext) dbOn() bool {
	return ctx.dbOnFlag
}

// SetDebug sets the debugging flag to true/false.
func (ctx *RelayContractContext) SetDebug(b bool) {
	ctx.dbOnFlag = b
}
