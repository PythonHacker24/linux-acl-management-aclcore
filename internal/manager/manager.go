package manager

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

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
func (s *ACLServer) Start(ctx context.Context, wg *sync.WaitGroup, maxQueue, maxConcurrent int) error {
	if err := os.RemoveAll(s.socketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket connection: %w", err)
	}

	s.listener = listener
	s.queueChan = make(chan net.Conn, maxQueue)

	/* worker pool */
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case conn, ok := <-s.queueChan:
					if !ok {
						return
					}
					acl.HandleConnection(conn)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		zap.L().Info("Shutting down ACL server")
		s.listener.Close()
		close(s.queueChan)
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				zap.L().Info("ACL server shutdown complete")
				return nil
			default:
				zap.L().Error("Accept error:",
					zap.Error(err),
				)
				continue
			}
		}

		select {
		case s.queueChan <- conn:
			/* connection enqueued successfully */
		default:
			/* connection queue is full */
			zap.L().Warn("Request dropped: queue full")
			errorResponse(conn, "server overloaded, try later")
		}
	}
}
