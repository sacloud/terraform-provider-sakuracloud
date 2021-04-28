// Copyright 2016-2021 The Libsacloud Authors
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
	"time"

	"github.com/sacloud/libsacloud/v2"
	"github.com/sacloud/libsacloud/v2/helper/defaults"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
	"github.com/sacloud/libsacloud/v2/sacloud/trace/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	TraceAPI             bool
	TraceHTTP            bool
	OpenTelemetry        bool
	OpenTelemetryOptions []otel.Option

	FakeMode      bool
	FakeStorePath string
}

// NewCaller 指定のオプションでsacloud.APICallerを構築して返す
func NewCaller(opts *CallerOptions) sacloud.APICaller {
	return newCaller(opts)
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
	if opts.OpenTelemetry {
		otel.Initialize(opts.OpenTelemetryOptions...)
		transport := caller.HTTPClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		caller.HTTPClient.Transport = otelhttp.NewTransport(transport)
	}

	if opts.FakeMode {
		if opts.FakeStorePath != "" {
			fake.DataStore = fake.NewJSONFileStore(opts.FakeStorePath)
		}
		// note: exact once
		fake.SwitchFactoryFuncToFake()

		SetupFakeDefaults()
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

func SetupFakeDefaults() {
	defaultInterval := 10 * time.Millisecond

	// update default polling intervals: libsacloud/sacloud
	sacloud.DefaultStatePollingInterval = defaultInterval
	sacloud.DefaultDBStatusPollingInterval = defaultInterval
	// update default polling intervals: libsacloud/helper/setup
	defaults.DefaultDeleteWaitInterval = defaultInterval
	defaults.DefaultProvisioningWaitInterval = defaultInterval
	defaults.DefaultPollingInterval = defaultInterval
	// update default polling intervals: libsacloud/helper/builder
	defaults.DefaultNICUpdateWaitDuration = defaultInterval
	// update default timeouts and span: libsacloud/helper/power
	defaults.DefaultPowerHelperBootRetrySpan = defaultInterval
	defaults.DefaultPowerHelperShutdownRetrySpan = defaultInterval
	defaults.DefaultPowerHelperInitialRequestRetrySpan = defaultInterval
	defaults.DefaultPowerHelperInitialRequestTimeout = defaultInterval * 100
}
