package server

const (
	DefaultPort = 8000
)

type ServerConfig struct {
	Port int
}

var DefaultServerConfig *ServerConfig = &ServerConfig{
	Port: DefaultPort,
}
