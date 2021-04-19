package pingdom

import (
	"errors"
	"github.com/nordcloud/go-pingdom/solarwinds"
	"log"
	"os"

	"github.com/nordcloud/go-pingdom/pingdom"
	"github.com/nordcloud/go-pingdom/pingdomext"
)

// Config respresents the client configuration
type Config struct {
	APIToken           string `mapstructure:"api_token"`
	SolarwindsUser     string `mapstructure:"solarwinds_user"`
	SolarwindsPassword string `mapstructure:"solarwinds_password"`
}

type Clients struct {
	Pingdom    *pingdom.Client
	PingdomExt *pingdomext.Client
	Solarwinds *solarwinds.Client
}

func (c *Config) Client() (*Clients, error) {
	pingdomClient, err := c.pingdomClient()
	if err != nil {
		return nil, err
	}
	if user := os.Getenv("SOLARWINDS_USER"); user != "" {
		if password := os.Getenv("SOLARWINDS_PASSWD"); password != "" {
			c.SolarwindsUser = user
			c.SolarwindsPassword = password
		} else {
			return nil, errors.New("user and password must be present together")
		}
	}
	solarwindsClient, err := solarwinds.NewClient(solarwinds.ClientConfig{
		Username: c.SolarwindsUser,
		Password: c.SolarwindsPassword,
	})
	if err != nil {
		return nil, err
	}
	err = solarwindsClient.Init()
	if err != nil {
		return nil, err
	}

	pingdomClientExt, err := pingdomext.NewClientWithConfig(pingdomext.ClientConfig{
		Username: c.SolarwindsUser,
		Password: c.SolarwindsPassword,
	})

	if err != nil {
		return nil, err
	}

	return &Clients{
		Pingdom:    pingdomClient,
		PingdomExt: pingdomClientExt,
		Solarwinds: solarwindsClient,
	}, nil
}

// Client returns a new client for accessing pingdom.
//
func (c *Config) pingdomClient() (*pingdom.Client, error) {
	if v := os.Getenv("PINGDOM_API_TOKEN"); v != "" {
		c.APIToken = v
	}

	client, _ := pingdom.NewClientWithConfig(pingdom.ClientConfig{APIToken: c.APIToken})

	log.Printf("[INFO] Pingdom Client configured.")

	return client, nil
}
