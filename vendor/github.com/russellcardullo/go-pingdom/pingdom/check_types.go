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
	SendToAndroid            bool              `json:"sendtoandroid,omitempty"`
	SendToEmail              bool              `json:"sendtoemail,omitempty"`
	SendToIPhone             bool              `json:"sendtoiphone,omitempty"`
	SendToSms                bool              `json:"sendtosms,omitempty"`
	SendToTwitter            bool              `json:"sendtotwitter,omitempty"`
	SendNotificationWhenDown int               `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int               `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool              `json:"notifywhenbackup,omitempty"`
	UseLegacyNotifications   bool              `json:"use_legacy_notifications,omitempty"`
	Url                      string            `json:"url,omitempty"`
	Encryption               bool              `json:"encryption,omitempty"`
	Port                     int               `json:"port,omitempty"`
	Username                 string            `json:"username,omitempty"`
	Password                 string            `json:"password,omitempty"`
	ShouldContain            string            `json:"shouldcontain,omitempty"`
	ShouldNotContain         string            `json:"shouldnotcontain,omitempty"`
	PostData                 string            `json:"postdata,omitempty"`
	RequestHeaders           map[string]string `json:"requestheaders,omitempty"`
	ContactIds               []int             `json:"contactids,omitempty"`
	IntegrationIds           []int             `json:"integrationids,omitempty"`
	Tags                     string            `json:"tags,omitempty"`
  ProbeFilters             string            `json:"probe_filters,omitempty"`
}

// PingCheck represents a Pingdom ping check
type PingCheck struct {
	Name                     string `json:"name"`
	Hostname                 string `json:"hostname,omitempty"`
	Resolution               int    `json:"resolution,omitempty"`
	Paused                   bool   `json:"paused,omitempty"`
	SendToAndroid            bool   `json:"sendtoandroid,omitempty"`
	SendToEmail              bool   `json:"sendtoemail,omitempty"`
	SendToIPhone             bool   `json:"sendtoiphone,omitempty"`
	SendToSms                bool   `json:"sendtosms,omitempty"`
	SendToTwitter            bool   `json:"sendtotwitter,omitempty"`
	SendNotificationWhenDown int    `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int    `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool   `json:"notifywhenbackup,omitempty"`
	UseLegacyNotifications   bool   `json:"use_legacy_notifications,omitempty"`
	ContactIds               []int  `json:"contactids,omitempty"`
	IntegrationIds           []int  `json:"integrationids,omitempty"`
  ProbeFilters             string `json:"probe_filters,omitempty"`
}

// Params returns a map of parameters for an HttpCheck that can be sent along
// with an HTTP PUT request
func (ck *HttpCheck) PutParams() map[string]string {
	m := map[string]string{
		"name":                     ck.Name,
		"host":                     ck.Hostname,
		"resolution":               strconv.Itoa(ck.Resolution),
		"paused":                   strconv.FormatBool(ck.Paused),
		"sendtoemail":              strconv.FormatBool(ck.SendToEmail),
		"sendtosms":                strconv.FormatBool(ck.SendToSms),
		"sendtotwitter":            strconv.FormatBool(ck.SendToTwitter),
		"sendtoiphone":             strconv.FormatBool(ck.SendToIPhone),
		"sendtoandroid":            strconv.FormatBool(ck.SendToAndroid),
		"sendnotificationwhendown": strconv.Itoa(ck.SendNotificationWhenDown),
		"notifyagainevery":         strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup":         strconv.FormatBool(ck.NotifyWhenBackup),
		"use_legacy_notifications": strconv.FormatBool(ck.UseLegacyNotifications),
		"url":        ck.Url,
		"encryption": strconv.FormatBool(ck.Encryption),
		"postdata":   ck.PostData,
		"contactids": intListToCDString(ck.ContactIds),
		"integrationids": intListToCDString(ck.IntegrationIds),
		"tags":       ck.Tags,
    "probe_filters": ck.ProbeFilters,
	}

	// Ignore port is not defined
	if ck.Port != 0 {
		m["port"] = strconv.Itoa(ck.Port)
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
	return map[string]string{
		"name":                     ck.Name,
		"host":                     ck.Hostname,
		"resolution":               strconv.Itoa(ck.Resolution),
		"paused":                   strconv.FormatBool(ck.Paused),
		"sendtoemail":              strconv.FormatBool(ck.SendToEmail),
		"sendtosms":                strconv.FormatBool(ck.SendToSms),
		"sendtotwitter":            strconv.FormatBool(ck.SendToTwitter),
		"sendtoiphone":             strconv.FormatBool(ck.SendToIPhone),
		"sendtoandroid":            strconv.FormatBool(ck.SendToAndroid),
		"sendnotificationwhendown": strconv.Itoa(ck.SendNotificationWhenDown),
		"notifyagainevery":         strconv.Itoa(ck.NotifyAgainEvery),
		"notifywhenbackup":         strconv.FormatBool(ck.NotifyWhenBackup),
		"use_legacy_notifications": strconv.FormatBool(ck.UseLegacyNotifications),
		"contactids":               intListToCDString(ck.ContactIds),
		"integrationids":           intListToCDString(ck.IntegrationIds),
    "probe_filters":            ck.ProbeFilters,
	}
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
