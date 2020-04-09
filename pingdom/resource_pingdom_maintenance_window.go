package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomUserCreate,
		Read:   resourcePingdomUserRead,
		Update: resourcePingdomUserUpdate,
		Delete: resourcePingdomUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required: false,
				Default:  "none",
			},
			"repeat_every": {
				Type:     schema.TypeInt,
				Required: false,
				Default:  0,
			},
			"effective_to": {
				Type:     schema.TypeInt,
				Required: false,
			},
			"uptimeids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"tmsids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

type commonMaintenanceWindowParams struct {
	Description string
	From        int64
	To          int64
}

func maintenanceWindowForResource(d *schema.ResourceData) (*pingdom.MaintenanceWindow, error) {
	windowParams := commonMaintenanceWindowParams{}

	// required
	if v, ok := d.GetOk("description"); ok {
		windowParams.Description = v.(string)
	}

	return &pingdom.MaintenanceWindow{
		Description: windowParams.Description,
	}, nil
}

func resourcePingdomMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	window, err := client.Maintenances.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving maintenance window: %s", err)
	}

	d.Set("descrition", window.Description)
	d.Set("from", window.From)
	d.Set("to", window.To)
	return nil
}

func resourcePingdomMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	window, err := maintenanceWindowForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Maintenance window create configuration: %#v", d.Get("description"))
	result, err := client.Maintenances.Create(window)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	return nil
}

func resourcePingdomMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	if _, err := client.Maintenances.Delete(id); err != nil {
		return fmt.Errorf("Error deleting maintenance window: %s", err)
	}
	return nil
}

func resourcePingdomMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	window, err := maintenanceWindowForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Maintenance window update configuration: %#v", d.Get("description"))

	if _, err = client.Maintenances.Update(id, window); err != nil {
		return fmt.Errorf("Error updating maintenance window: %s", err)
	}

	return nil
}
