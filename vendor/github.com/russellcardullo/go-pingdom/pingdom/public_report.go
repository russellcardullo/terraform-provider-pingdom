package pingdom

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// PublicReportService provides an interface to Pingdom reports.
type PublicReportService struct {
	client *Client
}

// List return a list of reports from Pingdom.
func (cs *PublicReportService) List() ([]PublicReportResponse, error) {
	req, err := cs.client.NewRequest("GET", "/reports.public", nil)
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

	t := &listPublicReportsJSONResponse{}
	err = json.Unmarshal([]byte(bodyString), &t)

	return t.Checks, err
}

// PublishCheck is used to activate a check on the public report.
func (cs *PublicReportService) PublishCheck(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("PUT", "/reports.public/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	t := &PingdomResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

// WithdrawlCheck is used to deactivate a check on the public report.
func (cs *PublicReportService) WithdrawlCheck(id int) (*PingdomResponse, error) {
	req, err := cs.client.NewRequest("DELETE", "/reports.public/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}

	t := &PingdomResponse{}
	_, err = cs.client.Do(req, t)
	if err != nil {
		return nil, err
	}
	return t, err
}
