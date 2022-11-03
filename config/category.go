package config

type Category int

const (
	General Category = iota
	Ethereum
	Electrum
	Network
	Storage
	ClientInfo
	Tbtc
	Maintainer
	Developer
)

// StartCmdCategories are categories needed for the start command.
var StartCmdCategories = []Category{
	General,
	Ethereum,
	Network,
	Storage,
	ClientInfo,
	Tbtc,
	Developer,
}

// AllCategories are all available categories.
var AllCategories = []Category{
	General,
	Ethereum,
	Electrum,
	Network,
	Storage,
	ClientInfo,
	Tbtc,
	Maintainer,
	Developer,
}
