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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
	"github.com/sacloud/libsacloud/v2/utils/builder"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

const (
	traceHTTP = "http"
	traceAPI  = "api"
)

const uaEnvVar = "SAKURACLOUD_APPEND_USER_AGENT"

var (
	fakeModeOnce                    sync.Once
	v2ClientOnce                    sync.Once
	deletionWaiterTimeout           = 30 * time.Minute
	deletionWaiterPollingInterval   = 5 * time.Second
	databaseWaitAfterCreateDuration = 1 * time.Minute
)

// Config type of SakuraCloud Config
type Config struct {
	AccessToken         string
	AccessTokenSecret   string
	Zone                string
	Zones               []string
	TraceMode           string
	FakeMode            string
	FakeStorePath       string
	AcceptLanguage      string
	APIRootURL          string
	RetryMax            int
	RetryInterval       int
	APIRequestTimeout   int
	APIRequestRateLimit int

	terraformVersion string
}

// APIClient for SakuraCloud API
type APIClient struct {
	sacloud.APICaller
	defaultZone                     string
	zones                           []string
	deletionWaiterTimeout           time.Duration
	deletionWaiterPollingInterval   time.Duration
	databaseWaitAfterCreateDuration time.Duration
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *APIClient {
	tfUserAgent := httpclient.TerraformUserAgent(c.terraformVersion)
	providerUserAgent := fmt.Sprintf("%s/v%s", "terraform-provider-sakuracloud", Version)
	ua := fmt.Sprintf("%s %s", tfUserAgent, providerUserAgent)
	if add := os.Getenv(uaEnvVar); add != "" {
		ua += " " + add
		log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(c.APIRequestTimeout) * time.Second,
		Transport: &sacloud.RateLimitRoundTripper{RateLimitPerSec: c.APIRequestRateLimit},
	}
	caller := &sacloud.Client{
		AccessToken:       c.AccessToken,
		AccessTokenSecret: c.AccessTokenSecret,
		UserAgent:         ua,
		AcceptLanguage:    c.AcceptLanguage,
		RetryMax:          c.RetryMax,
		RetryInterval:     time.Duration(c.RetryInterval) * time.Second,
		HTTPClient:        httpClient,
	}
	sacloud.DefaultStatePollingTimeout = 72 * time.Hour

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
		defaultInterval := 10 * time.Millisecond
		deletionWaiterPollingInterval = defaultInterval
		databaseWaitAfterCreateDuration = defaultInterval

		// update default polling intervals: libsacloud/sacloud
		sacloud.DefaultStatePollingInterval = defaultInterval
		sacloud.APIDefaultRetryInterval = defaultInterval
		// update default polling intervals: libsacloud/utils/setup
		setup.DefaultDeleteWaitInterval = defaultInterval
		setup.DefaultProvisioningWaitInterval = defaultInterval
		setup.DefaultPollingInterval = defaultInterval
		// update default polling intervals: libsacloud/utils/builder
		builder.DefaultNICUpdateWaitDuration = defaultInterval
	}

	zones := c.Zones
	if len(zones) == 0 {
		zones = defaultZones
	}
	if c.APIRootURL != "" {
		if strings.HasSuffix(c.APIRootURL, "/") {
			c.APIRootURL = strings.TrimRight(c.APIRootURL, "/")
		}
		sacloud.SakuraCloudAPIRoot = c.APIRootURL
	}

	return &APIClient{
		APICaller:                       caller,
		defaultZone:                     c.Zone,
		zones:                           zones,
		deletionWaiterTimeout:           deletionWaiterTimeout,
		deletionWaiterPollingInterval:   deletionWaiterPollingInterval,
		databaseWaitAfterCreateDuration: databaseWaitAfterCreateDuration,
	}
}
