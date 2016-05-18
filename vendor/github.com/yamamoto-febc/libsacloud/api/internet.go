package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

type InternetAPI struct {
	*baseAPI
}

func NewInternetAPI(client *Client) *InternetAPI {
	return &InternetAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "internet"
			},
		},
	}
}

func (api *InternetAPI) UpdateBandWidth(id string, bandWidth int) (*sacloud.Internet, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/bandwidth", api.getResourceURL(), id)
		body   = &sacloud.Request{}
	)
	body.Internet = &sacloud.Internet{BandWidthMbps: bandWidth}

	return api.request(func(res *sacloud.Response) error {
		return api.baseAPI.request(method, uri, body, res)
	})
}

func (api *InternetAPI) SleepWhileCreating(internetID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second

	var item *sacloud.Internet
	var err error
	//READ
	for item == nil && timeout > current {
		item, err = api.Read(internetID)

		if err != nil {
			time.Sleep(interval)
			current = current + interval
			err = nil
		}
	}

	if err != nil {
		return err
	}
	if current > timeout {
		return fmt.Errorf("Timeout: Can't read /internet/%s", internetID)
	}

	return nil

}
