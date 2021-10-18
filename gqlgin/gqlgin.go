package gqlgin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ubgo/gofm/ginserver"
	"github.com/ubgo/gqlgenfn"
)

type Gql struct {
	Config
}

type Config struct {
	Server    ginserver.Server
	GqlServer *handler.Server
	// Resolver *resolverfn.Resolver
}

// Defining the Graphql handler
func graphqlHandler(gserver *handler.Server) gin.HandlerFunc {
	// https://github.com/99designs/gqlgen/blob/master/docs/content/reference/introspection.md
	gserver.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {

		isAllowedPlayground := IsPlaygroundAllwedForContext(ctx)
		fmt.Println("isAllowedPlayground", isAllowedPlayground)
		if !isAllowedPlayground {
			graphql.GetOperationContext(ctx).DisableIntrospection = true
		}

		return next(ctx)
	})

	return func(c *gin.Context) {
		gserver.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func IsPlaygroundAllwedForContext(ctx context.Context) bool {
	gc, err := gqlgenfn.GinContextFromContext(ctx)
	if err != nil {
		return false
	}

	referer := gc.Request.Header.Get("Referer")
	if len(referer) == 0 {
		return false
	}

	u, err := url.Parse(referer)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	key := q.Get("key")
	if len(key) == 0 {
		return false
	}

	if key == viper.GetString("gql_playground_key") {
		return true
	}

	return false
}

func playgroundAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ok bool
		key, ok := c.GetQuery("key")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "API Key required.",
			})
			c.Abort()
			return
		}

		if key != viper.GetString("gql_playground_key") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Wrong key.",
			})
			c.Abort()
			return
		}

		c.Set("allow_playground", true)

		c.Next()
	}
}

func New(config Config) Gql {
	config.Server.Router.POST("/query", graphqlHandler(config.GqlServer))
	config.Server.Router.GET("/gql", playgroundAccessMiddleware(), playgroundHandler())
	gqlgin := Gql{
		Config: config,
	}

	return gqlgin
}
