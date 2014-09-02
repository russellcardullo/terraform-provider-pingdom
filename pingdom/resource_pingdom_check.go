package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/pingdom"
)

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomCheckCreate,
		Read:   resourcePingdomCheckRead,
		Delete: resourcePingdomCheckDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePingdomCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	name := d.Get("name").(string)
	host := d.Get("host").(string)
	check := pingdom.HttpCheck{Name: name, Host: host}

	log.Printf("[DEBUG] Check create configuration: %#v, %#v", name, host)

	ck, err := client.CreateCheck(check)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ck.ID))
	d.Set("hostname", ck.Hostname)
	d.Set("name", ck.Name)

	return nil
}

func resourcePingdomCheckDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Deleting Check: %v", id)

	_, err = client.DeleteCheck(id)
	if err != nil {
		return fmt.Errorf("Error deleting check: %s", err)
	}

	return nil
}

func resourcePingdomCheckRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	ck, err := client.ReadCheck(id)
	if err != nil {
		return fmt.Errorf("Error retrieving check: %s", err)
	}

	d.Set("hostname", ck.Hostname)
	d.Set("name", ck.Name)

	return nil
}
