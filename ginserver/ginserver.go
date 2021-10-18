package ginserver

import (
	"errors"
	"log"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	BeforeHandler gin.HandlerFunc
}

type Server struct {
	Config
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
	port := viper.GetString("server.port")
	url := "http://localhost:" + port
	log.Println("Http Sever started at " + color.CyanString(url))
	server.Router.Run(":" + port)
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"app":     viper.GetString("app_name"),
	})
}

func New(config Config) Server {
	if viper.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	if config.BeforeHandler != nil {
		router.Use(config.BeforeHandler)
	}
	router.GET("/ping", pingHandler)

	server := Server{
		Config:       config,
		Router:       router,
		RouterGroups: make(map[string]*gin.RouterGroup),
	}

	return server
}
