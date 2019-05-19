package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomTeamCreate,
		Read:   resourcePingdomTeamRead,
		Update: resourcePingdomTeamUpdate,
		Delete: resourcePingdomTeamDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePingdomTeamImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

type commonTeamParams struct {
	Name string
}

func teamForResource(d *schema.ResourceData) (*pingdom.TeamData, error) {
	teamParams := commonTeamParams{}

	// required
	if v, ok := d.GetOk("name"); ok {
		teamParams.Name = v.(string)
	}

	return &pingdom.TeamData{
		Name: teamParams.Name,
	}, nil
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

	d.SetId(result.ID)
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
		if team.ID == d.Id() {
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

	d.Set("name", team.Name)
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

func resourcePingdomTeamImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Importing key using ADDR ID %d", id)

	return []*schema.ResourceData{d}, nil
}
