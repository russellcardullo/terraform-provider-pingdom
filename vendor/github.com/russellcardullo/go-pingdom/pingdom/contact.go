package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// ContactService provides an interface to Pingdom notification contacts
type ContactService struct {
	client *Client
}

// Return a list of contacts from Pingdom.
// This returns type ContactResponse rather than Contact since the
// pingdom API does not return a complete representation of a contact.
func (cs *ContactService) List() ([]ContactResponse, error) {
	req, err := cs.client.NewRequest("GET", "/api/2.0/notification_contacts", nil)
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
	m := &listContactsJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Contacts, err
}

// Create a new contact.  This function will validate the given contact param
// to ensure that it contains correct values before submitting the request
// Returns a ContactResponse object representing the response from Pingdom.
func (cs *ContactService) Create(contact *Contact) (*ContactResponse, error) {
	if err := contact.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/api/2.0/notification_contacts", contact.PostParams())
	if err != nil {
		return nil, err
	}

	m := &contactDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Contact, err
}

// UpdateContact will update the contact represented by the given ID with the values
// in the given contact.  You should submit the complete list of values in
// the given check parameter, not just those that have changed.
func (cs *ContactService) Update(id int, contact *Contact) (*PingdomResponse, error) {
	if err := contact.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("PUT", "/api/2.0/notification_contacts/"+strconv.Itoa(id), contact.PutParams())
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
func (cs *ContactService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/api/2.0/notification_contacts/"+strconv.Itoa(id), nil)
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
