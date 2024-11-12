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
					return handlers.StatusHandler(p)
				},
			},
			// Add other queries here as needed.
		},
	})
}
