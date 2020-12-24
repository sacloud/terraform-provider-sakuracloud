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

package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sacloud/libsacloud/v2/helper/defaults"

	"github.com/sacloud/libsacloud/v2"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
)

// CallerOptions sacloud.APICallerを作成する際のオプション
type CallerOptions struct {
	AccessToken       string
	AccessTokenSecret string

	APIRootURL     string
	DefaultZone    string
	AcceptLanguage string

	HTTPClient *http.Client

	HTTPRequestTimeout   int
	HTTPRequestRateLimit int

	RetryMax     int
	RetryWaitMax int
	RetryWaitMin int

	UserAgent string

	TraceAPI  bool
	TraceHTTP bool

	FakeMode      bool
	FakeStorePath string
}

var clientOnce sync.Once
var cachedCaller sacloud.APICaller

// NewCaller 指定のオプションでsacloud.APICallerを構築して返す
// sacloud.APICallerはプロセス中はキャッシュされる(NewCallerを用いてオプションの違うCallerは取得不可)
func NewCaller(opts *CallerOptions) sacloud.APICaller {
	clientOnce.Do(func() {
		cachedCaller = newCaller(opts)
	})

	return cachedCaller
}

func newCaller(opts *CallerOptions) sacloud.APICaller {
	// build http client
	httpClient := http.DefaultClient
	if opts.HTTPClient != nil {
		httpClient = opts.HTTPClient
	}
	if opts.HTTPRequestTimeout > 0 {
		httpClient.Timeout = time.Duration(opts.HTTPRequestTimeout) * time.Second
	}
	if opts.HTTPRequestRateLimit > 0 {
		httpClient.Transport = &sacloud.RateLimitRoundTripper{RateLimitPerSec: opts.HTTPRequestRateLimit}
	}

	retryMax := sacloud.APIDefaultRetryMax
	if opts.RetryMax > 0 {
		retryMax = opts.RetryMax
	}

	retryWaitMax := sacloud.APIDefaultRetryWaitMax
	if opts.RetryWaitMax > 0 {
		retryWaitMax = time.Duration(opts.RetryWaitMax) * time.Second
	}

	retryWaitMin := sacloud.APIDefaultRetryWaitMin
	if opts.RetryWaitMin > 0 {
		retryWaitMin = time.Duration(opts.RetryWaitMin) * time.Second
	}

	ua := fmt.Sprintf("libsacloud/%s", libsacloud.Version)
	if opts.UserAgent != "" {
		ua = opts.UserAgent
	}

	caller := &sacloud.Client{
		AccessToken:       opts.AccessToken,
		AccessTokenSecret: opts.AccessTokenSecret,
		UserAgent:         ua,
		AcceptLanguage:    opts.AcceptLanguage,
		RetryMax:          retryMax,
		RetryWaitMax:      retryWaitMax,
		RetryWaitMin:      retryWaitMin,
		HTTPClient:        httpClient,
	}
	sacloud.DefaultStatePollingTimeout = 72 * time.Hour

	if opts.TraceAPI {
		// note: exact once
		trace.AddClientFactoryHooks()
	}
	if opts.TraceHTTP {
		caller.HTTPClient.Transport = &sacloud.TracingRoundTripper{
			Transport: caller.HTTPClient.Transport,
		}
	}

	if opts.FakeMode {
		if opts.FakeStorePath != "" {
			fake.DataStore = fake.NewJSONFileStore(opts.FakeStorePath)
		}
		// note: exact once
		fake.SwitchFactoryFuncToFake()

		defaultInterval := 10 * time.Millisecond

		// update default polling intervals: libsacloud/sacloud
		sacloud.DefaultStatePollingInterval = defaultInterval
		sacloud.DefaultDBStatusPollingInterval = defaultInterval
		// update default polling intervals: libsacloud/utils/setup
		defaults.DefaultDeleteWaitInterval = defaultInterval
		defaults.DefaultProvisioningWaitInterval = defaultInterval
		defaults.DefaultPollingInterval = defaultInterval
		// update default polling intervals: libsacloud/utils/builder
		defaults.DefaultNICUpdateWaitDuration = defaultInterval
	}

	if opts.DefaultZone != "" {
		sacloud.APIDefaultZone = opts.DefaultZone
	}

	if opts.APIRootURL != "" {
		if strings.HasSuffix(opts.APIRootURL, "/") {
			opts.APIRootURL = strings.TrimRight(opts.APIRootURL, "/")
		}
		sacloud.SakuraCloudAPIRoot = opts.APIRootURL
	}
	return caller
}
