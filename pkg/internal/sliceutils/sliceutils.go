package sliceutils

// Contains checks if slice of integers contains a given value.
// TODO Use this func in https://github.com/keep-network/keep-core/pull/379/files#diff-0d2fbe694cc2b75e75577f3df3e10b7aR45
func Contains(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}
