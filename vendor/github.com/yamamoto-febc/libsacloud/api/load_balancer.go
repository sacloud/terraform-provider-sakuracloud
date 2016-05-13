package api

import (
	"encoding/json"
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

//HACK: さくらのAPI側仕様: Applianceの内容によってJSONフォーマットが異なるため
//      ロードバランサ/VPCルータそれぞれでリクエスト/レスポンスデータ型を定義する。

type SearchLoadBalancerResponse struct {
	Total         int                    `json:",omitempty"`
	From          int                    `json:",omitempty"`
	Count         int                    `json:",omitempty"`
	LoadBalancers []sacloud.LoadBalancer `json:"Appliances,omitempty"`
}

type loadBalancerRequest struct {
	LoadBalander *sacloud.LoadBalancer  `json:"Appliance,omitempty"`
	From         int                    `json:",omitempty"`
	Count        int                    `json:",omitempty"`
	Sort         []string               `json:",omitempty"`
	Filter       map[string]interface{} `json:",omitempty"`
	Exclude      []string               `json:",omitempty"`
	Include      []string               `json:",omitempty"`
}

type loadBalanderResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.LoadBalancer `json:"Appliance,omitempty"`
	Success               interface{} `json:",omitempty"` //HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないためinterface{}
}

type LoadBalancerAPI struct {
	*baseAPI
}

func NewLoadBalancerAPI(client *Client) *LoadBalancerAPI {
	return &LoadBalancerAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "appliance"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Class", "loadbalancer")
				return res
			},
		},
	}
}

func (api *LoadBalancerAPI) Find() (*SearchLoadBalancerResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchLoadBalancerResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *LoadBalancerAPI) request(f func(*loadBalanderResponse) error) (*sacloud.LoadBalancer, error) {
	res := &loadBalanderResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.LoadBalancer, nil
}

func (api *LoadBalancerAPI) createRequest(value *sacloud.LoadBalancer) *loadBalanderResponse {
	return &loadBalanderResponse{LoadBalancer: value}
}

func (api *LoadBalancerAPI) New() *sacloud.LoadBalancer {
	return sacloud.CreateNewLoadBalancer()
}

func (api *LoadBalancerAPI) Create(value *sacloud.LoadBalancer) (*sacloud.LoadBalancer, error) {
	return api.request(func(res *loadBalanderResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *LoadBalancerAPI) Read(id string) (*sacloud.LoadBalancer, error) {
	return api.request(func(res *loadBalanderResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *LoadBalancerAPI) Update(id string, value *sacloud.LoadBalancer) (*sacloud.LoadBalancer, error) {
	return api.request(func(res *loadBalanderResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *LoadBalancerAPI) Delete(id string) (*sacloud.LoadBalancer, error) {
	return api.request(func(res *loadBalanderResponse) error {
		return api.delete(id, nil, res)
	})
}

// SleepWhileCopying wait until became to available
func (api *LoadBalancerAPI) SleepWhileCopying(loadBalancerID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {
		loadBalancer, err := api.Read(loadBalancerID)
		if err != nil {
			return err
		}

		if loadBalancer.IsAvailable() {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepWhileCopying")
		}
	}
}
