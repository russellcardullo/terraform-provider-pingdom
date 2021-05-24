package pingdom

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"solarwinds_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"solarwinds_passwd": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"solarwinds_org_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pingdom_check":       resourcePingdomCheck(),
			"pingdom_team":        resourcePingdomTeam(),
			"pingdom_contact":     resourcePingdomContact(),
			"pingdom_integration": resourcePingdomIntegration(),
			"pingdom_maintenance": resourcePingdomMaintenance(),
			"pingdom_user":        resourceSolarwindsUser(),
			"pingdom_occurrence":  resourcePingdomOccurrences(),
			"pingdom_tms_check":   resourcePingdomTmsCheck(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pingdom_contact":      dataSourcePingdomContact(),
			"pingdom_contacts":     dataSourcePingdomContacts(),
			"pingdom_team":         dataSourcePingdomTeam(),
			"pingdom_teams":        dataSourcePingdomTeams(),
			"pingdom_integration":  dataSourcePingdomIntegration(),
			"pingdom_integrations": dataSourcePingdomIntegrations(),
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
