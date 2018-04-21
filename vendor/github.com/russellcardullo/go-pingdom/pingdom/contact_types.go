package pingdom

import (
	"fmt"
	"strconv"
)

// Contact represents a Pingdom contact
type Contact struct {
	Name               string `json:"name"`
	Email              string `json:"email,omitempty"`
	Cellphone          string `json:"cellphone,omitempty"`
	CountryISO         string `json:"countryiso,omitempty"`
	CountryCode        string `json:"countrycode,omitempty"`
	DefaultSMSProvider string `json:"defaultsmsprovider,omitempty"`
	DirectTwitter      bool   `json:"directtwitter,omitempty"`
	TwitterUser        string `json:"twitteruser,omitempty"`
	IphoneTokens       string `json:"iphonetokens,omitempty"`
	AndroidTokens      string `json:"androidtokens,omitempty"`
	Paused             bool   `json:"paused,omitempty"`
}

// Params returns a map of parameters for a Contact that can be sent along
// with an HTTP PUT request
func (ct *Contact) PutParams() map[string]string {
	m := map[string]string{
		"name":               ct.Name,
		"email":              ct.Email,
		"cellphone":          ct.Cellphone,
		"countryiso":         ct.CountryISO,
		"countrycode":        ct.CountryCode,
		"defaultsmsprovider": ct.DefaultSMSProvider,
		"directtwitter":      strconv.FormatBool(ct.DirectTwitter),
		"twitteruser":        ct.TwitterUser,
	}

	return m
}

// Params returns a map of parameters for a Contact that can be sent along
// with an HTTP POST request. They are the same than the Put params, but
// empty strings cleared out, to avoid Pingdom API reject the request.
func (ct *Contact) PostParams() map[string]string {
	params := ct.PutParams()

	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}

	return params
}

// Determine whether the Contact contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ct *Contact) Valid() error {
	if ct.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	return nil
}
