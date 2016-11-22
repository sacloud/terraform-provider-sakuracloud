package api

import (
	"encoding/json"
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
)

// WebAccelAPI ウェブアクセラレータAPI
type WebAccelAPI struct {
	*baseAPI
}

// NewWebAccelAPI ウェブアクセラレータAPI
func NewWebAccelAPI(client *Client) *WebAccelAPI {
	return &WebAccelAPI{
		&baseAPI{
			client:        client,
			apiRootSuffix: sakuraWebAccelAPIRootSuffix,
			FuncGetResourceURL: func() string {
				return ""
			},
		},
	}
}

// WebAccelDeleteCacheResponse ウェブアクセラレータ キャッシュ削除レスポンス
type WebAccelDeleteCacheResponse struct {
	*sacloud.ResultFlagValue
	Results []*sacloud.DeleteCacheResult
}

// DeleteCache キャッシュ削除
func (api *WebAccelAPI) DeleteCache(urls ...string) (*WebAccelDeleteCacheResponse, error) {

	type request struct {
		// URL
		URL []string
	}

	uri := fmt.Sprintf("%s/deletecache", api.getResourceURL())

	data, err := api.client.newRequest("POST", uri, &request{URL: urls})
	if err != nil {
		return nil, err
	}

	var res WebAccelDeleteCacheResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
