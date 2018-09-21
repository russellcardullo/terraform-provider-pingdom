package pingdom

import (
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
		Importer: &schema.ResourceImporter{
			State: resourcePingdomCheckImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"publicreport": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"resolution": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"sendnotificationwhendown": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"notifyagainevery": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},

			"notifywhenbackup": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"integrationids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "/",
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldcontain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"shouldnotcontain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"postdata": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"requestheaders": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"probefilters": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"userids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"teamids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			"stringtosend": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"stringtoexpect": {
				Type:     schema.TypeString,
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
	SendNotificationWhenDown int
	NotifyAgainEvery         int
	NotifyWhenBackup         bool
	IntegrationIds           []int
	UserIds                  []int
	TeamIds                  []int
	Url                      string
	Encryption               bool
	Port                     int
	Username                 string
	Password                 string
	ShouldContain            string
	ShouldNotContain         string
	PostData                 string
	RequestHeaders           map[string]string
	Tags                     string
	ProbeFilters             string
	StringToSend             string
	StringToExpect           string
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

	if v, ok := d.GetOk("sendnotificationwhendown"); ok {
		checkParams.SendNotificationWhenDown = v.(int)
	}

	if v, ok := d.GetOk("notifyagainevery"); ok {
		checkParams.NotifyAgainEvery = v.(int)
	}

	if v, ok := d.GetOk("notifywhenbackup"); ok {
		checkParams.NotifyWhenBackup = v.(bool)
	}

	if v, ok := d.GetOk("integrationids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		checkParams.IntegrationIds = intSlice
	}

	if v, ok := d.GetOk("userids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		checkParams.UserIds = intSlice
	}

	if v, ok := d.GetOk("teamids"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var intSlice []int
		for i := range interfaceSlice {
			intSlice = append(intSlice, interfaceSlice[i].(int))
		}
		checkParams.TeamIds = intSlice
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
	if v, ok := d.GetOk("tags"); ok {
		checkParams.Tags = v.(string)
	}

	if v, ok := d.GetOk("probefilters"); ok {
		checkParams.ProbeFilters = v.(string)
	}

	if v, ok := d.GetOk("stringtosend"); ok {
		checkParams.StringToSend = v.(string)
	}

	if v, ok := d.GetOk("stringtoexpect"); ok {
		checkParams.StringToExpect = v.(string)
	}

	checkType := d.Get("type")
	switch checkType {
	case "http":
		return &pingdom.HttpCheck{
			Name:       checkParams.Name,
			Hostname:   checkParams.Hostname,
			Resolution: checkParams.Resolution,
			Paused:     checkParams.Paused,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			IntegrationIds:           checkParams.IntegrationIds,
			Encryption:               checkParams.Encryption,
			Url:                      checkParams.Url,
			Port:                     checkParams.Port,
			Username:                 checkParams.Username,
			Password:                 checkParams.Password,
			ShouldContain:            checkParams.ShouldContain,
			ShouldNotContain:         checkParams.ShouldNotContain,
			PostData:                 checkParams.PostData,
			RequestHeaders:           checkParams.RequestHeaders,
			Tags:                     checkParams.Tags,
			ProbeFilters:             checkParams.ProbeFilters,
			UserIds:                  checkParams.UserIds,
			TeamIds:                  checkParams.TeamIds,
		}, nil
	case "ping":
		return &pingdom.PingCheck{
			Name:       checkParams.Name,
			Hostname:   checkParams.Hostname,
			Resolution: checkParams.Resolution,
			Paused:     checkParams.Paused,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			IntegrationIds:           checkParams.IntegrationIds,
			UserIds:                  checkParams.UserIds,
			TeamIds:                  checkParams.TeamIds,
		}, nil
	case "tcp":
		return &pingdom.TCPCheck{
			Name:       checkParams.Name,
			Hostname:   checkParams.Hostname,
			Resolution: checkParams.Resolution,
			Paused:     checkParams.Paused,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			IntegrationIds:           checkParams.IntegrationIds,
			UserIds:                  checkParams.UserIds,
			TeamIds:                  checkParams.TeamIds,
			StringToSend:             checkParams.StringToSend,
			StringToExpect:           checkParams.StringToExpect,
		}, nil
	default:
		return nil, fmt.Errorf("unknown type for check '%v'", checkType)
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

	if v, ok := d.GetOk("publicreport"); ok && v.(bool) {
		_, _ := client.PublicReport.PublishCheck(ck.ID)
	}

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
	rl, err := client.PublicReport.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of public report checks: %s", err)
	}
	inPublicReport := false
	for _, ckid := range rl {
		if ckid.ID == id {
			inPublicReport = true
			break
		}
	}

	d.Set("host", ck.Hostname)
	d.Set("name", ck.Name)
	d.Set("resolution", ck.Resolution)
	d.Set("sendnotificationwhendown", ck.SendNotificationWhenDown)
	d.Set("notifyagainevery", ck.NotifyAgainEvery)
	d.Set("notifywhenbackup", ck.NotifyWhenBackup)
	d.Set("publicreport", inPublicReport)

	integids := schema.NewSet(
		func(integrationId interface{}) int { return integrationId.(int) },
		[]interface{}{},
	)
	for _, integrationId := range ck.IntegrationIds {
		integids.Add(integrationId)
	}
	d.Set("integrationids", integids)

	userids := schema.NewSet(
		func(userId interface{}) int { return userId.(int) },
		[]interface{}{},
	)
	for _, userId := range ck.UserIds {
		userids.Add(userId)
	}
	d.Set("userids", userids)

	teamids := schema.NewSet(
		func(userId interface{}) int { return userId.(int) },
		[]interface{}{},
	)
	for _, userId := range ck.TeamIds {
		teamids.Add(userId)
	}
	d.Set("teamids", teamids)

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

	if ck.Type.TCP == nil {
		ck.Type.TCP = &pingdom.CheckResponseTCPDetails{}
	}
	d.Set("port", ck.Type.TCP.Port)
	d.Set("stringtosend", ck.Type.TCP.StringToSend)
	d.Set("stringtoexpect", ck.Type.TCP.StringToExpect)

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

	if v, ok := d.GetOk("publicreport"); ok && v.(bool) {
		_, _ := client.PublicReport.PublishCheck(id)
	} else {
		_, _ := client.PublicReport.WithdrawlCheck(id)
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

func resourcePingdomCheckImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Importing key using ADDR ID %d", id)

	return []*schema.ResourceData{d}, nil
}
