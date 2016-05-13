package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

type ServerAPI struct {
	*baseAPI
}

func NewServerAPI(client *Client) *ServerAPI {
	return &ServerAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "server"
			},
		},
	}
}

func (api *ServerAPI) WithPlan(planID string) *ServerAPI {
	return api.FilterBy("ServerPlan.ID", planID)
}

func (api *ServerAPI) WithStatus(status string) *ServerAPI {
	return api.FilterBy("Instance.Status", status)
}
func (api *ServerAPI) WithStatusUp(status string) *ServerAPI {
	return api.WithStatus("up")
}
func (api *ServerAPI) WithStatusDown(status string) *ServerAPI {
	return api.WithStatus("down")
}

func (api *ServerAPI) WithISOImage(imageID string) *ServerAPI {
	return api.FilterBy("Instance.CDROM.ID", imageID)
}

func (api *ServerAPI) SortByCPU(reverse bool) *ServerAPI {
	api.sortBy("ServerPlan.CPU", reverse)
	return api
}

func (api *ServerAPI) SortByMemory(reverse bool) *ServerAPI {
	api.sortBy("ServerPlan.MemoryMB", reverse)
	return api
}

// CreateWithAdditionalIP create server
func (api *ServerAPI) CreateWithAdditionalIP(spec *sacloud.Server, addIPAddress string) (*sacloud.Server, error) {
	//TODO 高レベルAPIへ移動

	server, err := api.Create(spec)
	if err != nil {
		return nil, err
	}

	if addIPAddress != "" && len(server.Interfaces) > 1 {
		if err := api.updateIPAddress(server, addIPAddress); err != nil {
			return nil, err
		}
	}

	return server, nil
}

func (api *ServerAPI) updateIPAddress(server *sacloud.Server, ip string) error {
	//TODO 高レベルAPIへ移動
	var (
		method = "PUT"
		uri    = fmt.Sprintf("interface/%s", server.Interfaces[1].ID)
		body   = sacloud.Request{}
	)
	body.Interface = &sacloud.Interface{UserIPAddress: ip}

	_, err := api.client.newRequest(method, uri, body)
	if err != nil {
		return err
	}

	return nil

}

// State get server state
func (api *ServerAPI) State(id string) (string, error) {
	server, err := api.Read(id)
	if err != nil {
		return "", err
	}
	return server.Availability, nil
}

func (api *ServerAPI) IsUp(id string) (bool, error) {
	server, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return server.Instance.IsUp(), nil
}

func (api *ServerAPI) IsDown(id string) (bool, error) {
	server, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return server.Instance.IsDown(), nil
}

// Boot power on
func (api *ServerAPI) Boot(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown power off
func (api *ServerAPI) Shutdown(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop force shutdown
func (api *ServerAPI) Stop(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

func (api *ServerAPI) RebootForce(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

func (api *ServerAPI) SleepUntilUp(serverID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		up, err := api.IsUp(serverID)
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

func (api *ServerAPI) SleepUntilDown(serverID string, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		down, err := api.IsDown(serverID)
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

func (api *ServerAPI) ChangePlan(serverID string, planID string) (*sacloud.Server, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/to/plan/%s", api.getResourceURL(), serverID, planID)
	)

	return api.request(func(res *sacloud.Response) error {
		return api.baseAPI.request(method, uri, nil, res)
	})
}

func (api *ServerAPI) FindDisk(serverID string) ([]sacloud.Disk, error) {
	server, err := api.Read(serverID)
	if err != nil {
		return nil, err
	}
	return server.Disks, nil
}

func (api *ServerAPI) InsertCDROM(serverID string, cdromID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/cdrom", api.getResourceURL(), serverID)
	)

	req := &sacloud.Request{
		SakuraCloudResources: sacloud.SakuraCloudResources{
			CDROM: &sacloud.CDROM{Resource: &sacloud.Resource{ID: cdromID}},
		},
	}

	return api.modify(method, uri, req)
}

func (api *ServerAPI) EjectCDROM(serverID string, cdromID string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/cdrom", api.getResourceURL(), serverID)
	)

	req := &sacloud.Request{
		SakuraCloudResources: sacloud.SakuraCloudResources{
			CDROM: &sacloud.CDROM{Resource: &sacloud.Resource{ID: cdromID}},
		},
	}

	return api.modify(method, uri, req)
}

func (api *ServerAPI) NewKeyboardRequest() *sacloud.KeyboardRequest {
	return &sacloud.KeyboardRequest{}
}

func (api *ServerAPI) SendKey(serverID string, body *sacloud.KeyboardRequest) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/keyboard", api.getResourceURL(), serverID)
	)

	return api.modify(method, uri, body)
}

func (api *ServerAPI) NewMouseRequest() *sacloud.MouseRequest {
	return &sacloud.MouseRequest{
		Buttons: &sacloud.MouseRequestButtons{},
	}
}

func (api *ServerAPI) SendMouse(serverID string, mouseIndex string, body *sacloud.MouseRequest) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/mouse/%s", api.getResourceURL(), serverID, mouseIndex)
	)

	return api.modify(method, uri, body)
}

func (api *ServerAPI) NewVNCSnapshotRequest() *sacloud.VNCSnapshotRequest {
	return &sacloud.VNCSnapshotRequest{}
}

func (api *ServerAPI) GetVNCProxy(serverID string) (*sacloud.VNCProxyResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s/vnc/proxy", api.getResourceURL(), serverID)
		res    = &sacloud.VNCProxyResponse{}
	)
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ServerAPI) GetVNCSize(serverID string) (*sacloud.VNCSizeResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s/vnc/size", api.getResourceURL(), serverID)
		res    = &sacloud.VNCSizeResponse{}
	)
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ServerAPI) GetVNCSnapshot(serverID string, body *sacloud.VNCSnapshotRequest) (*sacloud.VNCSnapshotResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%s/vnc/snapshot", api.getResourceURL(), serverID)
		res    = &sacloud.VNCSnapshotResponse{}
	)
	err := api.baseAPI.request(method, uri, body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ServerAPI) Monitor(id string, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.monitor(id, body)
}
