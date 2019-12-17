// Copyright 2016-2019 The Libsacloud Authors
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

	"github.com/sacloud/libsacloud/v2"
)

var (
	// SakuraCloudAPIRoot APIリクエスト送信先ルートURL(末尾にスラッシュを含まない)
	SakuraCloudAPIRoot = "https://secure.sakura.ad.jp/cloud/zone"
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
	// APIDefaultRetryInterval デフォルトのリトライ間隔
	APIDefaultRetryInterval = 5 * time.Second
)

const (
	// APIAccessTokenEnvKey APIアクセストークンの環境変数名
	APIAccessTokenEnvKey = "SAKURACLOUD_ACCESS_TOKEN"
	// APIAccessSecretEnvKey APIアクセスシークレットの環境変数名
	APIAccessSecretEnvKey = "SAKURACLOUD_ACCESS_TOKEN_SECRET"
)

// APICaller API呼び出し時に利用するトランスポートのインターフェース
type APICaller interface {
	Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error)
}

// Client APIクライアント、APICallerインターフェースを実装する
//
// スレッドセーフではないため複数スレッドから利用する場合は複数のインスタンス生成を推奨
type Client struct {
	// AccessToken アクセストークン
	AccessToken string `validate:"required"`
	// AccessTokenSecret アクセストークンシークレット
	AccessTokenSecret string `validate:"required"`
	// ユーザーエージェント
	UserAgent string
	// Accept-Language
	AcceptLanguage string
	// 503エラー時のリトライ回数
	RetryMax int
	// 503エラー時のリトライ待ち時間
	RetryInterval time.Duration
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
		RetryInterval:     APIDefaultRetryInterval,
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

// Do APIコール実施
func (c *Client) Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error) {
	var (
		client = &retryableHTTPClient{
			Client:        c.HTTPClient,
			retryMax:      c.RetryMax,
			retryInterval: c.RetryInterval,
		}
		err     error
		req     *request
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
	req, err = newRequest(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error with request: %v - %q", url, err)
	}

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

type lenReader interface {
	Len() int
}

type request struct {
	// body is a seekable reader over the request body payload. This is
	// used to rewind the request data in between retries.
	body io.ReadSeeker

	// Embed an HTTP request directly. This makes a *Request act exactly
	// like an *http.Request so that all meta methods are supported.
	*http.Request
}

func newRequest(ctx context.Context, method, url string, body io.ReadSeeker) (*request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}

	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}

	if lr, ok := body.(lenReader); ok {
		httpReq.ContentLength = int64(lr.Len())
	}

	return &request{body, httpReq.WithContext(ctx)}, nil
}

type retryableHTTPClient struct {
	*http.Client
	retryInterval time.Duration
	retryMax      int
}

func (c *retryableHTTPClient) Do(req *request) (*http.Response, error) {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	for i := 0; ; i++ {
		if req.body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("failed to seek body: %v", err)
			}
		}

		res, err := c.Client.Do(req.Request)
		if res != nil && res.StatusCode != http.StatusServiceUnavailable && res.StatusCode != http.StatusLocked {
			return res, err
		}
		if res != nil && res.Body != nil {
			res.Body.Close()
		}

		if err != nil {
			return res, err
		}

		remain := c.retryMax - i
		if remain == 0 {
			break
		}
		time.Sleep(c.retryInterval)
	}

	return nil, fmt.Errorf("%s %s giving up after %d attempts",
		req.Method, req.URL, c.retryMax+1)
}
