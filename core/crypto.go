package core

type Crypto struct {
	Name     string
	Shortcut string
	Code     uint32
}

var (
	BITCOIN Crypto = Crypto{
		Name:     "BITCOIN",
		Shortcut: "BTC",
	}
	ETHEREUM Crypto = Crypto{
		Name:     "ETHEREUM",
		Shortcut: "ETH",
	}
	LITECOIN Crypto = Crypto{
		Name:     "LITECOIN",
		Shortcut: "LTC",
	}
)
