package queries

import (
	"github.com/graphql-go/graphql"
	fields "github.com/smook1980/medialocker/app/server/queries/fields"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"getNotTodos": fields.GetNotTodos,
	},
})
