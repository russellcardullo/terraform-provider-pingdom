package pingdom

import (
	"encoding/json"
	"io/ioutil"
)

// ProbeService provides an interface to Pingdom probes.
type ProbeService struct {
	client *Client
}

// List return a list of probes from Pingdom.
func (cs *ProbeService) List(params ...map[string]string) ([]ProbeResponse, error) {
	param := map[string]string{}
	if len(params) == 1 {
		param = params[0]
	}
	req, err := cs.client.NewRequest("GET", "/probes", param)
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

	p := &listProbesJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &p)

	return p.Probes, err
}
