// /api-graphql/internal/server/server.go
package server

import (
	"context"
	"crypto/tls"
	"fmt"
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
	config         *GraphQLConfig
	useTLS         bool
}

// GraphQLConfig holds the configuration for the GraphQL server.
type GraphQLConfig struct {
	Address      string `mapstructure:"address"`
	UseTLS       bool   `mapstructure:"use_tls"`
	CertFilePath string `mapstructure:"cert_file_path"`
	KeyFilePath  string `mapstructure:"key_file_path"`
}

// NewGraphQLServer creates a new instance of the GraphQLServer.
func NewGraphQLServer(obs *observability.Observability) services.Service {
	config := loadGraphQLConfig(obs)

	// Setup Security
	securityModule := setupSecurityModule(obs)
	if securityModule == nil {
		obs.Logger.Fatal("Security module could not be initialized")
	}

	// Setup server
	server := setupServer(obs, securityModule, config)

	// Create the GraphQLServer instance
	return &GraphQLServer{
		server:         server,
		name:           "GraphQL Server",
		observability:  obs,
		securityModule: securityModule,
		config:         config,
		useTLS:         config.UseTLS,
	}
}

// Name returns the service name.
func (s *GraphQLServer) Name() string {
	return s.name
}

// Initialize performs any initialization tasks needed by the service.
func (s *GraphQLServer) Initialize() error {
	s.observability.Logger.Info("Initializing GraphQL server", zap.String("service", s.name))

	// Setup Security Module again to ensure all security setups are reloaded if needed
	s.securityModule = setupSecurityModule(s.observability)
	if s.securityModule == nil {
		s.observability.Logger.Error("Failed to initialize security module during server initialization")
		return fmt.Errorf("security module initialization failed")
	}

	return nil
}

// Start starts the GraphQL server.
func (s *GraphQLServer) Start() error {
	s.observability.Logger.Info("Starting GraphQL server", zap.String("address", s.server.Addr))
	var err error
	if s.useTLS {
		err = s.server.ListenAndServeTLS(s.config.CertFilePath, s.config.KeyFilePath)
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		s.observability.Logger.Error("Failed to start GraphQL server", zap.Error(err))
		return err
	}
	if err == http.ErrServerClosed {
		s.observability.Logger.Info("GraphQL server has been stopped gracefully")
	}
	return nil
}

// Stop gracefully stops the GraphQL server.
func (s *GraphQLServer) Stop() error {
	s.observability.Logger.Info("Stopping GraphQL server", zap.String("service", s.name))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func loadGraphQLConfig(obs *observability.Observability) *GraphQLConfig {
	cfg := &GraphQLConfig{
		Address: ":8081", // Set default values here if needed
		UseTLS:  false,
	}
	err := config.LoadConfig("GraphQL", cfg, obs.Logger)
	if err != nil {
		obs.Logger.Warn("Failed to load GraphQL configuration, using defaults", zap.Error(err))
	}
	return cfg
}

func setupSecurityModule(obs *observability.Observability) *security.Security {
	securityModule, err := security.NewSecurity(obs.Logger)
	if err != nil {
		obs.Logger.Error("Failed to initialize security module", zap.Error(err))
		return nil
	}
	return securityModule
}

func setupServer(obs *observability.Observability, securityModule *security.Security, cfg *GraphQLConfig) *http.Server {
	var tlsConfig *tls.Config
	var err error

	// Configure TLS
	if securityModule != nil {
		tlsConfig, err = securityModule.CertLoader.LoadServerTLSConfig()
		if err != nil {
			obs.Logger.Warn("Failed to configure mTLS, proceeding without TLS", zap.Error(err))
		}
	} else {
		obs.Logger.Warn("Security module is nil, proceeding without TLS")
	}

	// Set up GraphQL schema
	schema, err := schemas.InitSchema()
	if err != nil {
		obs.Logger.Fatal("Failed to initialize GraphQL schema", zap.Error(err))
	}

	obs.Logger.Sugar().Infof("server address: %s", cfg.Address)

	// Create GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   schema,
		GraphiQL: true,
	})

	// Initialize HTTP server
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         tlsConfig,
	}

	return server
}
