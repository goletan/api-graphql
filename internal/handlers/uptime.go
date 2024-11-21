// /api-graphql/internal/handlers/uptime.go
package handlers

import (
	"time"
)

var startTime = time.Now()

// GetUptime provides the uptime of the server in seconds.
func GetUptime() string {
	uptime := time.Since(startTime)
	return uptime.String()
}
