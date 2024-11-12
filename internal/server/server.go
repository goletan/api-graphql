// /api-graphql/internal/server/server.go
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/goletan/api-graphql/internal/schemas"
	config "github.com/goletan/config/pkg"
	observability "github.com/goletan/observability/pkg"
	security "github.com/goletan/security/pkg"
	services "github.com/goletan/services/pkg"
	"github.com/graphql-go/handler"
	"go.uber.org/zap"
)

// GraphQLServer is an enhanced GraphQL server that implements the Service interface.
type GraphQLServer struct {
	server         *http.Server
	name           string
	observability  *observability.Observability
	securityModule *security.Security
	useTLS         bool
}

// GraphQLConfig holds the configuration for the GraphQL server.
type GraphQLConfig struct {
	Address      string `mapstructure:"address"`
	UseTLS       bool   `mapstructure:"use_tls"`
	CertFilePath string `mapstructure:"cert_file_path"`
	KeyFilePath  string `mapstructure:"key_file_path"`
}

var cfg *GraphQLConfig

// NewGraphQLServer creates a new instance of the GraphQLServer.
func NewGraphQLServer(obs *observability.Observability) services.Service {
	// Load configuration
	cfg = &GraphQLConfig{}
	err := config.LoadConfig("GraphQL", cfg, obs.Logger)
	if err != nil {
		obs.Logger.Warn("Failed to load GraphQL configuration, using defaults", zap.Error(err))
	}

	// Set up GraphQL schema
	schema, err := schemas.InitSchema()
	if err != nil {
		obs.Logger.Fatal("Failed to initialize GraphQL schema", zap.Error(err))
	}

	// Create GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   schema,
		GraphiQL: true,
	})

	// Load and configure security module (mTLS)
	securityModule, secErr := security.NewSecurity(obs.Logger)
	if secErr != nil {
		obs.Logger.Error("Failed to initialize security module", zap.Error(secErr))
	}

	// Configure TLS
	tlsConfig, err := securityModule.CertLoader.LoadServerTLSConfig()
	if err != nil {
		obs.Logger.Warn("Failed to configure mTLS, proceeding without TLS", zap.Error(err))
	}

	// Initialize HTTP server
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         tlsConfig,
	}

	// Create the GraphQLServer instance
	return &GraphQLServer{
		server:         server,
		name:           "GraphQL Server",
		observability:  obs,
		securityModule: securityModule,
		useTLS:         cfg.UseTLS,
	}
}

// Name returns the service name.
func (s *GraphQLServer) Name() string {
	return s.name
}

// Initialize performs any initialization tasks needed by the service.
func (s *GraphQLServer) Initialize() error {
	s.observability.Logger.Info("Initializing GraphQL server", zap.String("service", s.name))
	return nil
}

// Start starts the GraphQL server.
func (s *GraphQLServer) Start() error {
	go func() {
		s.observability.Logger.Info("Starting GraphQL server", zap.String("address", s.server.Addr))
		var err error
		if s.useTLS {
			err = s.server.ListenAndServeTLS(cfg.CertFilePath, cfg.KeyFilePath)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.observability.Logger.Error("Failed to start GraphQL server", zap.Error(err))
		}
	}()
	return nil
}

// Stop gracefully stops the GraphQL server.
func (s *GraphQLServer) Stop() error {
	s.observability.Logger.Info("Stopping GraphQL server", zap.String("service", s.name))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
