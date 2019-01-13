package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/smook1980/medialocker/app"
	"github.com/smook1980/medialocker/app/server/mutations"
	"github.com/smook1980/medialocker/app/server/queries"
)

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    queries.RootQuery,
	Mutation: mutations.RootMutation,
})

func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func graphqlHandler() http.Handler {
	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
}

func Start(ctx app.Context) {
	r := gin.Default()
	r.Any("/graphql", gin.WrapH(graphqlHandler()))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	listen := fmt.Sprintf("%s:%d", ctx.HTTPServerHost(), ctx.HTTPServerPort())
	r.Run(listen)
}
