package pingdom

import (
	"fmt"
	"strconv"
	"strings"

   "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func dataSourcePingdomMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePingdomMaintenanceWindowRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"recurrence_type": {
				Type:     schema.TypeString,
				Required: true,
				// Default:  "none",
			},
			"repeat_every": {
				Type:     schema.TypeInt,
				Required: true,
				// Default:  0,
			},
			"effective_to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"uptimeids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"tmsids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePingdomMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	id := d.Get("id").(int)
	window, err := client.Maintenances.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving maintenance window: %s", err)
	}

	if err = d.Set("description", window.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("from", window.From); err != nil {
		return fmt.Errorf("Error setting from : %s", err)
	}
	if err = d.Set("to", window.To); err != nil {
		return fmt.Errorf("Error setting to: %s", err)
	}
	if err = d.Set("recurence_type", window.RecurrenceType); err != nil {
		return fmt.Errorf("Error setting recurence_type: %s", err)
	}
	if err = d.Set("repeat_every", window.RepeatEvery); err != nil {
		return fmt.Errorf("Error setting repeat_every: %s", err)
	}
	if err =   d.Set("effective_to", window.EffectiveTo); err != nil {
		return fmt.Errorf("Error setting effective_to: %s", err)
	}

	uptimeids := make([]string, len(window.Checks.Uptime))
	for i, x := range window.Checks.Uptime {
		uptimeids[i] = strconv.Itoa(x)
	}
	if err = d.Set("uptimeids", strings.Join(uptimeids, ",")); err != nil {
		return fmt.Errorf("Error setting uptimeids: %s", err)
	}

	tmids := make([]string, len(window.Checks.Tms))
	for i, x := range window.Checks.Tms {
		tmids[i] = strconv.Itoa(x)
	}
	if err = d.Set("tmsids", strings.Join(tmids, ",")); err != nil {
		return fmt.Errorf("Error setting tmsids: %s", err)
	}

	return nil
}
