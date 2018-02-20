package RelayContract

type RelayContractContext struct {
	dbOnFlag bool
}

func (ctx *RelayContractContext) dbOn() bool {
	return ctx.dbOnFlag
}

func (ctx *RelayContractContext) SetDebug(b bool) {
	ctx.dbOnFlag = b
}
