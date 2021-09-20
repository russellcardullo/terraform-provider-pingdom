package pingdom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sms_notification": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:     schema.TypeString,
							Required: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "1",
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "nexmo",
						},
					},
				},
			},
			"email_notification": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func getNotificationMethods(d *schema.ResourceData) (pingdom.NotificationTargets, error) {
	base := pingdom.NotificationTargets{}

	// You must have both a low and a high severity notification for a user to be valid
	hasLowSeverity := false
	hasHighSeverity := false

	for _, raw := range d.Get("sms_notification").(*schema.Set).List() {
		input := raw.(map[string]interface{})
		sms := pingdom.SMSNotification{
			CountryCode: input["country_code"].(string),
			Number:      input["number"].(string),
			Provider:    input["provider"].(string),
			Severity:    input["severity"].(string),
		}
		if sms.Severity == "HIGH" {
			hasHighSeverity = true
		}
		if sms.Severity == "LOW" {
			hasLowSeverity = true
		}
		switch sms.Provider {
		case "nexmo", "bulksms", "esendex", "cellsynt":
			base.SMS = append(base.SMS, sms)
			continue
		}

		return base, fmt.Errorf("SMS provider must be one of: nexmo, bulksms, esendex, or cellsynt")
	}

	for _, raw := range d.Get("email_notification").(*schema.Set).List() {
		input := raw.(map[string]interface{})
		email := pingdom.EmailNotification{
			Address:  input["address"].(string),
			Severity: input["severity"].(string),
		}
		if email.Severity == "HIGH" {
			hasHighSeverity = true
		}
		if email.Severity == "LOW" {
			hasLowSeverity = true
		}
		base.Email = append(base.Email, email)
	}

	if !hasHighSeverity || !hasLowSeverity {
		return base, fmt.Errorf("You must provide both a high and low severity notification method")
	}

	return base, nil
}

func contactForResource(d *schema.ResourceData) (*pingdom.Contact, error) {
	contact := pingdom.Contact{}

	// required
	if v, ok := d.GetOk("name"); ok {
		contact.Name = v.(string)
	}

	notifications, err := getNotificationMethods(d)
	if err != nil {
		return nil, err
	}
	contact.NotificationTargets = notifications

	if v, ok := d.GetOk("paused"); ok {
		contact.Paused = v.(bool)
	}

	return &contact, nil
}

func updateResourceFromContactResponse(d *schema.ResourceData, c *pingdom.Contact) error {
	smsTargets := []map[string]string{}
	for _, raw := range c.NotificationTargets.SMS {
		sms := map[string]string{
			"country_code": raw.CountryCode,
			"number":       raw.Number,
			"severity":     raw.Severity,
			"provider":     raw.Provider,
		}
		smsTargets = append(smsTargets, sms)
	}
	if err := d.Set("sms_notification", smsTargets); err != nil {
		return err
	}

	emailTargets := []map[string]string{}
	for _, raw := range c.NotificationTargets.Email {
		email := map[string]string{
			"address":  raw.Address,
			"severity": raw.Severity,
		}
		emailTargets = append(emailTargets, email)
	}
	if err := d.Set("email_notification", emailTargets); err != nil {
		return err
	}

	if err := d.Set("paused", c.Paused); err != nil {
		return err
	}

	return nil
}

func resourcePingdomContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	contact, err := contactForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Contact create configuration: %#v", d.Get("name"))
	result, err := client.Contacts.Create(contact)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	return nil
}

func resourcePingdomContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	contact, err := client.Contacts.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving contact: %s", err)
	}

	if err := d.Set("name", contact.Name); err != nil {
		return err
	}

	if err := updateResourceFromContactResponse(d, contact); err != nil {
		return err
	}
	return nil
}

func resourcePingdomContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	contact, err := contactForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Contact update configuration: %#v", d.Get("name"))

	if _, err = client.Contacts.Update(id, contact); err != nil {
		return fmt.Errorf("Error updating contact: %s", err)
	}

	return nil
}

func resourcePingdomContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	if _, err := client.Contacts.Delete(id); err != nil {
		return fmt.Errorf("Error deleting contact: %s", err)
	}
	return nil
}
