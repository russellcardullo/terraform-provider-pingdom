package pingdom

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomContact() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomContactCreate,
		Read:   resourcePingdomContactRead,
		Update: resourcePingdomContactUpdate,
		Delete: resourcePingdomContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"severity_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

type commonContactParams struct {
	UserID        int
	Email         string
	Number        string
	PhoneProvider string
	CountryCode   string
	Severity      string
}

func contactForResource(d *schema.ResourceData) (pingdom.Contact, error) {
	contactParams := commonContactParams{}

	// required
	if v, ok := d.GetOk("user_id"); ok {
		contactParams.UserID = v.(int)
	}

	if v, ok := d.GetOk("email"); ok {
		contactParams.Email = v.(string)
	}

	if v, ok := d.GetOk("number"); ok {
		contactParams.Number = v.(string)
	}

	if v, ok := d.GetOk("country_code"); ok {
		contactParams.CountryCode = v.(string)
	}

	if v, ok := d.GetOk("phone_provider"); ok {
		contactParams.PhoneProvider = v.(string)
	}

	if v, ok := d.GetOk("severity_level"); ok {
		contactParams.Severity = strings.ToUpper(v.(string))
	}

	return pingdom.Contact{
		Severity:    contactParams.Severity,
		Email:       contactParams.Email,
		Number:      contactParams.Number,
		CountryCode: contactParams.CountryCode,
		Provider:    contactParams.PhoneProvider,
	}, nil
}

func resourcePingdomContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	contact, err := contactForResource(d)
	if err != nil {
		return err
	}

	userID, ok := d.Get("user_id").(int)
	if !ok {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[DEBUG] Contact create configuration: %#v", d.Get("Contactname"))
	result, err := client.Users.CreateContact(userID, contact)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", result.Id))
	return nil
}

func resourcePingdomContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	user := pingdom.UsersResponse{}
	users, err := client.Users.List()
	if err != nil {
		return fmt.Errorf("Error retrieving users: %s", err)
	}
	userFound := false
	for _, u := range users {
		for _, email := range u.Email {
			if email.Id == id {
				user = u
				userFound = true
				break
			}
		}
		for _, sms := range u.Sms {
			if sms.Id == id {
				user = u
				userFound = true
				break
			}
		}
		if userFound {
			break
		}
	}
	if !userFound {
		return fmt.Errorf("Error matching contact %d to a user", id)
	}

	for _, contact := range user.Email {
		if contact.Id == id {
			d.Set("email", contact.Address)
			d.Set("severity_level", contact.Severity)
			d.Set("user_id", user.Id)
			return nil
		}
	}
	for _, contact := range user.Sms {
		if contact.Id == id {
			d.Set("number", contact.Number)
			d.Set("country_code", contact.CountryCode)
			d.Set("phone_provider", contact.Provider)
			d.Set("severity_level", contact.Severity)
			d.Set("user_id", user.Id)
			return nil
		}
	}
	return nil
}

func resourcePingdomContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	userID, ok := d.Get("user_id").(int)
	if !ok {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	if _, err := client.Users.DeleteContact(userID, id); err != nil {
		return fmt.Errorf("Error deleting Contact: %s", err)
	}
	return nil
}

func resourcePingdomContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	contact, err := contactForResource(d)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	userID, ok := d.Get("user_id").(int)
	if !ok {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	if _, err := client.Users.UpdateContact(userID, id, contact); err != nil {
		return fmt.Errorf("Error updating Contact: %s", err)
	}
	return nil
}
