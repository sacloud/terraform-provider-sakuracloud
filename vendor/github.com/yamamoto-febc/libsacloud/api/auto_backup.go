package api

import (
	"encoding/json"
	//	"strings"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

type SearchAutoBackupResponse struct {
	Total                        int                  `json:",omitempty"`
	From                         int                  `json:",omitempty"`
	Count                        int                  `json:",omitempty"`
	CommonServiceAutoBackupItems []sacloud.AutoBackup `json:"CommonServiceItems,omitempty"`
}

type autoBackupRequest struct {
	CommonServiceAutoBackupItem *sacloud.AutoBackup    `json:"CommonServiceItem,omitempty"`
	From                        int                    `json:",omitempty"`
	Count                       int                    `json:",omitempty"`
	Sort                        []string               `json:",omitempty"`
	Filter                      map[string]interface{} `json:",omitempty"`
	Exclude                     []string               `json:",omitempty"`
	Include                     []string               `json:",omitempty"`
}

type autoBackupResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.AutoBackup `json:"CommonServiceItem,omitempty"`
}

// AutoBackupAPI API Client for SAKURA CLOUD AutoBackup
type AutoBackupAPI struct {
	*baseAPI
}

func NewAutoBackupAPI(client *Client) *AutoBackupAPI {
	return &AutoBackupAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "commonserviceitem"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Provider.Class", "autobackup")
				return res
			},
		},
	}
}

func (api *AutoBackupAPI) Find() (*SearchAutoBackupResponse, error) {

	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchAutoBackupResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *AutoBackupAPI) request(f func(*autoBackupResponse) error) (*sacloud.AutoBackup, error) {
	res := &autoBackupResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.AutoBackup, nil
}

func (api *AutoBackupAPI) createRequest(value *sacloud.AutoBackup) *autoBackupResponse {
	return &autoBackupResponse{AutoBackup: value}
}

func (api *AutoBackupAPI) New(name string, diskID string) *sacloud.AutoBackup {
	return sacloud.CreateNewAutoBackup(name, diskID)
}

func (api *AutoBackupAPI) Create(value *sacloud.AutoBackup) (*sacloud.AutoBackup, error) {
	return api.request(func(res *autoBackupResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *AutoBackupAPI) Read(id string) (*sacloud.AutoBackup, error) {
	return api.request(func(res *autoBackupResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *AutoBackupAPI) Update(id string, value *sacloud.AutoBackup) (*sacloud.AutoBackup, error) {
	return api.request(func(res *autoBackupResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *AutoBackupAPI) Delete(id string) (*sacloud.AutoBackup, error) {
	return api.request(func(res *autoBackupResponse) error {
		return api.delete(id, nil, res)
	})
}
