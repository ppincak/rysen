package security

type Secret struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}
