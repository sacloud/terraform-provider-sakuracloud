// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sacloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sacloud/libsacloud/v2"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

var (
	// SakuraCloudAPIRoot APIリクエスト送信先ルートURL(末尾にスラッシュを含まない)
	SakuraCloudAPIRoot = "https://secure.sakura.ad.jp/cloud/zone"

	// SakuraCloudZones 利用可能なゾーンのデフォルト値
	SakuraCloudZones = types.ZoneNames
)

var (
	// APIDefaultZone デフォルトゾーン、グローバルリソースなどで利用される
	APIDefaultZone = "is1a"
	// APIDefaultTimeoutDuration デフォルトのタイムアウト
	APIDefaultTimeoutDuration = 20 * time.Minute
	//APIDefaultUserAgent デフォルトのユーザーエージェント
	APIDefaultUserAgent = fmt.Sprintf("libsacloud/%s", libsacloud.Version)
	// APIDefaultAcceptLanguage デフォルトのAcceptLanguage
	APIDefaultAcceptLanguage = ""
	// APIDefaultRetryMax デフォルトのリトライ回数
	APIDefaultRetryMax = 0
	// APIDefaultRetryWaitMin デフォルトのリトライ間隔(最小)
	APIDefaultRetryWaitMin = 1 * time.Second
	// APIDefaultRetryWaitMax デフォルトのリトライ間隔(最大)
	APIDefaultRetryWaitMax = 64 * time.Second
)

const (
	// APIAccessTokenEnvKey APIアクセストークンの環境変数名
	APIAccessTokenEnvKey = "SAKURACLOUD_ACCESS_TOKEN"
	// APIAccessSecretEnvKey APIアクセスシークレットの環境変数名
	APIAccessSecretEnvKey = "SAKURACLOUD_ACCESS_TOKEN_SECRET"
)

// APICaller API呼び出し時に利用するトランスポートのインターフェース sacloud.Clientなどで実装される
type APICaller interface {
	Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error)
}

// Client APIクライアント、APICallerインターフェースを実装する
//
// レスポンスステータスコード423、または503を受け取った場合、RetryMax回リトライする
// リトライ間隔はRetryMinからRetryMaxまで指数的に増加する(Exponential Backoff)
//
// リトライ時にcontext.Canceled、またはcontext.DeadlineExceededの場合はリトライしない
type Client struct {
	// AccessToken アクセストークン
	AccessToken string `validate:"required"`
	// AccessTokenSecret アクセストークンシークレット
	AccessTokenSecret string `validate:"required"`
	// ユーザーエージェント
	UserAgent string
	// Accept-Language
	AcceptLanguage string
	// 423/503エラー時のリトライ回数
	RetryMax int
	// 423/503エラー時のリトライ待ち時間(最小)
	RetryWaitMin time.Duration
	// 423/503エラー時のリトライ待ち時間(最大)
	RetryWaitMax time.Duration
	// APIコール時に利用される*http.Client 未指定の場合http.DefaultClientが利用される
	HTTPClient *http.Client
}

// NewClient APIクライアント作成
func NewClient(token, secret string) *Client {
	c := &Client{
		AccessToken:       token,
		AccessTokenSecret: secret,
		UserAgent:         APIDefaultUserAgent,
		AcceptLanguage:    APIDefaultAcceptLanguage,
		RetryMax:          APIDefaultRetryMax,
		RetryWaitMin:      APIDefaultRetryWaitMin,
		RetryWaitMax:      APIDefaultRetryWaitMax,
	}
	return c
}

// NewClientFromEnv 環境変数からAPIキーを取得してAPIクライアントを作成する
func NewClientFromEnv() (*Client, error) {
	token := os.Getenv(APIAccessTokenEnvKey)
	if token == "" {
		return nil, fmt.Errorf("environment variable %q is required", APIAccessTokenEnvKey)
	}
	secret := os.Getenv(APIAccessSecretEnvKey)
	if secret == "" {
		return nil, fmt.Errorf("environment variable %q is required", APIAccessSecretEnvKey)
	}
	return NewClient(token, secret), nil
}

func (c *Client) isOkStatus(code int) bool {
	codes := map[int]bool{
		http.StatusOK:        true,
		http.StatusCreated:   true,
		http.StatusAccepted:  true,
		http.StatusNoContent: true,
	}
	_, ok := codes[code]
	return ok
}

func (c *Client) httpClient() *retryablehttp.Client {
	return &retryablehttp.Client{
		HTTPClient:   c.HTTPClient,
		RetryWaitMin: c.RetryWaitMin,
		RetryWaitMax: c.RetryWaitMax,
		RetryMax:     c.RetryMax,
		CheckRetry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			if err != nil {
				return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
			}
			if resp.StatusCode == 0 || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusLocked {
				return true, nil
			}
			return false, nil
		},
		Backoff: retryablehttp.DefaultBackoff,
	}
}

// Do APIコール実施
func (c *Client) Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error) {
	var (
		client  = c.httpClient()
		err     error
		strBody string
	)

	// setup url and body
	var url = uri
	var bodyReader io.ReadSeeker
	if body != nil {
		var bodyJSON []byte
		bodyJSON, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if method == "GET" {
			url = fmt.Sprintf("%s?%s", url, bytes.NewBuffer(bodyJSON))
		} else {
			bodyReader = bytes.NewReader(bodyJSON)
		}
	}
	req, err := retryablehttp.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error with request: %v - %q", url, err)
	}
	req = req.WithContext(ctx)

	// set headers
	req.SetBasicAuth(c.AccessToken, c.AccessTokenSecret)
	req.Header.Add("X-Sakura-Bigint-As-Int", "1") //Use BigInt on resource ids.
	req.Header.Add("User-Agent", c.UserAgent)
	if c.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", c.AcceptLanguage)
	}
	req.Method = method

	// API call
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint - ignore error

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !c.isOkStatus(resp.StatusCode) {
		errResponse := &APIErrorResponse{}
		err := json.Unmarshal(data, errResponse)
		if err != nil {
			return nil, fmt.Errorf("error in response: %s", string(data))
		}
		return nil, NewAPIError(req.Method, req.URL, strBody, resp.StatusCode, errResponse)
	}

	return data, nil
}
