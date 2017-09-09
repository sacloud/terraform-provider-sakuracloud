package api

import (
	"encoding/json"
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
)

// SearchNFSResponse NFS検索レスポンス
type SearchNFSResponse struct {
	// Total 総件数
	Total int `json:",omitempty"`
	// From ページング開始位置
	From int `json:",omitempty"`
	// Count 件数
	Count int `json:",omitempty"`
	// NFSs NFS リスト
	NFS []sacloud.NFS `json:"Appliances,omitempty"`
}

type nfsRequest struct {
	NFS     *sacloud.NFS           `json:"Appliance,omitempty"`
	From    int                    `json:",omitempty"`
	Count   int                    `json:",omitempty"`
	Sort    []string               `json:",omitempty"`
	Filter  map[string]interface{} `json:",omitempty"`
	Exclude []string               `json:",omitempty"`
	Include []string               `json:",omitempty"`
}

type nfsResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.NFS `json:"Appliance,omitempty"`
	Success      interface{} `json:",omitempty"` //HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないためinterface{}
}

// NFSAPI NFSAPI
type NFSAPI struct {
	*baseAPI
}

// NewNFSAPI NFSAPI作成
func NewNFSAPI(client *Client) *NFSAPI {
	return &NFSAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "appliance"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Class", "nfs")
				return res
			},
		},
	}
}

// Find 検索
func (api *NFSAPI) Find() (*SearchNFSResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchNFSResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *NFSAPI) request(f func(*nfsResponse) error) (*sacloud.NFS, error) {
	res := &nfsResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.NFS, nil
}

func (api *NFSAPI) createRequest(value *sacloud.NFS) *nfsResponse {
	return &nfsResponse{NFS: value}
}

//func (api *NFSAPI) New() *sacloud.NFS {
//	return sacloud.CreateNewNFS()
//}

// Create 新規作成
func (api *NFSAPI) Create(value *sacloud.NFS) (*sacloud.NFS, error) {
	return api.request(func(res *nfsResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

// Read 読み取り
func (api *NFSAPI) Read(id int64) (*sacloud.NFS, error) {
	return api.request(func(res *nfsResponse) error {
		return api.read(id, nil, res)
	})
}

// Update 更新
func (api *NFSAPI) Update(id int64, value *sacloud.NFS) (*sacloud.NFS, error) {
	return api.request(func(res *nfsResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

// Delete 削除
func (api *NFSAPI) Delete(id int64) (*sacloud.NFS, error) {
	return api.request(func(res *nfsResponse) error {
		return api.delete(id, nil, res)
	})
}

// Config 設定変更の反映
func (api *NFSAPI) Config(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/config", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// IsUp 起動しているか判定
func (api *NFSAPI) IsUp(id int64) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsUp(), nil
}

// IsDown ダウンしているか判定
func (api *NFSAPI) IsDown(id int64) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsDown(), nil
}

// Boot 起動
func (api *NFSAPI) Boot(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown シャットダウン(graceful)
func (api *NFSAPI) Shutdown(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop シャットダウン(force)
func (api *NFSAPI) Stop(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

// RebootForce 再起動
func (api *NFSAPI) RebootForce(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// ResetForce リセット
func (api *NFSAPI) ResetForce(id int64, recycleProcess bool) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"RecycleProcess": recycleProcess})
}

// SleepUntilUp 起動するまで待機
func (api *NFSAPI) SleepUntilUp(id int64, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		up, err := api.IsUp(id)
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

// SleepUntilDown ダウンするまで待機
func (api *NFSAPI) SleepUntilDown(id int64, timeout time.Duration) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	for {

		down, err := api.IsDown(id)
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

// SleepWhileCopying コピー終了まで待機
func (api *NFSAPI) SleepWhileCopying(id int64, timeout time.Duration, maxRetryCount int) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	errCount := 0

	for {
		nfs, err := api.Read(id)
		if err != nil {
			errCount++
			if errCount > maxRetryCount {
				return err
			}
		}

		if nfs != nil && nfs.IsAvailable() {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepWhileCopying")
		}
	}
}

// AsyncSleepWhileCopying コピー終了まで待機(非同期)
func (api *NFSAPI) AsyncSleepWhileCopying(id int64, timeout time.Duration, maxRetryCount int) (chan (*sacloud.NFS), chan (*sacloud.NFS), chan (error)) {
	complete := make(chan *sacloud.NFS)
	progress := make(chan *sacloud.NFS)
	err := make(chan error)
	errCount := 0

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				lb, e := api.Read(id)
				if e != nil {
					errCount++
					if errCount > maxRetryCount {
						err <- e
						return
					}
				}

				progress <- lb

				if lb.IsAvailable() {
					complete <- lb
					return
				}
				if lb.IsFailed() {
					err <- fmt.Errorf("Failed: Create NFS is failed: %#v", lb)
					return
				}

			case <-time.After(timeout):
				err <- fmt.Errorf("Timeout: AsyncSleepWhileCopying[ID:%d]", id)
				return
			}
		}
	}()
	return complete, progress, err
}

// MonitorNFS NFS固有項目アクティビティモニター取得
func (api *NFSAPI) MonitorNFS(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "nfs", 0, body)
}

// MonitorInterface NICアクティビティーモニター取得
func (api *NFSAPI) MonitorInterface(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "interface", 0, body)
}
