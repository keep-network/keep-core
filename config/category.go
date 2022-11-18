package config

type Category int

const (
	General Category = iota
	Ethereum
	BitcoinElectrum
	Network
	Storage
	ClientInfo
	Tbtc
	TbtcMaintainer
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
	BitcoinElectrum,
	Network,
	Storage,
	ClientInfo,
	Tbtc,
	TbtcMaintainer,
	Developer,
}
