package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// CheckService provides an interface to Pingdom checks
type CheckService struct {
	client *Client
}

// Check is an interface representing a pingdom check.
// Specific check types should implement the methods of this interface
type Check interface {
	PutParams() map[string]string
	PostParams() map[string]string
	Valid() error
}

// Return a list of checks from Pingdom.
// This returns type CheckResponse rather than Check since the
// pingdom API does not return a complete representation of a check.
func (cs *CheckService) List() ([]CheckResponse, error) {
	req, err := cs.client.NewRequest("GET", "/api/2.0/checks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &listChecksJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Checks, err
}

// Create a new check. This function will validate the given check param
// to ensure that it contains correct values before submitting the request
// Returns a CheckResponse object representing the response from Pingdom.
// Note that Pingdom does not return a full check object so in the returned
// object you should only use the ID field.
func (cs *CheckService) Create(check Check) (*CheckResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/api/2.0/checks", check.PostParams())
	if err != nil {
		return nil, err
	}

	m := &checkDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Check, err
}

// ReadCheck returns detailed information about a pingdom check given its ID.
// This returns type CheckResponse rather than Check since the
// pingdom API does not return a complete representation of a check.
func (cs *CheckService) Read(id int) (*CheckResponse, error) {
	req, err := cs.client.NewRequest("GET", "/api/2.0/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &checkDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m.Check, err
}

// UpdateCheck will update the check represented by the given ID with the values
// in the given check.  You should submit the complete list of values in
// the given check parameter, not just those that have changed.
func (cs *CheckService) Update(id int, check Check) (*PingdomResponse, error) {
	if err := check.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("PUT", "/api/2.0/checks/"+strconv.Itoa(id), check.PutParams())
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}

// DeleteCheck will delete the check for the given ID.
func (cs *CheckService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/api/2.0/checks/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &PingdomResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m, err
}
