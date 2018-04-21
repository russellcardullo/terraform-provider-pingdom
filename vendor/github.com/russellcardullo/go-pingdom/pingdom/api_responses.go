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
	SendToAndroid            bool               `json:"sendtoandroid,omitempty"`
	SendToEmail              bool               `json:"sendtoemail,omitempty"`
	SendToIPhone             bool               `json:"sendtoiphone,omitempty"`
	SendToSms                bool               `json:"sendtosms,omitempty"`
	SendToTwitter            bool               `json:"sendtotwitter,omitempty"`
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
	ContactIds               []int              `json:"contactids,omitempty"`
	IntegrationIds           []int              `json:"integrationids,omitempty"`
	Type                     CheckResponseType  `json:"type,omitempty"`
	Tags                     []CheckResponseTag `json:"tags,omitempty"`
}

type CheckResponseType struct {
	Name string                    `json:"-"`
	HTTP *CheckResponseHTTPDetails `json:"http,omitempty"`
}

type CheckResponseTag struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Count string `json:"count"`
}

type ContactResponse struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email,omitempty"`
	Cellphone          string `json:"cellphone,omitempty"`
	CountryISO         string `json:"countryiso,omitempty"`
	DefaultSMSProvider string `json:"defaultsmsprovider,omitempty"`
	DirectTwitter      bool   `json:"directtwitter,omitempty"`
	TwitterUser        string `json:"twitteruser,omitempty"`
	IphoneTokens       string `json:"iphonetokens,omitempty"`
	AndroidTokens      string `json:"androidtokens,omitempty"`
	Paused             bool   `json:"paused,omitempty"`
	Type               string `json:"type,omitempty"`
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
		for k, _ := range v {
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

// Return string representation of the PingdomError
func (r *PingdomError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

// private types used to unmarshall json responses from pingdom

type listChecksJsonResponse struct {
	Checks []CheckResponse `json:"checks"`
}

type checkDetailsJsonResponse struct {
	Check *CheckResponse `json:"check"`
}

type contactDetailsJsonResponse struct {
	Contact *ContactResponse `json:"contact"`
}

type listContactsJsonResponse struct {
	Contacts []ContactResponse `json:"contacts"`
}

type errorJsonResponse struct {
	Error *PingdomError `json:"error"`
}
