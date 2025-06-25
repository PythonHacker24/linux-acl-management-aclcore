package manager

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/PythonHacker24/linux-acl-management-aclcore/config"
	"github.com/PythonHacker24/linux-acl-management-aclcore/internal/acl"
	"go.uber.org/zap"
)

/* creates a new ACL server */
func NewACLServer(path string, errCh chan error) *ACLServer {
	return &ACLServer{
		socketPath: path,
		errCh:      errCh,
	}
}

/* starts the ACL server */
func (s *ACLServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	if err := os.RemoveAll(s.socketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket connection: %w", err)
	}

	s.listener = listener

	sem := make(chan struct{}, config.COREDConfig.DConfig.MaxConnPool)

	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		zap.L().Info("Shutting down ACL server")
		s.listener.Close()
		close(done)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				<-done
				fmt.Println("ACL server shutdown complete.")
				return nil
			default:
				fmt.Println("Accept error:", err)
				continue
			}
		}

		select {
		/* acquire a slot */
		case sem <- struct{}{}:
			/* add process the waitgroup */
			wg.Add(1)

			/* handle connection asynchronously */
			go func(c net.Conn) {
				defer wg.Done()

				/* release the slot */
				defer func() { <-sem }()

				acl.HandleConnection(c)
			}(conn)
		default:
			/* no slot available, let client know */
			fmt.Println(conn, `{"error": "server overloaded"}`)
			zap.L().Info("rejecting connection")
			conn.Close()
		}
	}
}
