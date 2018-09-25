package crypto

// Represents crypto asset eg: { Symbol: "Bitcoin", Amount: 1.256 }
type Asset struct {
	AssetName string  `json:"assetName"`
	Amount    float64 `json:"amount"`
}
