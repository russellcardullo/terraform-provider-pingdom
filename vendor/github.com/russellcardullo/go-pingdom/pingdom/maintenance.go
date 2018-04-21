package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

type MaintenanceService struct {
	client *Client
}

type Maintenance interface {
	PutParams() map[string]string
	PostParams() map[string]string
	Valid() error
}

type MaintenanceDelete interface {
	DeleteParams() map[string]string
	ValidDelete() error
}

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

	m := &listMaintenanceJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)

	return m.Maintenances, err
}

func (cs *MaintenanceService) Read(id int) (*MaintenanceResponse, error) {
	req, err := cs.client.NewRequest("GET", "/maintenance/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	m := &maintenanceDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}

	return m.Maintenance, err
}

func (cs *MaintenanceService) Create(maintenance Maintenance) (*MaintenanceResponse, error) {
	if err := maintenance.Valid(); err != nil {
		return nil, err
	}

	req, err := cs.client.NewRequest("POST", "/maintenance", maintenance.PostParams())
	if err != nil {
		return nil, err
	}

	m := &maintenanceDetailsJsonResponse{}
	_, err = cs.client.Do(req, m)
	if err != nil {
		return nil, err
	}
	return m.Maintenance, err
}

// Update is used to update existing maintenance window.Looks like you can only update 'Description',
// and 'To' fields.
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

	log.Println(maintenance.DeleteParams())

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
