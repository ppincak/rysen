package crypto

// Represents single crypto currency
type Currency struct {
	Name     string
	Shortcut string
	Code     uint32
}

// Main Crypto currencies
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
	TETHER Currency = Currency{
		Name:     "TETHER",
		Shortcut: "USDT",
	}
)
