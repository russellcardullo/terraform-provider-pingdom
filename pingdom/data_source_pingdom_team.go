package pingdom

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func dataSourcePingdomTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePingdomTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePingdomTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	name := d.Get("name").(string)
	teams, err := client.Teams.List()
	log.Printf("==== teams : %v", teams)
	if err != nil {
		return fmt.Errorf("Error retrieving team: %s", err)
	}
	var found *pingdom.TeamResponse
	for _, team := range teams {
		if team.Name == name {
			log.Printf("Team: %v", team)
			found = &team
			break
		}
	}
	if found == nil {
		return fmt.Errorf("User '%s' not found", name)
	}
	if err = d.Set("name", found.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	var memberIds []int
	for _, member := range found.Members {
		memberIds = append(memberIds, member.ID)
	}

	if err = d.Set("member_ids", memberIds); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", found.ID))
	return nil
}
