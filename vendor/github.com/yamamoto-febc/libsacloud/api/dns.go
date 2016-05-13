package api

import (
	"encoding/json"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"strings"
)

//HACK: さくらのAPI側仕様: CommonServiceItemsの内容によってJSONフォーマットが異なるため
//      DNS/GSLB/シンプル監視それぞれでリクエスト/レスポンスデータ型を定義する。

type SearchDNSResponse struct {
	Total                 int           `json:",omitempty"`
	From                  int           `json:",omitempty"`
	Count                 int           `json:",omitempty"`
	CommonServiceDNSItems []sacloud.DNS `json:"CommonServiceItems,omitempty"`
}
type dnsRequest struct {
	CommonServiceDNSItem *sacloud.DNS           `json:"CommonServiceItem,omitempty"`
	From                 int                    `json:",omitempty"`
	Count                int                    `json:",omitempty"`
	Sort                 []string               `json:",omitempty"`
	Filter               map[string]interface{} `json:",omitempty"`
	Exclude              []string               `json:",omitempty"`
	Include              []string               `json:",omitempty"`
}
type dnsResponse struct {
	*sacloud.ResultFlagValue
	*sacloud.DNS `json:"CommonServiceItem,omitempty"`
}

// DNSAPI API Client for SAKURA CLOUD DNS
type DNSAPI struct {
	*baseAPI
}

func NewDNSAPI(client *Client) *DNSAPI {
	return &DNSAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "commonserviceitem"
			},
			FuncBaseSearchCondition: func() *sacloud.Request {
				res := &sacloud.Request{}
				res.AddFilter("Provider.Class", "dns")
				return res
			},
		},
	}
}

func (api *DNSAPI) Find(condition *sacloud.Request) (*SearchDNSResponse, error) {

	data, err := api.client.newRequest("GET", api.getResourceURL(), api.getSearchState())
	if err != nil {
		return nil, err
	}
	var res SearchDNSResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *DNSAPI) request(f func(*dnsResponse) error) (*sacloud.DNS, error) {
	res := &dnsResponse{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.DNS, nil
}

func (api *DNSAPI) createRequest(value *sacloud.DNS) *dnsRequest {
	req := &dnsRequest{}
	req.CommonServiceDNSItem = value
	return req
}
func (api *DNSAPI) Create(value *sacloud.DNS) (*sacloud.DNS, error) {
	return api.request(func(res *dnsResponse) error {
		return api.create(api.createRequest(value), res)
	})
}

func (api *DNSAPI) New(zoneName string) *sacloud.DNS {
	return sacloud.CreateNewDNS(zoneName)
}

func (api *DNSAPI) Read(id string) (*sacloud.DNS, error) {
	return api.request(func(res *dnsResponse) error {
		return api.read(id, nil, res)
	})
}

func (api *DNSAPI) Update(id string, value *sacloud.DNS) (*sacloud.DNS, error) {
	return api.request(func(res *dnsResponse) error {
		return api.update(id, api.createRequest(value), res)
	})
}

func (api *DNSAPI) Delete(id string) (*sacloud.DNS, error) {
	return api.request(func(res *dnsResponse) error {
		return api.delete(id, nil, res)
	})
}

// SetupDNSRecord get dns zone commonserviceitem id
func (api *DNSAPI) SetupDNSRecord(zoneName string, hostName string, ip string) ([]string, error) {

	dnsItem, err := api.findOrCreateBy(zoneName)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(hostName, zoneName) {
		hostName = strings.Replace(hostName, zoneName, "", -1)
	}

	dnsItem.Settings.DNS.AddDNSRecordSet(hostName, ip)

	res, err := api.updateDNSRecord(dnsItem)
	if err != nil {
		return nil, err
	}

	if dnsItem.ID == "" {
		return res.Status.NS, nil
	}

	return nil, nil

}

// DeleteDNSRecord delete dns record
func (api *DNSAPI) DeleteDNSRecord(zoneName string, hostName string, ip string) error {
	dnsItem, err := api.findOrCreateBy(zoneName)
	if err != nil {
		return err
	}
	dnsItem.Settings.DNS.DeleteDNSRecordSet(hostName, ip)

	if dnsItem.HasDNSRecord() {
		_, err = api.updateDNSRecord(dnsItem)
		if err != nil {
			return err
		}

	} else {
		_, err = api.Delete(dnsItem.ID)
		if err != nil {
			return err
		}

	}
	return nil
}

func (api *DNSAPI) findOrCreateBy(zoneName string) (*sacloud.DNS, error) {

	req := &sacloud.Request{}
	req.AddFilter("Name", zoneName)
	res, err := api.Find(req)
	if err != nil {
		return nil, err
	}

	//すでに登録されている場合
	var dnsItem *sacloud.DNS
	if res.Count > 0 {
		dnsItem = &res.CommonServiceDNSItems[0]
	} else {
		dnsItem = sacloud.CreateNewDNS(zoneName)
	}

	return dnsItem, nil
}

func (api *DNSAPI) updateDNSRecord(dnsItem *sacloud.DNS) (*sacloud.DNS, error) {

	var item *sacloud.DNS
	var err error

	if dnsItem.ID == "" {
		item, err = api.Create(dnsItem)
	} else {
		item, err = api.Update(dnsItem.ID, dnsItem)
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}
