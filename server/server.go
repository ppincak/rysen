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

// Create new server
func NewServer(app *App, config *Config) *Server {
	if config == nil {
		config = DefaultConfig
	}

	// TODO: replace gin.Default
	return &Server{
		config: config,
		engine: gin.Default(),
		router: NewRouter(app),
	}
}

// Run the server
func (server *Server) Run() {
	port := strings.Join([]string{":", strconv.Itoa(server.config.Port)}, "")
	server.router.Init(server.engine)
	server.engine.Run(port)
}
