package pingdom

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func dataSourcePingdomMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePingdomUserRead,
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

func dataSourcePingdomMainteannceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	id := d.Get("id").(int)
	window, err := client.Maintenances.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving maintenance window: %s", err)
	}

	d.Set("descrition", window.Description)
	d.Set("from", window.From)
	d.Set("to", window.To)
	d.Set("recurence_type", window.RecurrenceType)
	d.Set("repeat_every", window.RepeatEvery)
	d.Set("effective_to", window.EffectiveTo)

	uptimeids := make([]string, len(window.Checks.Uptime))
	for i, x := range window.Checks.Uptime {
		uptimeids[i] = strconv.Itoa(x)
	}
	d.Set("uptimeids", strings.Join(uptimeids, ","))

	tmids := make([]string, len(window.Checks.Tms))
	for i, x := range window.Checks.Tms {
		tmids[i] = strconv.Itoa(x)
	}
	d.Set("tmids", strings.Join(tmids, ","))

	return nil
}
