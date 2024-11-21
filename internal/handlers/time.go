// /api-graphql/internal/handlers/time.go
package handlers

import "time"

// GetServerTime returns the current server time.
func GetServerTime() string {
	return time.Now().Format(time.RFC1123)
}
