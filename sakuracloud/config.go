// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

package sakuracloud

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"

	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
)

const (
	traceHTTP = "http"
	traceAPI  = "api"
)

var (
	fakeModeOnce                  sync.Once
	v2ClientOnce                  sync.Once
	deletionWaiterTimeout         = 30 * time.Minute
	deletionWaiterPollingInterval = 5 * time.Second
)

// Config type of SakuraCloud Config
type Config struct {
	AccessToken         string
	AccessTokenSecret   string
	Zone                string
	TimeoutMinute       int
	TraceMode           string
	FakeMode            string
	FakeStorePath       string
	AcceptLanguage      string
	APIRootURL          string
	RetryMax            int
	RetryInterval       int
	APIRequestTimeout   int
	APIRequestRateLimit int

	initOnce sync.Once
}

// APIClient for SakuraCloud API
type APIClient struct {
	sacloud.APICaller
	defaultZone                   string
	deletionWaiterTimeout         time.Duration
	deletionWaiterPollingInterval time.Duration
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *APIClient {

	httpClient := &http.Client{
		Timeout:   time.Duration(c.APIRequestTimeout) * time.Second,
		Transport: &sacloud.RateLimitRoundTripper{RateLimitPerSec: c.APIRequestRateLimit},
	}
	caller := &sacloud.Client{
		AccessToken:            c.AccessToken,
		AccessTokenSecret:      c.AccessTokenSecret,
		DefaultTimeoutDuration: time.Duration(c.TimeoutMinute) * time.Minute,
		UserAgent:              "Terraform for SakuraCloud/v" + Version,
		AcceptLanguage:         c.AcceptLanguage,
		RetryMax:               c.RetryMax,
		RetryInterval:          time.Duration(c.RetryInterval) * time.Second,
		HTTPClient:             httpClient,
	}
	sacloud.DefaultStatePollTimeout = time.Duration(c.TimeoutMinute) * time.Minute

	if c.FakeMode != "" {
		sacloud.DefaultStatePollInterval = 10 * time.Millisecond
		sacloud.APIDefaultRetryInterval = 10 * time.Millisecond
	}

	if c.TraceMode != "" {
		enableAPITrace := true
		enableHTTPTrace := true

		mode := strings.ToLower(c.TraceMode)
		switch mode {
		case traceAPI:
			enableHTTPTrace = false
		case traceHTTP:
			enableAPITrace = false
		}

		if enableAPITrace {
			v2ClientOnce.Do(func() {
				trace.AddClientFactoryHooks()
			})
		}
		if enableHTTPTrace {
			caller.HTTPClient.Transport = &sacloud.TracingRoundTripper{
				Transport: caller.HTTPClient.Transport,
			}
		}
	}

	if c.FakeMode != "" {
		if c.FakeStorePath != "" {
			fake.DataStore = fake.NewJSONFileStore(c.FakeStorePath)
		}
		fakeModeOnce.Do(func() {
			fake.SwitchFactoryFuncToFake()
		})

		// TODO パラメータ化
		deletionWaiterTimeout = 10 * time.Second
		deletionWaiterPollingInterval = 10 * time.Millisecond
	}

	return &APIClient{
		APICaller:                     caller,
		defaultZone:                   c.Zone,
		deletionWaiterTimeout:         deletionWaiterTimeout,
		deletionWaiterPollingInterval: deletionWaiterPollingInterval,
	}
}
