package api

import (
	"encoding/json"
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

//HACK: さくらのAPI側仕様: Applianceの内容によってJSONフォーマットが異なるため
//      ロードバランサ/VPCルータそれぞれでリクエスト/レスポンスデータ型を定義する。

type SearchVPCRouterResponse struct {
	Total      int                 `json:",omitempty"`
	From       int                 `json:",omitempty"`
	Count      int                 `json:",omitempty"`
	VPCRouters []sacloud.VPCRouter `json:"Appliances,omitempty"`
}

type vpcRouterRequest struct {
	VPCRouter *sacloud.VPCRouter     `json:"Appliance,omitempty"`
	From      int                    `json:",omitempty"`
	Count     int                    `json:",omitempty"`
	Sort      []string               `json:",omitempty"`
	Filter    map[string]interface{} `json:",omitempty"`
	Exclude   []string               `json:",omitempty"`
	Include   []string               `json:",omitempty"`
}

type vpcRouterResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.VPCRouter `json:"Appliance,omitempty"`
	Success            interface{} `json:",omitempty"` //HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないためinterface{}
}

type VPCRouterAPI struct {
	*baseAPI
}

func NewVPCRouterAPI(client *Client) *VPCRouterAPI {
	return &VPCRouterAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "appliance"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Class", "vpcrouter")
				return res
			},
		},
	}
}

func (api *VPCRouterAPI) Find() (*SearchVPCRouterResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchVPCRouterResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *VPCRouterAPI) request(f func(*vpcRouterResponse) error) (*sacloud.VPCRouter, error) {
	res := &vpcRouterResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.VPCRouter, nil
}

func (api *VPCRouterAPI) createRequest(value *sacloud.VPCRouter) *vpcRouterResponse {
	return &vpcRouterResponse{VPCRouter: value}
}

func (api *VPCRouterAPI) New() *sacloud.VPCRouter {
	return sacloud.CreateNewVPCRouter()
}

func (api *VPCRouterAPI) Create(value *sacloud.VPCRouter) (*sacloud.VPCRouter, error) {
	return api.request(func(res *vpcRouterResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *VPCRouterAPI) Read(id string) (*sacloud.VPCRouter, error) {
	return api.request(func(res *vpcRouterResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *VPCRouterAPI) Update(id string, value *sacloud.VPCRouter) (*sacloud.VPCRouter, error) {
	return api.request(func(res *vpcRouterResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *VPCRouterAPI) Delete(id string) (*sacloud.VPCRouter, error) {
	return api.request(func(res *vpcRouterResponse) error {
		return api.delete(id, nil, res)
	})
}

func (api *VPCRouterAPI) Config(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/config", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

func (api *VPCRouterAPI) ConnectToSwitch(id string, switchID string, nicIndex int) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/interface/%d/to/switch/%s", api.getResourceURL(), id, nicIndex, switchID)
	)
	return api.modify(method, uri, nil)
}

func (api *VPCRouterAPI) DisconnectFromSwitch(id string, switchID string, nicIndex int) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/interface/%d/to/switch/%s", api.getResourceURL(), id, nicIndex, switchID)
	)
	return api.modify(method, uri, nil)
}

func (api *VPCRouterAPI) IsUp(id string) (bool, error) {
	router, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return router.Instance.IsUp(), nil
}

func (api *VPCRouterAPI) IsDown(id string) (bool, error) {
	router, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return router.Instance.IsDown(), nil
}

// Boot power on
func (api *VPCRouterAPI) Boot(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown power off
func (api *VPCRouterAPI) Shutdown(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop force shutdown
func (api *VPCRouterAPI) Stop(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

func (api *VPCRouterAPI) RebootForce(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

func (api *VPCRouterAPI) SleepUntilUp(routerID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		up, err := api.IsUp(routerID)
		if err != nil {
			return err
		}

		if up {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: WaitforAvailable")
		}
	}
}

func (api *VPCRouterAPI) SleepUntilDown(routerID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		down, err := api.IsDown(routerID)
		if err != nil {
			return err
		}

		if down {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: WaitforAvailable")
		}
	}
}

// SleepWhileCopying wait until became to available
func (api *VPCRouterAPI) SleepWhileCopying(vpcRouterID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {
		router, err := api.Read(vpcRouterID)
		if err != nil {
			return err
		}

		if router.IsAvailable() {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepWhileCopying")
		}
	}
}

//func (api *VPCRouterAPI) AddStandardInterface(routerID string, switchID string, ipaddress string, maskLen int) error {
//	return api.addInterface(routerID, switchID, &sacloud.VPCRouterInterface{
//		IPAddress:        []string{ipaddress},
//		NetworkMaskLen:   maskLen,
//		VirtualIPAddress: "",
//	})
//}
//
//func (api *VPCRouterAPI) AddPremiumInterface(routerID string, switchID string, ipaddresses []string, maskLen int, virtualIP string) error {
//	return api.addInterface(routerID, switchID, &sacloud.VPCRouterInterface{
//		IPAddress:        ipaddresses,
//		NetworkMaskLen:   maskLen,
//		VirtualIPAddress: virtualIP,
//	})
//}
//
//func (api *VPCRouterAPI) addInterface(routerID string, switchID string, routerNIC *sacloud.VPCRouterInterface) error {
//	router, err := api.Read(routerID)
//	if err != nil {
//		return err
//	}
//
//	req := &sacloud.VPCRouter{Settings: &sacloud.VPCRouterSettings{}}
//
//	if router.Settings == nil {
//		req.Settings = &sacloud.VPCRouterSettings{
//			Router: &sacloud.VPCRouterSetting{
//				Interfaces: []*sacloud.VPCRouterInterface{nil},
//			},
//		}
//	} else {
//		req.Settings.Router = router.Settings.Router
//	}
//
//	//connect to switch
//	//_, err = api.ConnectToSwitch(routerID, switchID, len(req.Settings.Router.Interfaces)+1)
//	//if err != nil {
//	//	return err
//	//}
//
//	req.Settings.Router.Interfaces = append(req.Settings.Router.Interfaces, routerNIC)
//
//	_, err = api.Update(routerID, req)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
