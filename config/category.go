package config

type Category int

const (
	General Category = iota
	Ethereum
	Network
	Storage
	Metrics
	Diagnostics
	Tbtc
	Developer
)

var AllCategories = []Category{
	General,
	Ethereum,
	Network,
	Storage,
	Metrics,
	Diagnostics,
	Tbtc,
	Developer,
}
