package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomUser() *schema.Resource {
	resource := resourceUser{}
	return &schema.Resource{
		Create: resource.create,
		Read:   resource.read,
		Delete: resource.read,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

type resourceUser struct{}

func (ru *resourceUser) checkForResource(d *schema.ResourceData) (pingdom.UserApi, error) {
	u := &pingdom.User{}
	if v, ok := d.GetOk("username"); ok {
		u.Username = v.(string)
	}
	return u, nil
}

func (ru *resourceUser) create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	u, err := ru.checkForResource(d)
	if err != nil {
		return err
	}
	cu, err := client.Users.Create(u)
	if err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	d.SetId(strconv.Itoa(cu.Id))
	return nil
}

func (ru *resourceUser) read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	ul, err := client.Users.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of users: %s", err)
	}

	exists := false
	for _, uid := range ul {
		if uid.Id == id {
			exists = true
			break
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}

	u, err := client.Users.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %s", err)
	}

	d.Set("email", u.Email)
	d.Set("paused", u.Paused)
	d.Set("username", u.Username)
	return nil
}

func (ru *resourceUser) delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Deleting user: %v", id)
	_, err = client.Users.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}
	d.SetId("")
	return nil
}

func (ru *resourceUser) importer(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	return nil, nil
}
