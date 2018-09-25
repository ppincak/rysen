package server

const (
	DefaultPort = 8000
)

// Config container
type Config struct {
	Port int
}

// Default config
var DefaultConfig *Config = &Config{
	Port: DefaultPort,
}
