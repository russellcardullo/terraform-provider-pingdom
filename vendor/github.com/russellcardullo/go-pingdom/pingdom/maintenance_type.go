package pingdom

import (
	"fmt"
	"strconv"
)

// MaintenanceWindow represents a Pingdom Maintenance Window.
type MaintenanceWindow struct {
	Description    string `json:"description"`
	From           int64  `json:"from"`
	To             int64  `json:"to"`
	RecurrenceType string `json:"recurrencetype,omitempty"`
	RepeatEvery    int    `json:"repeatevery,omitempty"`
	EffectiveTo    int    `json:"effectiveto,omitempty"`
	UptimeIDs      string `json:"uptimeids,omitempty"`
	TmsIDs         string `json:"tmsids,omitempty"`
}

// MaintenanceWindowDelete represents delete request parameters
type MaintenanceWindowDelete struct {
	MaintenanceIDs string `json:"maintenanceids"`
}

// PutParams returns a map of parameters for an MaintenanceWindow that can be sent along
func (ck *MaintenanceWindow) PutParams() map[string]string {
	m := map[string]string{
		"description": ck.Description,
		"from":        strconv.FormatInt(ck.From, 10),
		"to":          strconv.FormatInt(ck.To, 10),
	}

	// Ignore if not defined
	if ck.RecurrenceType != "" {
		m["recurrencetype"] = ck.RecurrenceType
	}

	if ck.UptimeIDs != "" {
		m["uptimeids"] = ck.UptimeIDs
	}

	if ck.TmsIDs != "" {
		m["tmsids"] = ck.TmsIDs
	}

	if ck.RepeatEvery != 0 {
		m["repeatevery"] = strconv.Itoa(ck.RepeatEvery)
	}

	if ck.EffectiveTo != 0 {
		m["effectiveto"] = strconv.Itoa(ck.EffectiveTo)
	}

	return m
}

// PostParams returns a map of parameters for an Maintenance Window that can be sent along
// with an HTTP POST request. They are the same than the Put params, but
// empty strings cleared out, to avoid Pingdom API reject the request.
func (ck *MaintenanceWindow) PostParams() map[string]string {
	params := ck.PutParams()

	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}

	return params
}

// Valid Determine whether the MaintenanceWindow contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *MaintenanceWindow) Valid() error {
	if ck.Description == "" {
		return fmt.Errorf("Invalid value for `Description`.  Must contain non-empty string")
	}

	if ck.From == 0 {
		return fmt.Errorf("Invalid value for `From`.  Must contain time")
	}

	if ck.To == 0 {
		return fmt.Errorf("Invalid value for `To`.  Must contain time")
	}

	return nil
}

// DeleteParams returns a map of parameters for an MaintenanceWindow that can be sent along
func (ck *MaintenanceWindowDelete) DeleteParams() map[string]string {
	m := map[string]string{
		"maintenanceids": ck.MaintenanceIDs,
	}

	return m
}

func (ck *MaintenanceWindowDelete) ValidDelete() error {
	if ck.MaintenanceIDs == "" {
		return fmt.Errorf("Invalid value for `IDs`.  Must contain non-empty string")
	}

	return nil
}
