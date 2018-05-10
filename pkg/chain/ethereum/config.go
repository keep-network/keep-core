package ethereum

type EthereumAccount struct {
	// Example: "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	Address string

	// Full path to a file like
	// "UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	KeyFile string

	// Password for accessing the account (you can't read this in from the .toml
	// file.  You have to use the enviroment variable or set the environment
	// variable to 'prompt' for it to interactivly prompt for a password)
	KeyFilePassword string
}

// This is the information that is read in from the .toml file
type EthereumConfig struct {
	// Example: "ws://192.168.0.157:8546"
	URL string

	// Names and addresses for contracts that can be called or for events received.
	ContractAddresses map[string]string

	Account EthereumAccount
}
