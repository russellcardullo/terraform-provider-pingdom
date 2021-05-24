package pingdom

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePingdomIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePingdomIntegrationsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePingdomIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).PingdomExt
	integrations, err := client.Integrations.List()
	log.Printf("[DEBUG] integrations : %v", integrations)
	if err != nil {
		return diag.Errorf("Error retrieving team: %s", err)
	}
	var ids []int
	var names []string
	for _, integration := range integrations {
		ids = append(ids, integration.ID)
		names = append(names, integration.UserData["name"])
	}

	d.SetId(fmt.Sprintf("%d", len(integrations)))
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("names", names); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
