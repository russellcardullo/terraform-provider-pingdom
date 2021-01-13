package pingdom

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pingdom_check":              resourcePingdomCheck(),
			"pingdom_team":               resourcePingdomTeam(),
			"pingdom_contact":            resourcePingdomContact(),
			"pingdom_maintenance_window": resourcePingdomMaintenanceWindow(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pingdom_contact": dataSourcePingdomContact(),
			"pingdom_team":    dataSourcePingdomTeam(),
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
