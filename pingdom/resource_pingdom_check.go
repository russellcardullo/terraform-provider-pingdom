package pingdom

import (
	"errors"
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

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"uselegacynotifications": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

type commonCheckParams struct {
	Name                     string
	Hostname                 string
	Resolution               int
	Paused                   bool
	SendToAndroid            bool
	SendToEmail              bool
	SendToIPhone             bool
	SendToSms                bool
	SendToTwitter            bool
	SendNotificationWhenDown int
	NotifyAgainEvery         int
	NotifyWhenBackup         bool
	UseLegacyNotifications   bool
}

func checkForResource(d *schema.ResourceData) (pingdom.Check, error) {
	checkParams := commonCheckParams{}

	// required
	if v, ok := d.GetOk("name"); ok {
		checkParams.Name = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		checkParams.Hostname = v.(string)
	}

	if v, ok := d.GetOk("resolution"); ok {
		checkParams.Resolution = v.(int)
	}

	// optional
	if v, ok := d.GetOk("sendtoemail"); ok {
		checkParams.SendToEmail = v.(bool)
	}

	if v, ok := d.GetOk("sendtosms"); ok {
		checkParams.SendToSms = v.(bool)
	}

	if v, ok := d.GetOk("sendtoiphone"); ok {
		checkParams.SendToIPhone = v.(bool)
	}

	if v, ok := d.GetOk("sendtoandroid"); ok {
		checkParams.SendToAndroid = v.(bool)
	}

	if v, ok := d.GetOk("sendnotificationwhendown"); ok {
		checkParams.SendNotificationWhenDown = v.(int)
	}

	if v, ok := d.GetOk("notifyagainevery"); ok {
		checkParams.NotifyAgainEvery = v.(int)
	}

	if v, ok := d.GetOk("notifywhenbackup"); ok {
		checkParams.NotifyWhenBackup = v.(bool)
	}

	if v, ok := d.GetOk("uselegacynotifications"); ok {
		checkParams.UseLegacyNotifications = v.(bool)
	}

	checkType := d.Get("type")
	switch checkType {
	case "http":
		return &pingdom.HttpCheck{
			Name:                     checkParams.Name,
			Hostname:                 checkParams.Hostname,
			Resolution:               checkParams.Resolution,
			Paused:                   checkParams.Paused,
			SendToAndroid:            checkParams.SendToAndroid,
			SendToEmail:              checkParams.SendToEmail,
			SendToIPhone:             checkParams.SendToIPhone,
			SendToSms:                checkParams.SendToSms,
			SendToTwitter:            checkParams.SendToTwitter,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			UseLegacyNotifications:   checkParams.UseLegacyNotifications,
		}, nil
	case "ping":
		return &pingdom.PingCheck{
			Name:                     checkParams.Name,
			Hostname:                 checkParams.Hostname,
			Resolution:               checkParams.Resolution,
			Paused:                   checkParams.Paused,
			SendToAndroid:            checkParams.SendToAndroid,
			SendToEmail:              checkParams.SendToEmail,
			SendToIPhone:             checkParams.SendToIPhone,
			SendToSms:                checkParams.SendToSms,
			SendToTwitter:            checkParams.SendToTwitter,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			UseLegacyNotifications:   checkParams.UseLegacyNotifications,
		}, nil
	default:
		errString := fmt.Sprintf("unknown type for check '%v'", checkType)
		return nil, errors.New(errString)
	}
}

func resourcePingdomCheckCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	check, err := checkForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Check create configuration: %#v, %#v", d.Get("name"), d.Get("hostname"))

	ck, err := client.Checks.Create(check)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ck.ID))

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

	check, err := checkForResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Check update configuration: %#v, %#v", d.Get("name"), d.Get("hostname"))

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
