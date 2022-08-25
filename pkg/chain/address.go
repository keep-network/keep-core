package chain

// Address is a chain-agnostic representation of a chain address.
type Address string

func (a Address) String() string {
	return string(a)
}

// Addresses is a list of Address.
type Addresses []Address

// Set transform Addresses into a set of unique items.
func (a Addresses) Set() map[Address]bool {
	set := make(map[Address]bool)

	for _, address := range a {
		set[address] = true
	}

	return set
}
