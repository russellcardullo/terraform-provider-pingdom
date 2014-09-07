package pingdom

import (
	"log"
	"os"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	APIKey   string `mapstructure:"api_key"`
}

// Client() returns a new client for accessing pingdom.
//
func (c *Config) Client() (*pingdom.Client, error) {

	if v := os.Getenv("PINGDOM_USER"); v != "" {
		c.User = v
	}
	if v := os.Getenv("PINGDOM_PASSWORD"); v != "" {
		c.Password = v
	}
	if v := os.Getenv("PINGDOM_API_KEY"); v != "" {
		c.APIKey = v
	}

	client := pingdom.NewClient(c.User, c.Password, c.APIKey)

	log.Printf("[INFO] Pingdom Client configured for user: %s", c.User)

	return client, nil
}
