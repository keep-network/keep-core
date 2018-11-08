package local

type localStakeMonitor struct {
}

func (lsm *localStakeMonitor) HasMinimumStake(address string) (bool, error) {
	return true, nil
}
