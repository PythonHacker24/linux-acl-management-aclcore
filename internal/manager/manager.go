package manager

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/PythonHacker24/linux-acl-management-aclcore/config"
	"github.com/PythonHacker24/linux-acl-management-aclcore/internal/acl"
)

/*
TODO: follow the case select pattern for context handling here
fix the whole code till it's production ready
*/
func ConnPool(ctx context.Context, errCh chan<- error) error {

	/* declare a semaphore for maximum connections */
	sem := make(chan struct{}, config.COREDConfig.DConfig.MaxConnPool)
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil
		default:
			conn, err := net.Dial("unix", config.COREDConfig.DConfig.SocketPath)
			if err != nil {
				select {
				case errCh <- fmt.Errorf("failed to dial Unix socket: %w", err):
				default:
					/* drop error if channel is full (must be rare) */
				}
				/* give it some rest */
				time.Sleep(time.Second)
				/* continue with the next connection */
				continue
			}

			/* connection is now established */

			/* acquire semaphore slot (or wait until slot is free and stop if ctx is done) */
			select {
			case sem <- struct{}{}:
				/* semaphore slot granted, continue */
			case <-ctx.Done():
				conn.Close()
				return nil
			}

			/* add to waitgroup */
			wg.Add(1)

			/* call handler asynchronously */
			go func(c net.Conn) {
				defer wg.Done()
				defer c.Close()
				defer func() { <-sem }()

				if err := acl.HandleConnection(c); err != nil {
					select {
					case errCh <- fmt.Errorf("connection handler error: %w", err):
					default:
						/* drop if error channel is full (rare case) */
					}
				}
			}(conn)
		}
	}
}
