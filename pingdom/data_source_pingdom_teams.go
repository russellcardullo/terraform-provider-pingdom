package pingdom

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePingdomTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePingdomTeamsRead,

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

func dataSourcePingdomTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom
	teams, err := client.Teams.List()
	if err != nil {
		return diag.Errorf("Error retrieving teams: %s", err)
	}

	var ids = make([]int, 0, len(teams))
	var names = make([]string, 0, len(teams))
	for _, team := range teams {
		ids = append(ids, team.ID)
		names = append(names, team.Name)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("names", names); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
