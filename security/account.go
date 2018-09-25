package security

type Account struct {
	Name           string          `json:"name"`
	Authentication *Authentication `json:"authentication"`
}

type Authentication struct {
	*Secret
}
