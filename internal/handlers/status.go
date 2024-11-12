// /api-graphql/internal/handlers/status.go
package handlers

import (
	"github.com/graphql-go/graphql"
)

// StatusHandler handles the status field.
func StatusHandler(p graphql.ResolveParams) (interface{}, error) {
	return "GraphQL API is running smoothly!", nil
}
