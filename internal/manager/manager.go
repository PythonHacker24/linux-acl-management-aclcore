package manager

import (
	"context"
	"fmt"
	"net"

	"github.com/PythonHacker24/linux-acl-management-aclcore/config"
	"github.com/PythonHacker24/linux-acl-management-aclcore/internal/acl"
)

/* 
	TODO: follow the case select pattern for context handling here
*/
func ConnPool(ctx context.Context, errCh chan<-error) error {
	/* dial into specified Unix socket */
	conn, err := net.Dial("unix", config.COREDConfig.DConfig.SocketPath)
	if err != nil {
		return fmt.Errorf("failed to dial Unix socket: %w", err)
	}

	/* close the connection before function ends */
	defer conn.Close()

	for {
		go acl.HandleConnection(conn)
	}

	ctx.Done()

	return nil
}
