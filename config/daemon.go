package config

/* daemon config */
type DConfig struct {
	DebugMode 	bool		`yaml:"debug_mode,omitempty"`   
	SocketPath  string		`yaml:"socket_path,omitempty"`
	MaxConnPool int 		`yaml:"max_conn_pool,omitempty"`	
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

	/* if maximum connection pool is not set or less than equal to 0, set it to 1 */
	if d.MaxConnPool <= 0 {
		/* perform single connection operation */
		d.MaxConnPool = 1
	}

	return nil
}
