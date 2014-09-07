package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomCheckCreate,
		Read:   resourcePingdomCheckRead,
		Update: resourcePingdomCheckUpdate,
		Delete: resourcePingdomCheckDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"resolution": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"sendtoemail": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"sendtosms": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"sendtotwitter": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"sendtoiphone": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"sendtoandroid": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"sendnotificationwhendown": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			"notifyagainevery": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			"notifywhenbackup": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func checkForResource(d *schema.ResourceData) *pingdom.Check {
	check := &pingdom.Check{}
	// required
	if v, ok := d.GetOk("name"); ok {
		check.Name = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		check.Hostname = v.(string)
	}

	if v, ok := d.GetOk("resolution"); ok {
		check.Resolution = v.(int)
	}

	// optional
	if v, ok := d.GetOk("sendtoemail"); ok {
		check.SendToEmail = v.(bool)
	}

	if v, ok := d.GetOk("sendtosms"); ok {
		check.SendToSms = v.(bool)
	}

	if v, ok := d.GetOk("sendtoiphone"); ok {
		check.SendToIPhone = v.(bool)
	}

	if v, ok := d.GetOk("sendtoandroid"); ok {
		check.SendToAndroid = v.(bool)
	}

	if v, ok := d.GetOk("sendnotificationwhendown"); ok {
		check.SendNotificationWhenDown = v.(int)
	}

	if v, ok := d.GetOk("notifyagainevery"); ok {
		check.NotifyAgainEvery = v.(int)
	}

	if v, ok := d.GetOk("notifywhenbackup"); ok {
		check.NotifyWhenBackup = v.(bool)
	}
	return check
}

func resourcePingdomCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	check := checkForResource(d)
	log.Printf("[DEBUG] Check create configuration: %#v, %#v", check.Name, check.Hostname)

	ck, err := client.Checks.Create(check)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ck.ID))
	d.Set("hostname", ck.Hostname)
	d.Set("name", ck.Name)

	return nil
}

func resourcePingdomCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	ck, err := client.Checks.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving check: %s", err)
	}

	d.Set("hostname", ck.Hostname)
	d.Set("name", ck.Name)
	d.Set("resolution", ck.Resolution)

	return nil
}

func resourcePingdomCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	check := checkForResource(d)
	log.Printf("[DEBUG] Check update configuration: %#v, %#v", check.Name, check.Hostname)

	_, err = client.Checks.Update(id, check)
	if err != nil {
		return fmt.Errorf("Error updating check: %s", err)
	}

	return nil
}

func resourcePingdomCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Deleting Check: %v", id)

	_, err = client.Checks.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting check: %s", err)
	}

	return nil
}
