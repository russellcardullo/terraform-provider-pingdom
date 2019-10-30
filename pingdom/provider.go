package pingdom

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"account_email": {
				Type:      schema.TypeString,
				Default:   "",
				Optional:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pingdom_check":   resourcePingdomCheck(),
			"pingdom_team":    resourcePingdomTeam(),
			"pingdom_user":    resourcePingdomUser(),
			"pingdom_contact": resourcePingdomContact(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pingdom_user": dataSourcePingdomUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}

	log.Println("[INFO] Initializing Pingdom client")
	return config.Client()
}
