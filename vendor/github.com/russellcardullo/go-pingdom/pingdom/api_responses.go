package pingdom

import (
	"encoding/json"
	"fmt"
)

// PingdomResponse represents a general response from the Pingdom API
type PingdomResponse struct {
	Message string `json:"message"`
}

// PingdomError represents an error response from the Pingdom API
type PingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

// CheckResponse represents the json response for a check from the Pingdom API
type CheckResponse struct {
	ID                       int                `json:"id"`
	Name                     string             `json:"name"`
	Resolution               int                `json:"resolution,omitempty"`
	SendNotificationWhenDown int                `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int                `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool               `json:"notifywhenbackup,omitempty"`
	Created                  int64              `json:"created,omitempty"`
	Hostname                 string             `json:"hostname,omitempty"`
	Status                   string             `json:"status,omitempty"`
	LastErrorTime            int64              `json:"lasterrortime,omitempty"`
	LastTestTime             int64              `json:"lasttesttime,omitempty"`
	LastResponseTime         int64              `json:"lastresponsetime,omitempty"`
	Paused                   bool               `json:"paused,omitempty"`
	IntegrationIds           []int              `json:"integrationids,omitempty"`
	Type                     CheckResponseType  `json:"type,omitempty"`
	Tags                     []CheckResponseTag `json:"tags,omitempty"`
	UserIds                  []int              `json:"userids,omitempty"`
	TeamIds                  []int              `json:"teamids,omitempty"`
}

type CheckResponseType struct {
	Name string                    `json:"-"`
	HTTP *CheckResponseHTTPDetails `json:"http,omitempty"`
	TCP  *CheckResponseTCPDetails  `json:"tcp,omitempty"`
}

type CheckResponseTag struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Count interface{} `json:"count"`
}

// MaintenanceResponse represents the json response for a maintenance from the Pingdom API
type MaintenanceResponse struct {
	ID             int                      `json:"id"`
	Description    string                   `json:"description"`
	From           int64                    `json:"from"`
	To             int64                    `json:"to"`
	RecurrenceType string                   `json:"recurrencetype"`
	RepeatEvery    int                      `json:"repeatevery"`
	EffectiveTo    int64                    `json:"effectiveto"`
	Checks         MaintenanceCheckResponse `json:"checks"`
}

// MaintenanceCheckResponse represents Check reply in json MaintenanceResponse
type MaintenanceCheckResponse struct {
	Uptime []int `json:"uptime"`
	Tms    []int `json:"tms"`
}

// ProbeResponse represents the json response for probes from the PIngdom API
type ProbeResponse struct {
	ID         int    `json:"id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	IPv6       string `json:"ipv6"`
	CountryISO string `json:"countryiso"`
	Region     string `json:"region"`
}

// TeamResponse represents the json response for teams from the PIngdom API
type TeamResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Users []TeamUserResponse
}

// TeamUserResponse represents the json response for users in teams from the PIngdom API
type TeamUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TeamDeleteResponse represents the json response for delete team from the PIngdom API
type TeamDeleteResponse struct {
	Success bool `json:"success"`
}

type PublicReportResponse struct {
	ID        int    `json:"checkid"`
	Name      string `json:"checkname"`
	ReportURL string `json:"reporturl"`
}

func (c *CheckResponseType) UnmarshalJSON(b []byte) error {
	var raw interface{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	switch v := raw.(type) {
	case string:
		c.Name = v
	case map[string]interface{}:
		if len(v) != 1 {
			return fmt.Errorf("Check detailed response `check.type` contains more than one object: %+v", v)
		}
		for k := range v {
			c.Name = k
		}

		// Allow continue use json.Unmarshall using a type != Unmarshaller
		// This avoid enter in a infinite loop
		type t CheckResponseType
		var rawCheckDetails t

		err := json.Unmarshal(b, &rawCheckDetails)
		if err != nil {
			return err
		}
		c.HTTP = rawCheckDetails.HTTP
	}
	return nil
}

// HttpCheck represents a Pingdom http check.
type CheckResponseHTTPDetails struct {
	Url              string            `json:"url,omitempty"`
	Encryption       bool              `json:"encryption,omitempty"`
	Port             int               `json:"port,omitempty"`
	Username         string            `json:"username,omitempty"`
	Password         string            `json:"password,omitempty"`
	ShouldContain    string            `json:"shouldcontain,omitempty"`
	ShouldNotContain string            `json:"shouldnotcontain,omitempty"`
	PostData         string            `json:"postdata,omitempty"`
	RequestHeaders   map[string]string `json:"requestheaders,omitempty"`
}

// HttpCheck represents a Pingdom http check.
type CheckResponseTCPDetails struct {
	Port           int    `json:"port,omitempty"`
	StringToSend   string `json:"stringtosend,omitempty"`
	StringToExpect string `json:"stringtoexpect,omitempty"`
}

// Return string representation of the PingdomError
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall json responses from pingdom

type listChecksJsonResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type listMaintenanceJsonResponse struct {
	Maintenances []MaintenanceResponse `json:"maintenance"`
}

type listProbesJsonResponse struct {
	Probes []ProbeResponse `json:"probes"`
}

type listTeamsJsonResponse struct {
	Teams []TeamResponse `json:"teams"`
}

type listPublicReportsJsonResponse struct {
	Checks []PublicReportResponse `json:"public"`
}

type checkDetailsJsonResponse struct {
	Check *CheckResponse `json:"check"`
}

type maintenanceDetailsJsonResponse struct {
	Maintenance *MaintenanceResponse `json:"maintenance"`
}

type teamDetailsJsonResponse struct {
	Team *TeamResponse `json:"team"`
}

type errorJsonResponse struct {
	Error *PingdomError `json:"error"`
}
