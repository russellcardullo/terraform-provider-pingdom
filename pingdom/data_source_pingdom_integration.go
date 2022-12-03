package pingdom

import (
	"context"
	"fmt"
	"log"

	"github.com/DrFaust92/go-pingdom/pingdomext"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePingdomIntegration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePingdomIntegrationRead,

		Schema: map[string]*schema.Schema{
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePingdomIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).PingdomExt
	name := d.Get("name").(string)
	integrations, err := client.Integrations.List()
	log.Printf("[DEBUG] integrations : %v", integrations)
	if err != nil {
		return diag.Errorf("Error retrieving team: %s", err)
	}
	var found *pingdomext.IntegrationGetResponse
	for _, integration := range integrations {
		if integration.UserData["name"] == name {
			log.Printf("Integration: %v", integration)
			found = &integration
			break
		}
	}
	if found == nil {
		return diag.Errorf("Integration '%s' not found", name)
	}

	if err := d.Set("provider_name", found.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", found.ActivatedAt != 0); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", found.UserData["name"]); err != nil {
		return diag.FromErr(err)
	}
	if found.Name == WEBHOOK {
		if err := d.Set("url", found.UserData["url"]); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%d", found.ID))
	return nil
}
