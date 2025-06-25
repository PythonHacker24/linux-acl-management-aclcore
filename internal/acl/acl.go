package acl

import (
	"encoding/json"
	"net"
	"os/exec"
	"sync"
)

/* maintains locks on file which are actively under ACL modifications */
var pathLocks sync.Map

/* locks a given file path */
func getPathLock(path string) *sync.Mutex {
	mtx, _ := pathLocks.LoadOrStore(path, &sync.Mutex{})
	return mtx.(*sync.Mutex)
}

/* handle connection for ACL requests */
func HandleConnection(conn net.Conn) error {
	/* close the connection before function ends */
	defer conn.Close()

	/* create 1 KB buffer (MAKE THIS MODIFIABLE) */
	buf := make([]byte, 1024)
	data, err := conn.Read(buf)
	if err != nil {
		return err
	}

	/* unmarshal the JSON request */
	var req ACLRequest
	if err := json.Unmarshal(buf[:data], &req); err != nil {
		sendResponse(conn, false, "Invalid JSON")
		return err
	}

	/* lock the file path for thread safety (ensure unlock even on panic) */
	lock := getPathLock(req.Path)
	lock.Lock()
	defer lock.Unlock()

	/* execute the ACL modifications with acl commands */
	var cmd *exec.Cmd
	switch req.Action {
	case "add", "modify":
		cmd = exec.Command("setfacl", "-m", req.Entry, req.Path)
	case "remove":
		cmd = exec.Command("setfacl", "-x", req.Entry, req.Path)
	default:
		sendResponse(conn, false, "Unsupported action: "+req.Action)
		return nil
	}

	/* retrive the output and send it via connection */
	output, err := cmd.CombinedOutput()
	if err != nil {
		sendResponse(conn, false, string(output))
	} else {
		sendResponse(conn, true, string(output))
	}

	/* no errors, return nil */
	return nil
}

/* send structured data over the socket */
func sendResponse(conn net.Conn, success bool, msg string) {
	resp := ACLResponse{Success: success, Message: msg}
	data, _ := json.Marshal(resp)
	conn.Write(data)
}
