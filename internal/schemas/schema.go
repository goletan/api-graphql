// /api-graphql/internal/schemas/schema.go
package schemas

import (
	"github.com/graphql-go/graphql"
)

// InitSchema initializes the GraphQL schema with queries and mutations.
func InitSchema() (*graphql.Schema, error) {
	rootQuery := DefineRootQuery()

	// We can define root mutations here if needed later
	schemaConfig := graphql.SchemaConfig{
		Query: rootQuery,
		// Mutation: rootMutation, // Uncomment when mutations are defined
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}
