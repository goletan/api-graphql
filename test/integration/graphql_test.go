// /kernel/internal/tests/integration/graphql_test.go
package integration_test

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGraphQLAPIHealth(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}

	query := `{"query": "{ status }"}`
	req, err := http.NewRequest("POST", "https://localhost:9443/graphql", bytes.NewBuffer([]byte(query)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err, "Should be able to call GraphQL endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "GraphQL endpoint should return status OK")
}
