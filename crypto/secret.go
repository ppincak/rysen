package crypto

// Account used to acces an exchange
type Secret interface {
	ApiKey() string
	SecretKey() string
}
