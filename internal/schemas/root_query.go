// /api-graphql/internal/schemas/root_query.go
package schemas

import (
	"github.com/goletan/api-graphql/internal/handlers"
	"github.com/graphql-go/graphql"
)

// DefineRootQuery initializes the root query for GraphQL.
func DefineRootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type:        graphql.String,
				Description: "Get the status of the GraphQL API",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Use the GetStatus function to get the status information
					return handlers.GetStatus(), nil
				},
			},
			"uptime": &graphql.Field{
				Type:        graphql.String,
				Description: "Get the uptime of the server",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return handlers.GetUptime(), nil
				},
			},
			"version": &graphql.Field{
				Type:        graphql.String,
				Description: "Get the current version of the GraphQL API",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return handlers.GetVersion(), nil
				},
			},
			"serverTime": &graphql.Field{
				Type:        graphql.String,
				Description: "Get the current server time",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return handlers.GetServerTime(), nil
				},
			},
		},
	})
}
