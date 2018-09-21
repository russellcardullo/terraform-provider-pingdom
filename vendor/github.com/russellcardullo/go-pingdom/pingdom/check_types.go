package pingdom

import (
	"fmt"
	"sort"
	"strconv"
)

// HttpCheck represents a Pingdom http check.
type HttpCheck struct {
	Name                     string            `json:"name"`
	Hostname                 string            `json:"hostname,omitempty"`
	Resolution               int               `json:"resolution,omitempty"`
	Paused                   bool              `json:"paused,omitempty"`
	SendNotificationWhenDown int               `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int               `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool              `json:"notifywhenbackup,omitempty"`
	Url                      string            `json:"url,omitempty"`
	Encryption               bool              `json:"encryption,omitempty"`
	Port                     int               `json:"port,omitempty"`
	Username                 string            `json:"username,omitempty"`
	Password                 string            `json:"password,omitempty"`
	ShouldContain            string            `json:"shouldcontain,omitempty"`
	ShouldNotContain         string            `json:"shouldnotcontain,omitempty"`
	PostData                 string            `json:"postdata,omitempty"`
	RequestHeaders           map[string]string `json:"requestheaders,omitempty"`
	IntegrationIds           []int             `json:"integrationids,omitempty"`
	Tags                     string            `json:"tags,omitempty"`
	ProbeFilters             string            `json:"probe_filters,omitempty"`
	UserIds                  []int             `json:"userids,omitempty"`
	TeamIds                  []int             `json:"teamids,omitempty"`
}

// PingCheck represents a Pingdom ping check
type PingCheck struct {
	Name                     string `json:"name"`
	Hostname                 string `json:"hostname,omitempty"`
	Resolution               int    `json:"resolution,omitempty"`
	Paused                   bool   `json:"paused,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
	IntegrationIds           []int  `json:"integrationids,omitempty"`
	ProbeFilters             string `json:"probe_filters,omitempty"`
	UserIds                  []int  `json:"userids,omitempty"`
	TeamIds                  []int  `json:"teamids,omitempty"`
}

// TCPCheck represents a Pingdom TCP check
type TCPCheck struct {
	Name                     string `json:"name"`
	Hostname                 string `json:"hostname,omitempty"`
	Resolution               int    `json:"resolution,omitempty"`
	Paused                   bool   `json:"paused,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
	IntegrationIds           []int  `json:"integrationids,omitempty"`
	ProbeFilters             string `json:"probe_filters,omitempty"`
	UserIds                  []int  `json:"userids,omitempty"`
	TeamIds                  []int  `json:"teamids,omitempty"`
	Port                     int    `json:"port"`
	StringToSend             string `json:"stringtosend,omitempty"`
	StringToExpect           string `json:"stringtoexpect,omitempty"`
}

// Params returns a map of parameters for an HttpCheck that can be sent along
// with an HTTP PUT request
func (ck *HttpCheck) PutParams() map[string]string {
	m := map[string]string{
		"name":             ck.Name,
		"host":             ck.Hostname,
		"resolution":       strconv.Itoa(ck.Resolution),
		"paused":           strconv.FormatBool(ck.Paused),
		"notifyagainevery": strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup": strconv.FormatBool(ck.NotifyWhenBackup),
		"url":              ck.Url,
		"encryption":       strconv.FormatBool(ck.Encryption),
		"postdata":         ck.PostData,
		"integrationids":   intListToCDString(ck.IntegrationIds),
		"tags":             ck.Tags,
		"probe_filters":    ck.ProbeFilters,
		"userids":          intListToCDString(ck.UserIds),
		"teamids":          intListToCDString(ck.TeamIds),
	}

	// Ignore zero values
	if ck.Port != 0 {
		m["port"] = strconv.Itoa(ck.Port)
	}

	if ck.SendNotificationWhenDown != 0 {
		m["sendnotificationwhendown"] = strconv.Itoa(ck.SendNotificationWhenDown)
	}

	// ShouldContain and ShouldNotContain are mutually exclusive.
	// But we must define one so they can be emptied if required.
	if ck.ShouldContain != "" {
		m["shouldcontain"] = ck.ShouldContain
	} else {
		m["shouldnotcontain"] = ck.ShouldNotContain
	}

	// Convert auth
	if ck.Username != "" {
		m["auth"] = fmt.Sprintf("%s:%s", ck.Username, ck.Password)
	}

	// Convert headers
	var headers []string
	for k := range ck.RequestHeaders {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for i, k := range headers {
		m[fmt.Sprintf("requestheader%d", i)] = fmt.Sprintf("%s:%s", k, ck.RequestHeaders[k])
	}

	return m
}

// Params returns a map of parameters for an HttpCheck that can be sent along
// with an HTTP POST request. They are the same than the Put params, but
// empty strings cleared out, to avoid Pingdom API reject the request.
func (ck *HttpCheck) PostParams() map[string]string {
	params := ck.PutParams()

	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}
	params["type"] = "http"

	return params
}

// Determine whether the HttpCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *HttpCheck) Valid() error {
	if ck.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	if ck.Hostname == "" {
		return fmt.Errorf("Invalid value for `Hostname`.  Must contain non-empty string")
	}

	if ck.Resolution != 1 && ck.Resolution != 5 && ck.Resolution != 15 &&
		ck.Resolution != 30 && ck.Resolution != 60 {
		return fmt.Errorf("Invalid value %v for `Resolution`.  Allowed values are [1,5,15,30,60].", ck.Resolution)
	}

	if ck.ShouldContain != "" && ck.ShouldNotContain != "" {
		return fmt.Errorf("`ShouldContain` and `ShouldNotContain` must not be declared at the same time")
	}

	return nil
}

// Params returns a map of parameters for a PingCheck that can be sent along
// with an HTTP PUT request
func (ck *PingCheck) PutParams() map[string]string {
	m := map[string]string{
		"name":             ck.Name,
		"host":             ck.Hostname,
		"resolution":       strconv.Itoa(ck.Resolution),
		"paused":           strconv.FormatBool(ck.Paused),
		"notifyagainevery": strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup": strconv.FormatBool(ck.NotifyWhenBackup),
		"integrationids":   intListToCDString(ck.IntegrationIds),
		"probe_filters":    ck.ProbeFilters,
		"userids":          intListToCDString(ck.UserIds),
		"teamids":          intListToCDString(ck.TeamIds),
	}

	if ck.SendNotificationWhenDown != 0 {
		m["sendnotificationwhendown"] = strconv.Itoa(ck.SendNotificationWhenDown)
	}

	return m
}

// Params returns a map of parameters for a PingCheck that can be sent along
// with an HTTP POST request. Same as PUT.
func (ck *PingCheck) PostParams() map[string]string {
	params := ck.PutParams()
	params["type"] = "ping"
	return params
}

// Determine whether the PingCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *PingCheck) Valid() error {
	if ck.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	if ck.Hostname == "" {
		return fmt.Errorf("Invalid value for `Hostname`.  Must contain non-empty string")
	}

	if ck.Resolution != 1 && ck.Resolution != 5 && ck.Resolution != 15 &&
		ck.Resolution != 30 && ck.Resolution != 60 {
		return fmt.Errorf("Invalid value %v for `Resolution`.  Allowed values are [1,5,15,30,60].", ck.Resolution)
	}
	return nil
}

// Params returns a map of parameters for a TCPCheck that can be sent along
// with an HTTP PUT request
func (ck *TCPCheck) PutParams() map[string]string {
	m := map[string]string{
		"name":             ck.Name,
		"host":             ck.Hostname,
		"resolution":       strconv.Itoa(ck.Resolution),
		"paused":           strconv.FormatBool(ck.Paused),
		"notifyagainevery": strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup": strconv.FormatBool(ck.NotifyWhenBackup),
		"integrationids":   intListToCDString(ck.IntegrationIds),
		"probe_filters":    ck.ProbeFilters,
		"userids":          intListToCDString(ck.UserIds),
		"teamids":          intListToCDString(ck.TeamIds),
		"port":             strconv.Itoa(ck.Port),
	}

	if ck.SendNotificationWhenDown != 0 {
		m["sendnotificationwhendown"] = strconv.Itoa(ck.SendNotificationWhenDown)
	}

	if ck.StringToSend != "" {
		m["stringtosend"] = ck.StringToSend
	}

	if ck.StringToExpect != "" {
		m["stringtoexpect"] = ck.StringToExpect
	}

	return m
}

// Params returns a map of parameters for a TCPCheck that can be sent along
// with an HTTP POST request. Same as PUT.
func (ck *TCPCheck) PostParams() map[string]string {
	params := ck.PutParams()
	params["type"] = "tcp"
	return params
}

// Determine whether the TCPCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *TCPCheck) Valid() error {
	if ck.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	if ck.Hostname == "" {
		return fmt.Errorf("Invalid value for `Hostname`.  Must contain non-empty string")
	}

	if ck.Resolution != 1 && ck.Resolution != 5 && ck.Resolution != 15 &&
		ck.Resolution != 30 && ck.Resolution != 60 {
		return fmt.Errorf("Invalid value %v for `Resolution`.  Allowed values are [1,5,15,30,60].", ck.Resolution)
	}

	if ck.Port < 1 {
		return fmt.Errorf("Invalid value for `Port`.  Must contain an integer >= 1")
	}

	return nil
}

func intListToCDString(integers []int) string {
	var CDString string
	for i, item := range integers {
		if i == 0 {
			CDString = strconv.Itoa(item)
		} else {
			CDString = fmt.Sprintf("%v,%d", CDString, item)
		}
	}
	return CDString
}
