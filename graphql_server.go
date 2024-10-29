// graphql_server.go
package graphql

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/goletan/services"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GraphQLServer is an enhanced GraphQL server that implements the Service interface.
type GraphQLServer struct {
	server *http.Server
	name   string
}

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
	schemaConfig := graphql.SchemaConfig{Query: rootQuery}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Errorf("Failed to create GraphQL schema: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		GraphiQL: true,
	})

	return &GraphQLServer{
		server: &http.Server{
			Addr:    ":8081",
			Handler: h,
		},
		name: "GraphQL Server",
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
		log.Printf("Starting %s on :8081", s.name)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
