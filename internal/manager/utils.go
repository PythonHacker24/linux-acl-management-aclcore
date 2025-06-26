package manager

import (
	"encoding/json"
	"net"
)

/* sends a structured error response through given socket */
func errorResponse(conn net.Conn, message string) {
	_ = json.NewEncoder(conn).Encode(map[string]string{"error": message})
	_ = conn.Close()
}
