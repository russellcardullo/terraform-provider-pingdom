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
		Create: resourcePingdomMaintenanceWindowCreate,
		Read:   resourcePingdomMaintenanceWindowRead,
		Update: resourcePingdomMaintenanceWindowUpdate,
		Delete: resourcePingdomMaintenanceWindowDelete,
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
				Optional: true,
				Default:  "none",
			},
			"repeat_every": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"effective_to": {
				Type:     schema.TypeInt,
				Optional: true,
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

type commonMaintenanceWindowParams struct {
	Description   string
	From          int64
	To            int64
	RecurenceType string
	RepeatEvery   int
	EffectiveTo   int
	UptimeIds     []int
	TmsIds        []int
}

func maintenanceWindowForResource(d *schema.ResourceData) (*pingdom.MaintenanceWindow, error) {
	windowParams := commonMaintenanceWindowParams{}

	if v, ok := d.GetOk("description"); ok {
		windowParams.Description = v.(string)
	}
	if v, ok := d.GetOk("from"); ok {
		windowParams.From = v.(int64)
	}
	if v, ok := d.GetOk("to"); ok {
		windowParams.To = v.(int64)
	}
	if v, ok := d.GetOk("recurrence_type"); ok {
		windowParams.RecurenceType = v.(string)
	}
	if v, ok := d.GetOk("repeat_every"); ok {
		windowParams.RepeatEvery = v.(int)
	}
	if v, ok := d.GetOk("effective_to"); ok {
		windowParams.EffectiveTo = v.(int)
	}
	if v, ok := d.GetOk("uptimeids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		windowParams.UptimeIds = intSlice
	}
	if v, ok := d.GetOk("tmsids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		windowParams.TmsIds = intSlice
	}

	uptimeIds := ""
	for _, id := range windowParams.UptimeIds {
		uptimeIds += strconv.Itoa(id)
	}

	tmsIds := ""
	for _, id := range windowParams.TmsIds {
		tmsIds += strconv.Itoa(id)
	}

	return &pingdom.MaintenanceWindow{
		Description:    windowParams.Description,
		From:           windowParams.From,
		To:             windowParams.To,
		RecurrenceType: windowParams.RecurenceType,
		RepeatEvery:    windowParams.RepeatEvery,
		EffectiveTo:    windowParams.EffectiveTo,
		UptimeIDs:      uptimeIds,
		TmsIDs:         tmsIds,
	}, nil
}

func resourcePingdomMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	log.Printf("[DEBUG] Read Maintenance window with ID: %v", d.Id())
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	// log.Printf("========================== Maintenances")
	// maintenances, err := client.Maintenances.List()
	// log.Println("Maintenances: ", maintenances) // [{ID Description} ...]
	// for _, m := range maintenances {
	// 	log.Println("Maintenance: ", m)
	// }

	window, err := client.Maintenances.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving maintenance window: %s", err)
	}
	log.Printf("[DEBUG] Maintenance: %+v\n", window)

	d.Set("descrition", window.Description)
	d.Set("from", window.From)
	d.Set("to", window.To)
	d.Set("recurence_type", window.RecurrenceType)
	d.Set("repeat_every", window.RepeatEvery)
	d.Set("effective_to", window.EffectiveTo)

	uptimeIds := schema.NewSet(
		func(id interface{}) int { return id.(int) },
		[]interface{}{},
	)
	for _, id := range window.Checks.Uptime {
		log.Println("Uptime: ", id)
		uptimeIds.Add(id)
	}
	d.Set("uptimeids", uptimeIds)

	tmsIds := schema.NewSet(
		func(id interface{}) int { return id.(int) },
		[]interface{}{},
	)
	for _, id := range window.Checks.Tms {
		log.Println("Tm: ", id)
		tmsIds.Add(id)
	}
	d.Set("tmsids", tmsIds)

	return nil
}

func resourcePingdomMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	log.Printf("[DEBUG] Create Maintenance window with ID: %v", d.Id())
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

	log.Printf("[DEBUG] Delete Maintenance window with ID: %v", d.Id())
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

	log.Printf("[DEBUG] Update Maintenance window with ID: %v", d.Id())
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
