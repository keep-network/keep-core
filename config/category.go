package config

type Category int

const (
	General Category = iota
	Ethereum
	Network
	Storage
	ClientInfo
	Tbtc
	Developer
)

var AllCategories = []Category{
	General,
	Ethereum,
	Network,
	Storage,
	ClientInfo,
	Tbtc,
	Developer,
}
