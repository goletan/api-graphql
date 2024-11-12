// /api-graphql/pkg/graphql.go
package graphql

import (
	observability "github.com/goletan/observability/pkg"
	services "github.com/goletan/services/pkg"
	"go.uber.org/zap"

	"github.com/goletan/api-graphql/internal/server"
)

// NewGraphQLService creates a new GraphQL service that implements the Goletan service interface.
func NewGraphQLService(logger *zap.Logger) services.Service {
	obs, err := observability.NewObserver()
	if err != nil {
		logger.Error("Failed to initialize observability", zap.Error(err))
	}

	return server.NewGraphQLServer(obs)
}
