package api

import (
	"encoding/json"
	//	"strings"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

//HACK: さくらのAPI側仕様: CommonServiceItemsの内容によってJSONフォーマットが異なるため
//      DNS/GSLB/シンプル監視それぞれでリクエスト/レスポンスデータ型を定義する。

type SearchSimpleMonitorResponse struct {
	Total          int                     `json:",omitempty"`
	From           int                     `json:",omitempty"`
	Count          int                     `json:",omitempty"`
	SimpleMonitors []sacloud.SimpleMonitor `json:"CommonServiceItems,omitempty"`
}

type simpleMonitorRequest struct {
	SimpleMonitor *sacloud.SimpleMonitor `json:"CommonServiceItem,omitempty"`
	From          int                    `json:",omitempty"`
	Count         int                    `json:",omitempty"`
	Sort          []string               `json:",omitempty"`
	Filter        map[string]interface{} `json:",omitempty"`
	Exclude       []string               `json:",omitempty"`
	Include       []string               `json:",omitempty"`
}

type simpleMonitorResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.SimpleMonitor `json:"CommonServiceItem,omitempty"`
}

type SimpleMonitorAPI struct {
	*baseAPI
}

func NewSimpleMonitorAPI(client *Client) *SimpleMonitorAPI {
	return &SimpleMonitorAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "commonserviceitem"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Provider.Class", "simplemon")
				return res
			},
		},
	}
}

func (api *SimpleMonitorAPI) Find() (*SearchSimpleMonitorResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchSimpleMonitorResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *SimpleMonitorAPI) request(f func(*simpleMonitorResponse) error) (*sacloud.SimpleMonitor, error) {
	res := &simpleMonitorResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.SimpleMonitor, nil
}

func (api *SimpleMonitorAPI) createRequest(value *sacloud.SimpleMonitor) *simpleMonitorResponse {
	return &simpleMonitorResponse{SimpleMonitor: value}
}

func (api *SimpleMonitorAPI) New(target string) *sacloud.SimpleMonitor {
	return sacloud.CreateNewSimpleMonitor(target)
}

func (api *SimpleMonitorAPI) Create(value *sacloud.SimpleMonitor) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *SimpleMonitorAPI) Read(id string) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *SimpleMonitorAPI) Update(id string, value *sacloud.SimpleMonitor) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *SimpleMonitorAPI) Delete(id string) (*sacloud.SimpleMonitor, error) {
	return api.request(func(res *simpleMonitorResponse) error {
		return api.delete(id, nil, res)
	})
}
