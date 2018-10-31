package local

type localStakeMonitoring struct {
}

func (lsm *localStakeMonitoring) HasMinimumStake(address string) (bool, error) {
	return true, nil
}
