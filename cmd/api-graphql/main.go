// /api-graphql/cmd/api-graphql/main.go
package main

import (
	"github.com/goletan/api-graphql/internal/server"
	config "github.com/goletan/config/pkg"
	obsConfig "github.com/goletan/observability/config"
	observability "github.com/goletan/observability/pkg"
	"go.uber.org/zap"
)

func main() {
	var obsCfg *obsConfig.ObservabilityConfig
	err := config.LoadConfig("Observability", &obsCfg, nil)
	if err != nil {
		panic(err)
	}

	// Initialize observability with the configuration
	obs, err := observability.NewObserver(obsCfg)
	if err != nil
		obs.Logger.Error("Failed to initialize observability: %v", zap.Error(err))
	}

	// Create a new GraphQL server instance
	graphqlServer := server.NewGraphQLServer(obs)

	// Initialize the GraphQL server
	if err := graphqlServer.Initialize(); err != nil {
		obs.Logger.Error("Failed to initialize GraphQL server: %v", zap.Error(err))
	}

	// Start the service
	if err := graphqlServer.Start(); err != nil {
		panic(err)
	}
}
