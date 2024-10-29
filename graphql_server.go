// graphql_server.go
package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/goletan/config"
	"github.com/goletan/security/mtls"
	"github.com/goletan/services"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GraphQLServer is an enhanced GraphQL server that implements the Service interface.
type GraphQLServer struct {
	server *http.Server
	name   string
	useTLS bool
}

// GraphQLConfig holds the configuration for the GraphQL server.
type GraphQLConfig struct {
	Address      string `mapstructure:"address"`
	UseTLS       bool   `mapstructure:"use_tls"`
	CertFilePath string `mapstructure:"cert_file_path"`
	KeyFilePath  string `mapstructure:"key_file_path"`
}

var cfg *GraphQLConfig

// Define a simple root query object.
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "API is running smoothly!", nil
			},
		},
	},
})

// NewGraphQLServer creates a new instance of the GraphQLServer.
func NewGraphQLServer() services.Service {
	cfg = &GraphQLConfig{}
	err := config.LoadConfig("GraphQL", cfg, nil)
	if err != nil {
		fmt.Printf("Warning: failed to load GraphQL configuration, using defaults: %v\n", err)
	}

	schemaConfig := graphql.SchemaConfig{Query: rootQuery}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Errorf("Failed to create GraphQL schema: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
	})

	tlsConfig, err := mtls.ConfigureMTLS()
	if err != nil {
		fmt.Printf("Warning: failed to configure mTLS, proceeding without TLS: %v\n", err)
	}

	return &GraphQLServer{
		server: &http.Server{
			Addr:      cfg.Address,
			Handler:   h,
			TLSConfig: tlsConfig,
		},
		name:   "GraphQL Server",
		useTLS: cfg.UseTLS,
	}
}

// Name returns the service name.
func (s *GraphQLServer) Name() string {
	return s.name
}

// Initialize performs any initialization tasks needed by the service.
func (s *GraphQLServer) Initialize() error {
	log.Printf("Initializing %s", s.name)
	return nil
}

// Start starts the GraphQL server.
func (s *GraphQLServer) Start() error {
	go func() {
		log.Printf("Starting %s on %s", s.name, s.server.Addr)
		var err error
		if s.useTLS {
			err = s.server.ListenAndServeTLS(cfg.CertFilePath, cfg.KeyFilePath)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start %s: %v", s.name, err)
		}
	}()
	return nil
}

// Stop gracefully stops the GraphQL server.
func (s *GraphQLServer) Stop() error {
	log.Printf("Stopping %s", s.name)
	return s.server.Shutdown(context.Background())
}
