package server

const (
	DefaultPort = 8000
)

type Config struct {
	Port int
}

var DefaultConfig *Config = &Config{
	Port: DefaultPort,
}
