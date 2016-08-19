package api

import (
	"encoding/json"
	"fmt"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

//HACK: さくらのAPI側仕様: Applianceの内容によってJSONフォーマットが異なるため
//      ロードバランサ/VPCルータそれぞれでリクエスト/レスポンスデータ型を定義する。

type SearchDatabaseResponse struct {
	Total     int                `json:",omitempty"`
	From      int                `json:",omitempty"`
	Count     int                `json:",omitempty"`
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
	Success           interface{} `json:",omitempty"` //HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないためinterface{}
}

type DatabaseAPI struct {
	*baseAPI
}

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

func (api *DatabaseAPI) New(values *sacloud.CreateDatabaseValue) *sacloud.Database {
	return sacloud.CreateNewPostgreSQLDatabase(values)
}

func (api *DatabaseAPI) Create(value *sacloud.Database) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *DatabaseAPI) Read(id string) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *DatabaseAPI) Update(id string, value *sacloud.Database) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *DatabaseAPI) UpdateSetting(id string, value *sacloud.Database) (*sacloud.Database, error) {
	req := &sacloud.Database{
		Settings: value.Settings,
	}
	return api.request(func(res *databaseResponse) error {
		return api.update(id, api.createRequest(req), res)
	})
}

func (api *DatabaseAPI) Delete(id string) (*sacloud.Database, error) {
	return api.request(func(res *databaseResponse) error {
		return api.delete(id, nil, res)
	})
}

func (api *DatabaseAPI) Config(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/config", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

func (api *DatabaseAPI) IsUp(id string) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsUp(), nil
}

func (api *DatabaseAPI) IsDown(id string) (bool, error) {
	lb, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return lb.Instance.IsDown(), nil
}

// Boot power on
func (api *DatabaseAPI) Boot(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown power off
func (api *DatabaseAPI) Shutdown(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop force shutdown
func (api *DatabaseAPI) Stop(id string) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%s/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

func (api *DatabaseAPI) RebootForce(id string) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

func (api *DatabaseAPI) ResetForce(id string, recycleProcess bool) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%s/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"RecycleProcess": recycleProcess})
}

func (api *DatabaseAPI) SleepUntilUp(id string, timeout time.Duration) error {
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

func (api *DatabaseAPI) SleepUntilDown(id string, timeout time.Duration) error {
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

// SleepWhileCopying wait until became to available
func (api *DatabaseAPI) SleepWhileCopying(id string, timeout time.Duration, maxRetryCount int) error {
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
