package pingdom

import (
	"log"
	"os"

	"github.com/nordcloud/go-pingdom/pingdom"
)

// Config respresents the client configuration
type Config struct {
	APIToken string `mapstructure:"api_token"`
}

// Client returns a new client for accessing pingdom.
//
func (c *Config) Client() (*pingdom.Client, error) {

	if v := os.Getenv("PINGDOM_API_TOKEN"); v != "" {
		c.APIToken = v
	}

	client, _ := pingdom.NewClientWithConfig(pingdom.ClientConfig{APIToken: c.APIToken})

	log.Printf("[INFO] Pingdom Client configured.")

	return client, nil
}
