package server

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app    *App
	config *Config
	engine *gin.Engine
	router *Router
}

func NewServer(app *App, config *Config) *Server {
	if config == nil {
		config = DefaultConfig
	}

	return &Server{
		config: config,
		engine: gin.Default(),
		router: NewRouter(app),
	}
}

func (server *Server) Run() {
	port := strings.Join([]string{":", strconv.Itoa(server.config.Port)}, "")
	server.router.Init(server.engine)
	server.engine.Run(port)
}
