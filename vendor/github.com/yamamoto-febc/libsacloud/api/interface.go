package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

type InterfaceAPI struct {
	*baseAPI
}

func NewInterfaceAPI(client *Client) *InterfaceAPI {
	return &InterfaceAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "interface"
			},
		},
	}
}

func (api *InterfaceAPI) CreateAndConnectToServer(serverID string) (*sacloud.Interface, error) {
	iface := api.New()
	iface.Server = &sacloud.Server{
		Resource: &sacloud.Resource{ID: serverID},
	}
	return api.Create(iface)
}

func (api *InterfaceAPI) ConnectToSwitch(interfaceID string, switchID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/switch/%s", api.getResourceURL(), interfaceID, switchID)
	)
	return api.modify(method, uri, nil)
}

func (api *InterfaceAPI) ConnectToSharedSegment(interfaceID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/switch/shared", api.getResourceURL(), interfaceID)
	)
	return api.modify(method, uri, nil)
}

func (api *InterfaceAPI) DisconnectFromSwitch(interfaceID string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/to/switch", api.getResourceURL(), interfaceID)
	)
	return api.modify(method, uri, nil)
}

func (api *InterfaceAPI) Monitor(id string, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.monitor(id, body)
}
