package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

type ArchiveAPI struct {
	*baseAPI
}

func NewArchiveAPI(client *Client) *ArchiveAPI {
	return &ArchiveAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "archive"
			},
		},
	}
}

func (api *ArchiveAPI) OpenFTP(id string, reset bool) (*sacloud.FTPServer, error) {
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

func (api *ArchiveAPI) CloseFTP(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/ftp", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)

}

func (api *ArchiveAPI) SleepWhileCopying(id string, timeout time.Duration) error {

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
