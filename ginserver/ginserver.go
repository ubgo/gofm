package ginserver

import (
	"errors"
	"log"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config
	Router       *gin.Engine
	RouterGroups map[string]*gin.RouterGroup
}

func (server Server) AddRouterGroup(name string, path string) *gin.RouterGroup {
	server.RouterGroups[name] = server.Router.Group(path)
	return server.RouterGroups[name]
}

func (server Server) GetRouterGroup(name string) (*gin.RouterGroup, error) {
	group := server.RouterGroups[name]
	if group == nil {
		return nil, errors.New("no route group found")
	}
	return server.RouterGroups[name], nil
}

func (server Server) Start() {
	port := server.config.Port
	url := "http://localhost:" + port
	log.Println("Http Sever started at " + color.CyanString(url))
	server.Router.Run(":" + port)
}

func (server Server) pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"app":     server.config.AppName,
	})
}

func New(opts ...Option) Server {
	cfg := config{
		AppName: "app",
		Port:    "7001",
		IsProd:  true,
	}

	cfg.options(opts...)

	if cfg.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	if cfg.BeforeHandler != nil {
		router.Use(cfg.BeforeHandler)
	}

	server := Server{
		config:       cfg,
		Router:       router,
		RouterGroups: make(map[string]*gin.RouterGroup),
	}

	router.GET("/ping", server.pingHandler)

	return server
}
