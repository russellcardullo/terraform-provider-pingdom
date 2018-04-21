package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// MaintenanceService provides an interface to Pingdom maintenance windows.
type MaintenanceService struct {
	client *Client
}

// Maintenance is a Pingdom maintenance window.
type Maintenance interface {
	PutParams() map[string]string
	PostParams() map[string]string
	Valid() error
}

// MaintenanceDelete is the set of parameters to a Pingdom maintenance delete request.
type MaintenanceDelete interface {
	DeleteParams() map[string]string
	ValidDelete() error
}

// List returns the response holding a list of Maintenance windows.
func (cs *MaintenanceService) List(params ...map[string]string) ([]MaintenanceResponse, error) {
	param := map[string]string{}
	if len(params) != 0 {
		for _, m := range params {
			for k, v := range m {
				param[k] = v
			}
		}
	}
	req, err := cs.client.NewRequest("GET", "/maintenance", param)
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

	m := &listMaintenanceJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Maintenances, err
}

// Read returns a Maintenance for a given ID.
func (cs *MaintenanceService) Read(id int) (*MaintenanceResponse, error) {
	req, err := cs.client.NewRequest("GET", "/maintenance/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &maintenanceDetailsJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m.Maintenance, err
}

// Create creates a new Maintenance.
func (cs *MaintenanceService) Create(maintenance Maintenance) (*MaintenanceResponse, error) {
	if err := maintenance.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/maintenance", maintenance.PostParams())
	if err != nil {
		return nil, err
	}

	m := &maintenanceDetailsJSONResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Maintenance, err
}

// Update is used to update an existing Maintenance. Only the 'Description',
// and 'To' fields can be updated.
func (cs *MaintenanceService) Update(id int, maintenance Maintenance) (*PingdomResponse, error) {
	if err := maintenance.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("PUT", "/maintenance/"+strconv.Itoa(id), maintenance.PutParams())
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

// MultiDelete will delete the Maintenance for the given ID.
func (cs *MaintenanceService) MultiDelete(maintenance MaintenanceDelete) (*PingdomResponse, error) {
	if err := maintenance.ValidDelete(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("DELETE", "/maintenance/", maintenance.DeleteParams())
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

// Delete will delete the Maintenance for the given ID.
func (cs *MaintenanceService) Delete(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/maintenance/"+strconv.Itoa(id), nil)
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
