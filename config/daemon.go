package config

/* daemon config */
type DConfig struct {
	DebugMode 			bool		`yaml:"debug_mode,omitempty"`   
	SocketPath  		string		`yaml:"socket_path,omitempty"`
	MaxConnQueueLen 	int 		`yaml:"max_conn_queue_len,omitempty"`	
	MaxConcurrentConn	int 		`yaml:"max_conncurrent_conn,omitempty"`	
}

/* normalization function */
func (d *DConfig) Normalize() error {
	
	/* 
		debug_mode is false by default
		daemon will run on production mode by default
	*/

	/* if SocketPath is empty, use "/var/run/laclm-daemon.sock" as default */
	if d.SocketPath == "" {
		d.SocketPath = "/var/run/laclm-daemon.sock"
	}

	/* if maximum connection queue length is not set or less than equal to 0, set it to 500 */
	if d.MaxConnQueueLen <= 0 {
		/* set queue length to 500 */
		d.MaxConnQueueLen = 500
	}

	/* if maximum concurrent connection is not set or less than equal to 0, set it to 5 */
	if d.MaxConcurrentConn <= 0 {
		/* set max concurrent connections to 5 */
		d.MaxConcurrentConn = 5
	}

	return nil
}
