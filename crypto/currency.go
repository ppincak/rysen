package crypto

type Currency struct {
	Name     string
	Shortcut string
	Code     uint32
}

var (
	BITCOIN Currency = Currency{
		Name:     "BITCOIN",
		Shortcut: "BTC",
	}
	ETHEREUM Currency = Currency{
		Name:     "ETHEREUM",
		Shortcut: "ETH",
	}
	LITECOIN Currency = Currency{
		Name:     "LITECOIN",
		Shortcut: "LTC",
	}
)
