package pingdom

import (
	"encoding/json"
	"fmt"
)

// PingdomResponse represents a general response from the Pingdom API.
type PingdomResponse struct {
	Message string `json:"message"`
}

// PingdomError represents an error response from the Pingdom API.
type PingdomError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

// CheckResponse represents the JSON response for a check from the Pingdom API.
type CheckResponse struct {
	ID                       int                 `json:"id"`
	Name                     string              `json:"name"`
	Resolution               int                 `json:"resolution,omitempty"`
	SendNotificationWhenDown int                 `json:"sendnotificationwhendown,omitempty"`
	NotifyAgainEvery         int                 `json:"notifyagainevery,omitempty"`
	NotifyWhenBackup         bool                `json:"notifywhenbackup,omitempty"`
	Created                  int64               `json:"created,omitempty"`
	Hostname                 string              `json:"hostname,omitempty"`
	Status                   string              `json:"status,omitempty"`
	LastErrorTime            int64               `json:"lasterrortime,omitempty"`
	LastTestTime             int64               `json:"lasttesttime,omitempty"`
	LastResponseTime         int64               `json:"lastresponsetime,omitempty"`
	Paused                   bool                `json:"paused,omitempty"`
	IntegrationIds           []int               `json:"integrationids,omitempty"`
	Type                     CheckResponseType   `json:"type,omitempty"`
	Tags                     []CheckResponseTag  `json:"tags,omitempty"`
	UserIds                  []int               `json:"userids,omitempty"`
	Teams                    []CheckTeamResponse `json:"teams,omitempty"`
	ResponseTimeThreshold    int                 `json:"responsetime_threshold,omitempty"`

	// Legacy; this is not returned by the API, we backfill the value from the
	// Teams field.
	TeamIds []int
}

// CheckTeamResponse is a Team returned inside of a Check instance. (We can't
// use TeamResponse because the ID returned here is an int, not a string).
type CheckTeamResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CheckResponseType is the type of the Pingdom check.
type CheckResponseType struct {
	Name string                    `json:"-"`
	HTTP *CheckResponseHTTPDetails `json:"http,omitempty"`
	TCP  *CheckResponseTCPDetails  `json:"tcp,omitempty"`
}

// CheckResponseTag is an optional tag that can be added to checks.
type CheckResponseTag struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Count interface{} `json:"count"`
}

// MaintenanceResponse represents the JSON response for a maintenance from the Pingdom API.
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

// MaintenanceCheckResponse represents Check reply in json MaintenanceResponse.
type MaintenanceCheckResponse struct {
	Uptime []int `json:"uptime"`
	Tms    []int `json:"tms"`
}

// ProbeResponse represents the JSON response for probes from the Pingdom API.
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

// TeamResponse represents the JSON response for teams from the Pingdom API.
type TeamResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Users []TeamUserResponse
}

// TeamUserResponse represents the JSON response for users in teams from the Pingdom API.
type TeamUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TeamDeleteResponse represents the JSON response for delete team from the Pingdom API.
type TeamDeleteResponse struct {
	Success bool `json:"success"`
}

// PublicReportResponse represents the JSON response for a public report from the Pingdom API.
type PublicReportResponse struct {
	ID        int    `json:"checkid"`
	Name      string `json:"checkname"`
	ReportURL string `json:"reporturl"`
}

// SummaryPerformanceResponse represents the JSON response for a summary performance from the Pingdom API.
type SummaryPerformanceResponse struct {
	Summary SummaryPerformanceMap `json:"summary"`
}

// SummaryPerformanceMap is the performance broken down over different time intervals.
type SummaryPerformanceMap struct {
	Hours []SummaryPerformanceSummary `json:"hours,omitempty"`
	Days  []SummaryPerformanceSummary `json:"days,omitempty"`
	Weeks []SummaryPerformanceSummary `json:"weeks,omitempty"`
}

// SummaryPerformanceSummary is the metrics for a performance summary.
type SummaryPerformanceSummary struct {
	AvgResponse int `json:"avgresponse"`
	Downtime    int `json:"downtime"`
	StartTime   int `json:"starttime"`
	Unmonitored int `json:"unmonitored"`
	Uptime      int `json:"uptime"`
}

// UserSmsResponse represents the JSON response for a user SMS contact.
type UserSmsResponse struct {
	Id          int    `json:"id"`
	Severity    string `json:"severity"`
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Provider    string `json:"provider"`
}

// UserEmailResponse represents the JSON response for a user email contact.
type UserEmailResponse struct {
	Id       int    `json:"id"`
	Severity string `json:"severity"`
	Address  string `json:"address"`
}

// CreateUserContactResponse represents the JSON response for a user contact.
type CreateUserContactResponse struct {
	Id int `json:"id"`
}

// UsersResponse represents the JSON response for a Pingom User.
type UsersResponse struct {
	Id       int                 `json:"id"`
	Paused   string              `json:"paused,omitempty"`
	Username string              `json:"name,omitempty"`
	Sms      []UserSmsResponse   `json:"sms,omitempty"`
	Email    []UserEmailResponse `json:"email,omitempty"`
}

// UnmarshalJSON converts a byte array into a CheckResponseType.
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
		c.TCP = rawCheckDetails.TCP
	}
	return nil
}

// CheckResponseHTTPDetails represents the details specific to HTTP checks.
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

// CheckResponseTCPDetails represents the details specific to TCP checks.
type CheckResponseTCPDetails struct {
	Port           int    `json:"port,omitempty"`
	StringToSend   string `json:"stringtosend,omitempty"`
	StringToExpect string `json:"stringtoexpect,omitempty"`
}

// Return string representation of the PingdomError.
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall JSON responses from Pingdom.

type listChecksJSONResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type listMaintenanceJSONResponse struct {
	Maintenances []MaintenanceResponse `json:"maintenance"`
}

type listProbesJSONResponse struct {
	Probes []ProbeResponse `json:"probes"`
}

type listTeamsJSONResponse struct {
	Teams []TeamResponse `json:"teams"`
}

type listPublicReportsJSONResponse struct {
	Checks []PublicReportResponse `json:"public"`
}

type checkDetailsJSONResponse struct {
	Check *CheckResponse `json:"check"`
}

type maintenanceDetailsJSONResponse struct {
	Maintenance *MaintenanceResponse `json:"maintenance"`
}

type teamDetailsJSONResponse struct {
	Team *TeamResponse `json:"team"`
}

type createUserContactJSONResponse struct {
	Contact *CreateUserContactResponse `json:"contact_target"`
}

type createUserJSONResponse struct {
	User *UsersResponse `json:"user"`
}

type listUsersJSONResponse struct {
	Users []UsersResponse `json:"users"`
}

type errorJSONResponse struct {
	Error *PingdomError `json:"error"`
}
