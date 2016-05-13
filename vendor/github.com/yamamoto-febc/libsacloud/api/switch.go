package api

import "fmt"

type SwitchAPI struct {
	*baseAPI
}

func NewSwitchAPI(client *Client) *SwitchAPI {
	return &SwitchAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "switch"
			},
		},
	}
}

func (api *SwitchAPI) DisconnectFromBridge(switchID string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/to/bridge", api.getResourceURL(), switchID)
	)
	return api.modify(method, uri, nil)
}

func (api *SwitchAPI) ConnectToBridge(switchID string, bridgeID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/bridge/%s", api.getResourceURL(), switchID, bridgeID)
	)
	return api.modify(method, uri, nil)
}
