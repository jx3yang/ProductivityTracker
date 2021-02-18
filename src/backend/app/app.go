package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jx3yang/ProductivityTracker/src/backend/config"
	db "github.com/jx3yang/ProductivityTracker/src/backend/database"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph"
	"github.com/jx3yang/ProductivityTracker/src/backend/graph/generated"
	coll_handler "github.com/jx3yang/ProductivityTracker/src/backend/handler"
)

const devConfigFile = "dev/config.yml"

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type server interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func makeGinHandler(serv server) gin.HandlerFunc {
	return func(c *gin.Context) {
		serv.ServeHTTP(c.Writer, c.Request)
	}
}

func graphqlHandler() gin.HandlerFunc {
	serv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	return makeGinHandler(serv)
}

func playgroundHandler() gin.HandlerFunc {
	serv := playground.Handler("GraphQL playground", "/query")
	return makeGinHandler(serv)
}

func Run() {
	var configFile *string = new(string)
	if os.Getenv("ENV") == config.Prod {
		configFile = nil
	} else {
		*configFile = devConfigFile
	}
	fmt.Println(*configFile)
	configuration, err := config.GetConfig(configFile)
	handleErr(err)

	conn, err := db.ConnectionFromConfig(configuration)
	handleErr(err)

	database := conn.InitDatabase(db.PT)
	err = coll_handler.InitHandlers(database)
	handleErr(err)

	r := gin.Default()
	r.GET("/", playgroundHandler())
	r.POST("/query", graphqlHandler())
	r.Run(":" + configuration.Port)
}