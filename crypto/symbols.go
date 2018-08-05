package crypto

// Represents collection for storing symbols/currencies available for trading
type Symbols struct {
	Assets  map[string][]string `json:"assets"`
	Symbols []string            `json:"symbols"`
}
