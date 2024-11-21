// /api-graphql/internal/handlers/version.go
package handlers

// Version represents the current version of the API.
const Version = "v1.0.0"

// GetVersion provides the version information of the API.
func GetVersion() string {
	return Version
}
