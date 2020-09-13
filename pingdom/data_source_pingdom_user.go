package pingdom

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func dataSourcePingdomUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePingdomUserRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePingdomUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	username := d.Get("username").(string)
	users, err := client.Users.List()
	log.Printf("==== users : %v", users)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %s", err)
	}
	var found pingdom.UsersResponse
	for _, user := range users {
		if user.Username == username {
			log.Printf("User: %v", user)
			found = user
		}
	}
	err = d.Set("username", found.Username)
	if err != nil {
		return fmt.Errorf("Error setting username: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", found.Id))
	return nil
}
