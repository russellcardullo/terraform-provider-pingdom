package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomUserCreate,
		Read:   resourcePingdomUserRead,
		Update: resourcePingdomUserUpdate,
		Delete: resourcePingdomUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type commonUserParams struct {
	Username string
}

func userForResource(d *schema.ResourceData) (*pingdom.User, error) {
	userParams := commonUserParams{}

	// required
	if v, ok := d.GetOk("username"); ok {
		userParams.Username = v.(string)
	}

	return &pingdom.User{
		Username: userParams.Username,
	}, nil
}

func resourcePingdomUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	user, err := userForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] User create configuration: %#v", d.Get("username"))
	result, err := client.Users.Create(user)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", result.Id))
	return nil
}

func resourcePingdomUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	user, err := client.Users.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %s", err)
	}

	if err := d.Set("username", user.Username); err != nil {
		return err
	}
	return nil
}

func resourcePingdomUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	user, err := userForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] User update configuration: %#v", d.Get("username"))

	if _, err = client.Users.Update(id, user); err != nil {
		return fmt.Errorf("Error updating user: %s", err)
	}

	return nil
}

func resourcePingdomUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	if _, err := client.Users.Delete(id); err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}
	return nil
}
