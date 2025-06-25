package manager

import (
	"net"
)

/* acl server struct */
type ACLServer struct {
	socketPath string
	errCh      chan error
	listener   net.Listener
}
