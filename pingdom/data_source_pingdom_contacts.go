package pingdom

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePingdomContacts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePingdomContactsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePingdomContactsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Clients).Pingdom
	contacts, err := client.Contacts.List()
	if err != nil {
		return diag.Errorf("Error retrieving contacts: %s", err)
	}

	var ids = make([]int, 0, len(contacts))
	var names = make([]string, 0, len(contacts))
	var types = make([]string, 0, len(contacts))
	for _, contact := range contacts {
		ids = append(ids, contact.ID)
		names = append(names, contact.Name)
		types = append(types, contact.Type)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("names", names); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("types", types); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
