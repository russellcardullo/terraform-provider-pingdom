package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/nordcloud/go-pingdom/pingdom"
)

func resourcePingdomTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomTeamCreate,
		Read:   resourcePingdomTeamRead,
		Update: resourcePingdomTeamUpdate,
		Delete: resourcePingdomTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func teamForResource(d *schema.ResourceData) (*pingdom.Team, error) {
	team := pingdom.Team{}

	// required
	if v, ok := d.GetOk("name"); ok {
		team.Name = v.(string)
	}

	if v, ok := d.GetOk("member_ids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))

		}
		team.MemberIDs = intSlice
	}

	return &team, nil
}

func resourcePingdomTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	team, err := teamForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Team create configuration: %#v", d.Get("name"))
	result, err := client.Teams.Create(team)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return nil
}

func resourcePingdomTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	teams, err := client.Teams.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of teams: %s", err)
	}
	exists := false
	for _, team := range teams {
		if strconv.Itoa(team.ID) == d.Id() {
			exists = true
			break
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	team, err := client.Teams.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving team: %s", err)
	}

	if err := d.Set("name", team.Name); err != nil {
		return err
	}

	memberids := schema.NewSet(
		func(memberId interface{}) int { return memberId.(int) },
		[]interface{}{},
	)
	for _, member := range team.Members {
		memberids.Add(member.ID)
	}
	if err := d.Set("member_ids", memberids); err != nil {
		return err
	}

	return nil
}

func resourcePingdomTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	team, err := teamForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Team update configuration: %#v", d.Get("name"))

	if _, err = client.Teams.Update(id, team); err != nil {
		return fmt.Errorf("Error updating team: %s", err)
	}
	return nil
}

func resourcePingdomTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	if _, err = client.Teams.Delete(id); err != nil {
		return fmt.Errorf("Error deleting team: %s", err)
	}

	return nil
}
