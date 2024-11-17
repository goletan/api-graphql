// /api-graphql/pkg/graphql.go
package graphql

import (
	observability "github.com/goletan/observability/pkg"
	services "github.com/goletan/services/pkg"

	"github.com/goletan/api-graphql/internal/server"
)

// NewGraphQLService creates a new GraphQL service that implements the Goletan service interface.
func NewGraphQLService(obs *observability.Observability) services.Service {
	return server.NewGraphQLServer(obs)
}
