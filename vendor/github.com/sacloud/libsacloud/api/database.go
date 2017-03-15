package api

import (
	"encoding/json"
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
)

//HACK: さくらのAPI側仕様: Applianceの内容によってJSONフォーマットが異なるため
//      ロードバランサ/VPCルータそれぞれでリクエスト/レスポンスデータ型を定義する。

// SearchDatabaseResponse データベース検索レスポンス
type SearchDatabaseResponse struct {
	// Total 総件数
	Total int `json:",omitempty"`
	// From ページング開始位置
	From int `json:",omitempty"`
	// Count 件数
	Count int `json:",omitempty"`
	// Databases データベースリスト
	Databases []sacloud.Database `json:"Appliances,omitempty"`
}

type databaseRequest struct {
	Database *sacloud.Database      `json:"Appliance,omitempty"`
	From     int                    `json:",omitempty"`
	Count    int                    `json:",omitempty"`
	Sort     []string               `json:",omitempty"`
	Filter   map[string]interface{} `json:",omitempty"`
	Exclude  []string               `json:",omitempty"`
	Include  []string               `json:",omitempty"`
}

type databaseResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.Database `json:"Appliance,omitempty"`
	// Success
	Success interface{} `json:",omitempty"` //HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないためinterface{}で受ける
}

type databaseSettingsResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.Database `json:"Appliance,omitempty"`
}

// DatabaseAPI データベースAPI
type DatabaseAPI struct {
	*baseAPI
}

// NewDatabaseAPI データベースAPI作成
func NewDatabaseAPI(client *Client) *DatabaseAPI {
	return &DatabaseAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "appliance"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Class", "database")
				return res
			},
		},
	}
}

// Find 検索
func (api *DatabaseAPI) Find() (*SearchDatabaseResponse, error) {
	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchDatabaseResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *DatabaseAPI) request(f func(*databaseResponse) error) (*sacloud.Database, error) {
	res := &databaseResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.Database, nil
}

func (api *DatabaseAPI) createRequest(value *sacloud.Database) *databaseResponse {
	return &databaseResponse{Database: value}
}

// New 新規作成用パラメーター作成
func (api *DatabaseAPI) New(values *sacloud.CreateDatabaseValue) *sacloud.Database {
	return sacloud.CreateNewDatabase(values)
}

// Create 新規作成
func (api *DatabaseAPI) Create(value *sacloud.Database) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

// Read 読み取り
func (api *DatabaseAPI) Read(id int64) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.read(id, nil, res)
	})
}

// Status DBの設定/起動状態の取得
func (api *DatabaseAPI) Status(id int64) (*sacloud.DatabaseSettingsResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/status", api.getResourceURL(), id)
	)

	res := &databaseResponse{}
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res.SettingsResponse, nil
}

// Update 更新
func (api *DatabaseAPI) Update(id int64, value *sacloud.Database) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

// UpdateSetting 設定更新
func (api *DatabaseAPI) UpdateSetting(id int64, value *sacloud.Database) (*sacloud.Database, error) {
	req := &sacloud.Database{
		// Settings
		Settings: value.Settings,
	}
	return api.request(func(res *databaseResponse) error {
		return api.update(id, api.createRequest(req), res)
	})
}

// Delete 削除
func (api *DatabaseAPI) Delete(id int64) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.delete(id, nil, res)
	})
}

// Config 設定変更の反映
func (api *DatabaseAPI) Config(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/config", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// IsUp 起動しているか判定
func (api *DatabaseAPI) IsUp(id int64) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsUp(), nil
}

// IsDown ダウンしているか判定
func (api *DatabaseAPI) IsDown(id int64) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsDown(), nil
}

// IsDatabaseRunning データベースプロセスが起動しているか判定
func (api *DatabaseAPI) IsDatabaseRunning(id int64) (bool, error) {
	db, err := api.Status(id)
	if err != nil {
		return false, err
	}
	return db.IsUp(), nil

}

// Boot 起動
func (api *DatabaseAPI) Boot(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown シャットダウン(graceful)
func (api *DatabaseAPI) Shutdown(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop シャットダウン(force)
func (api *DatabaseAPI) Stop(id int64) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

// RebootForce 再起動
func (api *DatabaseAPI) RebootForce(id int64) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// ResetForce リセット
func (api *DatabaseAPI) ResetForce(id int64, recycleProcess bool) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"RecycleProcess": recycleProcess})
}

// SleepUntilUp 起動するまで待機
func (api *DatabaseAPI) SleepUntilUp(id int64, timeout time.Duration) error {
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

// SleepUntilDatabaseRunning 起動するまで待機
func (api *DatabaseAPI) SleepUntilDatabaseRunning(id int64, timeout time.Duration, maxRetryCount int) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	errCount := 0

	for {
		isUp, err := api.IsUp(id)
		if err != nil {
			errCount++
			if errCount > maxRetryCount {
				return err
			}
		}

		if isUp {
			return nil
		}
		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: SleepUntilDatabaseRunning")
		}
	}
}

// SleepUntilDown ダウンするまで待機
func (api *DatabaseAPI) SleepUntilDown(id int64, timeout time.Duration) error {
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
func (api *DatabaseAPI) SleepWhileCopying(id int64, timeout time.Duration, maxRetryCount int) error {
	current := 0 * time.Second
	interval := 5 * time.Second
	errCount := 0

	for {
		database, err := api.Read(id)
		if err != nil {
			errCount++
			if errCount > maxRetryCount {
				return err
			}
		}

		if database != nil && database.IsFailed() {
			return fmt.Errorf("Create database Failed")
		}

		if database != nil && database.IsAvailable() {
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
func (api *DatabaseAPI) AsyncSleepWhileCopying(id int64, timeout time.Duration, maxRetryCount int) (chan (*sacloud.Database), chan (*sacloud.Database), chan (error)) {
	complete := make(chan *sacloud.Database)
	progress := make(chan *sacloud.Database)
	err := make(chan error)
	errCount := 0

	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				db, e := api.Read(id)
				if e != nil {
					errCount++
					if errCount > maxRetryCount {
						err <- e
						return
					}
				}

				progress <- db

				if db.IsAvailable() {
					complete <- db
					return
				}
				if db.IsFailed() {
					err <- fmt.Errorf("Failed: Create Database is failed: %#v", db)
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

// MonitorCPU CPUアクティビティーモニター取得
func (api *DatabaseAPI) MonitorCPU(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "cpu", 0, body)
}

// MonitorDatabase データーベース固有項目アクティビティモニター取得
func (api *DatabaseAPI) MonitorDatabase(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "database", 0, body)
}

// MonitorInterface NICアクティビティーモニター取得
func (api *DatabaseAPI) MonitorInterface(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "interface", 0, body)
}

// MonitorSystemDisk システムディスクアクティビティーモニター取得
func (api *DatabaseAPI) MonitorSystemDisk(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "disk", 1, body)
}

// MonitorBackupDisk バックアップディスクアクティビティーモニター取得
func (api *DatabaseAPI) MonitorBackupDisk(id int64, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.applianceMonitorBy(id, "disk", 2, body)
}
