package pingdom

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomCheckCreate,
		Read:   resourcePingdomCheckRead,
		Update: resourcePingdomCheckUpdate,
		Delete: resourcePingdomCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"paused": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},

			"responsetime_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
				StateFunc: func(val interface{}) string {
					return sortString(val.(string), ",")
				},
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
	ResponseTimeThreshold    int
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

func sortString(input string, seperator string) string {
	list := strings.Split(input, seperator)
	sort.Strings(list)
	return strings.Join(list, seperator)
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

	if v, ok := d.GetOk("paused"); ok {
		checkParams.Paused = v.(bool)
	}

	if v, ok := d.GetOk("resolution"); ok {
		checkParams.Resolution = v.(int)
	}

	if v, ok := d.GetOk("responsetime_threshold"); ok {
		checkParams.ResponseTimeThreshold = v.(int)
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
		// Sort alphabetically before contionuing
		checkParams.Tags = sortString(v.(string), ",")
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
			Name:                     checkParams.Name,
			Hostname:                 checkParams.Hostname,
			Resolution:               checkParams.Resolution,
			Paused:                   checkParams.Paused,
			ResponseTimeThreshold:    checkParams.ResponseTimeThreshold,
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
			Name:                     checkParams.Name,
			Hostname:                 checkParams.Hostname,
			Resolution:               checkParams.Resolution,
			Paused:                   checkParams.Paused,
			ResponseTimeThreshold:    checkParams.ResponseTimeThreshold,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			IntegrationIds:           checkParams.IntegrationIds,
			Tags:                     checkParams.Tags,
			ProbeFilters:             checkParams.ProbeFilters,
			UserIds:                  checkParams.UserIds,
			TeamIds:                  checkParams.TeamIds,
		}, nil
	case "tcp":
		return &pingdom.TCPCheck{
			Name:                     checkParams.Name,
			Hostname:                 checkParams.Hostname,
			Resolution:               checkParams.Resolution,
			Paused:                   checkParams.Paused,
			SendNotificationWhenDown: checkParams.SendNotificationWhenDown,
			NotifyAgainEvery:         checkParams.NotifyAgainEvery,
			NotifyWhenBackup:         checkParams.NotifyWhenBackup,
			IntegrationIds:           checkParams.IntegrationIds,
			Tags:                     checkParams.Tags,
			ProbeFilters:             checkParams.ProbeFilters,
			UserIds:                  checkParams.UserIds,
			TeamIds:                  checkParams.TeamIds,
			Port:                     checkParams.Port,
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

	return resourcePingdomCheckRead(d, meta)
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

	if err := d.Set("host", ck.Hostname); err != nil {
		return err
	}

	if err := d.Set("name", ck.Name); err != nil {
		return err
	}

	if err := d.Set("resolution", ck.Resolution); err != nil {
		return err
	}

	if err := d.Set("responsetime_threshold", ck.ResponseTimeThreshold); err != nil {
		return err
	}

	if err := d.Set("sendnotificationwhendown", ck.SendNotificationWhenDown); err != nil {
		return err
	}

	if err := d.Set("notifyagainevery", ck.NotifyAgainEvery); err != nil {
		return err
	}

	if err := d.Set("notifywhenbackup", ck.NotifyWhenBackup); err != nil {
		return err
	}

	tags := []string{}
	for _, tag := range ck.Tags {
		tags = append(tags, tag.Name)
	}

	// We need to sort the strings here as the pingdom API returns them sorted by
	//number of occurances across all checks
	sort.Strings(tags)
	if err := d.Set("tags", strings.Join(tags, ",")); err != nil {
		return err
	}

	if ck.Status == "paused" {
		if err := d.Set("paused", true); err != nil {
			return err
		}
	}

	integids := schema.NewSet(
		func(integrationId interface{}) int { return integrationId.(int) },
		[]interface{}{},
	)
	for _, integrationId := range ck.IntegrationIds {
		integids.Add(integrationId)
	}
	if err := d.Set("integrationids", integids); err != nil {
		return err
	}

	userids := schema.NewSet(
		func(userId interface{}) int { return userId.(int) },
		[]interface{}{},
	)
	for _, userId := range ck.UserIds {
		userids.Add(userId)
	}
	if err := d.Set("userids", userids); err != nil {
		return err
	}

	teamids := schema.NewSet(
		func(userId interface{}) int { return userId.(int) },
		[]interface{}{},
	)
	for _, userId := range ck.TeamIds {
		teamids.Add(userId)
	}
	if err := d.Set("teamids", teamids); err != nil {
		return err
	}

	if probefilters := ck.ProbeFilters; len(probefilters) > 0 {
		// normalise: "region: NA" -> "region:NA"
		if err := d.Set("probefilters", strings.Replace(probefilters[0], ": ", ":", 1)); err != nil {
			return err
		}
	}

	if ck.Type.HTTP != nil {
		if err := d.Set("type", "http"); err != nil {
			return err
		}
		if err := d.Set("responsetime_threshold", ck.ResponseTimeThreshold); err != nil {
			return err
		}
		if err := d.Set("url", ck.Type.HTTP.Url); err != nil {
			return err
		}
		if err := d.Set("encryption", ck.Type.HTTP.Encryption); err != nil {
			return err
		}
		if err := d.Set("port", ck.Type.HTTP.Port); err != nil {
			return err
		}
		if err := d.Set("username", ck.Type.HTTP.Username); err != nil {
			return err
		}
		if err := d.Set("password", ck.Type.HTTP.Password); err != nil {
			return err
		}
		if err := d.Set("shouldcontain", ck.Type.HTTP.ShouldContain); err != nil {
			return err
		}
		if err := d.Set("shouldnotcontain", ck.Type.HTTP.ShouldNotContain); err != nil {
			return err
		}
		if err := d.Set("postdata", ck.Type.HTTP.PostData); err != nil {
			return err
		}

		if v, ok := ck.Type.HTTP.RequestHeaders["User-Agent"]; ok {
			if strings.HasPrefix(v, "Pingdom.com_bot_version_") {
				delete(ck.Type.HTTP.RequestHeaders, "User-Agent")
			}
		}
		if err := d.Set("requestheaders", ck.Type.HTTP.RequestHeaders); err != nil {
			return err
		}
	} else if ck.Type.TCP != nil {
		if err := d.Set("type", "tcp"); err != nil {
			return err
		}
		if err := d.Set("port", ck.Type.TCP.Port); err != nil {
			return err
		}
		if err := d.Set("stringtosend", ck.Type.TCP.StringToSend); err != nil {
			return err
		}
		if err := d.Set("stringtoexpect", ck.Type.TCP.StringToExpect); err != nil {
			return err
		}
	} else {
		if err := d.Set("type", "ping"); err != nil {
			return err
		}
	}

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

	return resourcePingdomCheckRead(d, meta)
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
