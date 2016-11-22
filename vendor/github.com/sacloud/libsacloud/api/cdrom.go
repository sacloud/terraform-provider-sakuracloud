package api

import (
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
)

// CDROMAPI ISOイメージAPI
type CDROMAPI struct {
	*baseAPI
}

// NewCDROMAPI ISOイメージAPI新規作成
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

// Create 新規作成
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

// OpenFTP FTP接続開始
func (api *CDROMAPI) OpenFTP(id int64, reset bool) (*sacloud.FTPServer, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
		body   = map[string]bool{"ChangePassword": reset}
		res    = &sacloud.Response{}
	)

	result, err := api.action(method, uri, body, res)
	if !result || err != nil {
		return nil, err
	}

	return res.FTPServer, nil
}

// CloseFTP FTP接続終了
func (api *CDROMAPI) CloseFTP(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)

}

// SleepWhileCopying コピー終了まで待機
func (api *CDROMAPI) SleepWhileCopying(id int64, timeout time.Duration) error {

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
			return fmt.Errorf("Timeout: SleepWhileCopying[disk:%d]", id)
		}
	}
}
