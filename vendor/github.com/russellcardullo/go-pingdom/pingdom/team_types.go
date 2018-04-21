package pingdom

import (
	"fmt"
)

// TeamData represents a Pingdom Team Data.
type TeamData struct {
	Name    string `json:"name"`
	UserIds string `json:"userids,omitempty"`
}

// PutParams returns a map of parameters for an Team that can be sent along.
func (ck *TeamData) PutParams() map[string]string {
	t := map[string]string{
		"name": ck.Name,
	}

	// Ignore if not defined
	if ck.UserIds != "" {
		t["userids"] = ck.UserIds
	}

	return t
}

// PostParams returns a map of parameters for an Team that can be sent along
// with an HTTP POST request. They are the same than the Put params, but
// empty strings cleared out, to avoid Pingdom API reject the request.
func (ck *TeamData) PostParams() map[string]string {
	params := ck.PutParams()

	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}

	return params
}

// Valid Determine whether the Team contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API.
func (ck *TeamData) Valid() error {
	if ck.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	return nil
}
