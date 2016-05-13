package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

type CDROMAPI struct {
	*baseAPI
}

func NewCDROMAPI(client *Client) *CDROMAPI {
	return &CDROMAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "cdrom"
			},
		},
	}
}

func (api *CDROMAPI) Create(value *sacloud.CDROM) (*sacloud.CDROM, *sacloud.FTPServer, error) {
	f := func(res *sacloud.Response) error {
		return api.create(api.createRequest(value), res)
	}
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, nil, err
	}
	return res.CDROM, res.FTPServer, nil
}

func (api *CDROMAPI) OpenFTP(id string, reset bool) (*sacloud.FTPServer, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/ftp", api.getResourceURL(), id)
		body   = map[string]bool{"ChangePassword": reset}
		res    = &sacloud.Response{}
	)

	result, err := api.action(method, uri, body, res)
	if !result || err != nil {
		return nil, err
	}

	return res.FTPServer, nil
}

func (api *CDROMAPI) CloseFTP(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/ftp", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)

}

func (api *CDROMAPI) SleepWhileCopying(id string, timeout time.Duration) error {

	current := 0 * time.Second
	interval := 5 * time.Second
	for {
		archive, err := api.Read(id)
		if err != nil {
			return err
		}

		if archive.IsAvailable() {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepWhileCopying[disk:%s]", id)
		}
	}
}
