package ethereum

// GetBlockNo return the block number last seen
func (ebc *ethereumBlockCounter) GetBlockNo() int {
	return ebc.latestBlockHeight
}
