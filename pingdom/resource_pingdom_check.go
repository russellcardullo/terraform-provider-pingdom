package pingdom

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

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
				Computed: true,
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
				Computed: true,
			},

			"uselegacynotifications": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"encryption": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldcontain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldnotcontain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"postdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"requestheaders": &schema.Schema{
				Type:     schema.TypeMap,
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
	Url                      string
	Encryption               bool
	Port                     int
	Username                 string
	Password                 string
	ShouldContain            string
	ShouldNotContain         string
	PostData                 string
	RequestHeaders           map[string]string
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

	if v, ok := d.GetOk("url"); ok {
		checkParams.Url = v.(string)
	}

	if v, ok := d.GetOk("encryption"); ok {
		checkParams.Encryption = v.(bool)
	}

	if v, ok := d.GetOk("port"); ok {
		checkParams.Port = v.(int)
	}

	if v, ok := d.GetOk("username"); ok {
		checkParams.Username = v.(string)
	}

	if v, ok := d.GetOk("password"); ok {
		checkParams.Password = v.(string)
	}

	if v, ok := d.GetOk("shouldcontain"); ok {
		checkParams.ShouldContain = v.(string)
	}

	if v, ok := d.GetOk("shouldnotcontain"); ok {
		checkParams.ShouldNotContain = v.(string)
	}

	if v, ok := d.GetOk("postdata"); ok {
		checkParams.PostData = v.(string)
	}

	if m, ok := d.GetOk("requestheaders"); ok {
		checkParams.RequestHeaders = make(map[string]string)
		for k, v := range m.(map[string]interface{}) {
			checkParams.RequestHeaders[k] = v.(string)
		}
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
			Encryption:               checkParams.Encryption,
			Url:                      checkParams.Url,
			Port:                     checkParams.Port,
			Username:                 checkParams.Username,
			Password:                 checkParams.Password,
			ShouldContain:            checkParams.ShouldContain,
			ShouldNotContain:         checkParams.ShouldNotContain,
			PostData:                 checkParams.PostData,
			RequestHeaders:           checkParams.RequestHeaders,
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
	cl, err := client.Checks.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of checks: %s", err)
	}
	exists := false
	for _, ckid := range cl {
		if ckid.ID == id {
			exists = true
			break
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}
	ck, err := client.Checks.Read(id)
	if err != nil {
		return fmt.Errorf("Error retrieving check: %s", err)
	}

	d.Set("hostname", ck.Hostname)
	d.Set("name", ck.Name)
	d.Set("resolution", ck.Resolution)
	d.Set("sendtoandroid", ck.SendToAndroid)
	d.Set("sendtoemail", ck.SendToEmail)
	d.Set("sendtoiphone", ck.SendToIPhone)
	d.Set("sendtosms", ck.SendToSms)
	d.Set("sendtotwitter", ck.SendToTwitter)
	d.Set("sendnotificationwhendown", ck.SendNotificationWhenDown)
	d.Set("notifyagainevery", ck.NotifyAgainEvery)
	d.Set("notifywhenbackup", ck.NotifyWhenBackup)
	d.Set("hostname", ck.Hostname)

	if ck.Type.HTTP == nil {
		ck.Type.HTTP = &pingdom.CheckResponseHTTPDetails{}
	}
	d.Set("url", ck.Type.HTTP.Url)
	d.Set("encryption", ck.Type.HTTP.Encryption)
	d.Set("port", ck.Type.HTTP.Port)
	d.Set("username", ck.Type.HTTP.Username)
	d.Set("password", ck.Type.HTTP.Password)
	d.Set("shouldcontain", ck.Type.HTTP.ShouldContain)
	d.Set("shouldnotcontain", ck.Type.HTTP.ShouldNotContain)
	d.Set("postdata", ck.Type.HTTP.PostData)

	if v, ok := ck.Type.HTTP.RequestHeaders["User-Agent"]; ok {
		if strings.HasPrefix(v, "Pingdom.com_bot_version_") {
			delete(ck.Type.HTTP.RequestHeaders, "User-Agent")
		}
	}
	d.Set("requestheaders", ck.Type.HTTP.RequestHeaders)

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
