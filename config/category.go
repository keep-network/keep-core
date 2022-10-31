package config

type Category int

const (
	General Category = iota
	Ethereum
	Bitcoin
	Network
	Storage
	ClientInfo
	Tbtc
	Developer
	Maintainer
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
