package core

type Symbols struct {
	Assets  map[string][]string `json:"assets"`
	Symbols []string            `json:"symbols"`
}
