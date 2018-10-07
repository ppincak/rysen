package crypto

// Represents collection for storing symbols/currencies available for trading
type Symbols struct {
	// Map of crypto to available trading options/symbols
	Assets map[string][]string `json:"assets"`
	// List of all symbols
	Symbols []string `json:"symbols"`
}
