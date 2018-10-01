package security

type Secret struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

type Account struct {
	Name           string          `json:"name"`
	Exchange       string          `json:"exchange"`
	Authentication *Authentication `json:"authentication"`
}

type Authentication struct {
	*Secret
}
