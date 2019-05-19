package pingdom

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomContact() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomContactCreate,
		Read:   resourcePingdomContactRead,
		Update: resourcePingdomContactUpdate,
		Delete: resourcePingdomContactDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: resourcePingdomContactImporter,
		// },
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
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
	UserID        string
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
		contactParams.UserID = v.(string)
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

	userID, err := strconv.Atoi(d.Get("user_id").(string))
	if err != nil {
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

	userID, err := strconv.Atoi(d.Get("user_id").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	user, err := client.Users.Read(userID)
	if err != nil {
		return fmt.Errorf("Error retrieving Contact: %s", err)
	}

	for _, contact := range user.Email {
		if contact.Id == id {
			d.Set("email", contact.Address)
			d.Set("severity", contact.Severity)
			return nil
		}
	}
	for _, contact := range user.Sms {
		if contact.Id == id {
			d.Set("number", contact.Number)
			d.Set("country_code", contact.CountryCode)
			d.Set("phone_provider", contact.Provider)
			d.Set("severity_level", contact.Severity)
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

	userID, err := strconv.Atoi(d.Get("user_id").(string))
	if err != nil {
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

	userID, err := strconv.Atoi(d.Get("user_id").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	if _, err := client.Users.UpdateContact(userID, id, contact); err != nil {
		return fmt.Errorf("Error updating Contact: %s", err)
	}
	return nil
}
