package pingdom

import (
	"context"
	"strconv"
	"strings"

	"github.com/DrFaust92/go-pingdom/pingdom"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePingdomMaintenance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePingdomMaintenanceCreate,
		ReadContext:   resourcePingdomMaintenanceRead,
		UpdateContext: resourcePingdomMaintenanceUpdate,
		DeleteContext: resourcePingdomMaintenanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Required: true,
			},
			"effectiveto": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"recurrencetype": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "day", "week", "month"}, false),
			},
			"repeatevery": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tmsids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"uptimeids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func maintenanceForResource(d *schema.ResourceData) (*pingdom.MaintenanceWindow, error) {
	maintenance := pingdom.MaintenanceWindow{}

	// required
	if v, ok := d.GetOk("description"); ok {
		maintenance.Description = v.(string)
	}

	if v, ok, err := getTime("from", d); err != nil {
		return nil, err
	} else if ok {
		maintenance.From = v
	}

	if v, ok, err := getTime("to", d); err != nil {
		return nil, err
	} else if ok {
		maintenance.To = v
	}

	if v, ok, err := getTime("effectiveto", d); err != nil {
		return nil, err
	} else if ok {
		maintenance.EffectiveTo = v
	}

	if v, ok := d.GetOk("recurrencetype"); ok {
		maintenance.RecurrenceType = v.(string)
	}

	if v, ok := d.GetOk("repeatevery"); ok {
		maintenance.RepeatEvery = v.(int)
	}

	if v, ok := d.GetOk("tmsids"); ok {
		maintenance.TmsIDs = convertIntInterfaceSliceToString(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("uptimeids"); ok {
		maintenance.UptimeIDs = convertIntInterfaceSliceToString(v.(*schema.Set).List())
	}

	return &maintenance, nil
}

func updateResourceFromMaintenanceResponse(d *schema.ResourceData, m *pingdom.MaintenanceResponse) error {
	if err := d.Set("description", m.Description); err != nil {
		return err
	}

	if err := d.Set("from", timeFormat(m.From)); err != nil {
		return err
	}

	if err := d.Set("to", timeFormat(m.To)); err != nil {
		return err
	}

	if err := d.Set("effectiveto", timeFormat(m.EffectiveTo)); err != nil {
		return err
	}

	if err := d.Set("recurrencetype", m.RecurrenceType); err != nil {
		return err
	}

	if err := d.Set("repeatevery", m.RepeatEvery); err != nil {
		return err
	}

	tmsids := schema.NewSet(
		func(tmsId interface{}) int { return tmsId.(int) },
		[]interface{}{},
	)
	for _, tms := range m.Checks.Tms {
		tmsids.Add(tms)
	}
	if err := d.Set("tmsids", tmsids); err != nil {
		return err
	}

	uptimeids := schema.NewSet(
		func(uptimeId interface{}) int { return uptimeId.(int) },
		[]interface{}{},
	)
	for _, uptime := range m.Checks.Uptime {
		uptimeids.Add(uptime)
	}
	if err := d.Set("uptimeids", uptimeids); err != nil {
		return err
	}

	return nil
}

func resourcePingdomMaintenanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	maintenance, err := maintenanceForResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Maintenances.Create(maintenance)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(result.ID))

	return nil
}

func resourcePingdomMaintenanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving id for resource: %s", err)
	}
	maintenance, err := client.Maintenances.Read(id)
	if err != nil {
		return diag.Errorf("Error retrieving maintenance: %s", err)
	}

	if err := updateResourceFromMaintenanceResponse(d, maintenance); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePingdomMaintenanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving id for resource: %s", err)
	}
	maintenance, err := maintenanceForResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err = client.Maintenances.Update(id, maintenance); err != nil {
		return diag.Errorf("Error updating maintenance: %s", err)
	}

	return nil
}

func resourcePingdomMaintenanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error retrieving id for resource: %s", err)
	}
	if _, err := client.Maintenances.Delete(id); err != nil {
		return diag.Errorf("Error deleting maintenance: %s", err)
	}
	return nil
}

func convertIntInterfaceSliceToString(slice []interface{}) string {
	stringSlice := make([]string, len(slice))
	for i := range slice {
		stringSlice[i] = strconv.Itoa(slice[i].(int))
	}
	return strings.Join(stringSlice, ",")
}
