// /api-graphql/internal/handlers/health.go
package handlers

import (
	"github.com/graphql-go/graphql"
)

// HealthHandler returns a basic health status.
func HealthHandler(p graphql.ResolveParams) (interface{}, error) {
	return "Service is healthy", nil
}
