package chain

// Address is a chain-agnostic representation of a chain address.
type Address string

func (a Address) String() string {
	return string(a)
}