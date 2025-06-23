package config

import "fmt"

var (
	COREDConfig CoreDConfig
)

/* config struct for aclapi */
type CoreDConfig struct {
	DConfig 	DConfig 	`yaml:"daemon,omitempty"`
	Logging		Logging 	`yaml:"logs,omitempty"`
}

/* complete config normalizer function */
func (c *CoreDConfig) Normalize() error {
	
	if err := c.DConfig.Normalize(); err != nil {
		return fmt.Errorf("daemon configuration error: %w", err)
	}

	if err := c.Logging.Normalize(); err != nil {
		return fmt.Errorf("logging configuration error: %w", err)
	}

	return nil 
}
