package server

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	config *ServerConfig
}

func NewHttpServer(config *ServerConfig) *HttpServer {
	if config == nil {
		config = DefaultServerConfig
	}

	return &HttpServer{
		config: config,
	}
}

func (server *HttpServer) Start() {
	router := gin.Default()
	port := strings.Join([]string{":", strconv.Itoa(server.config.Port)}, "")

	Init(router)

	router.Run(port)
}
