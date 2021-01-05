// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/sacloud/libsacloud/v2/helper/defaults"
	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/profile"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
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
	Profile             string
	AccessToken         string
	AccessTokenSecret   string
	Zone                string
	Zones               []string
	DefaultZone         string
	TraceMode           string
	FakeMode            string
	FakeStorePath       string
	AcceptLanguage      string
	APIRootURL          string
	RetryMax            int
	RetryWaitMin        int
	RetryWaitMax        int
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

func (c *APIClient) checkReferencedOption() query.CheckReferencedOption {
	return query.CheckReferencedOption{
		Tick:    c.deletionWaiterPollingInterval,
		Timeout: c.deletionWaiterTimeout,
	}
}

func (c *Config) loadFromProfile() error {
	if c.Profile == "" {
		c.Profile = profile.DefaultProfileName
	}
	if c.Profile != profile.DefaultProfileName {
		log.Printf("[DEBUG] using profile %q", c.Profile)
	}

	pcv := &profile.ConfigValue{}
	if err := profile.Load(c.Profile, pcv); err != nil {
		return fmt.Errorf("loading profile %q is failed: %s", c.Profile, err)
	}

	if c.AccessToken == "" {
		c.AccessToken = pcv.AccessToken
	}
	if c.AccessTokenSecret == "" {
		c.AccessTokenSecret = pcv.AccessTokenSecret
	}
	if c.Zone == defaultZone && pcv.Zone != "" {
		c.Zone = pcv.Zone
	}

	sort.Strings(c.Zones)
	sort.Strings(pcv.Zones)
	if reflect.DeepEqual(defaultZones, c.Zones) && !reflect.DeepEqual(c.Zones, pcv.Zones) && len(pcv.Zones) > 0 {
		c.Zones = pcv.Zones
	}
	if c.DefaultZone != "" {
		sacloud.APIDefaultZone = c.DefaultZone
	}
	if c.TraceMode == "" {
		c.TraceMode = pcv.TraceMode
	}
	if c.FakeMode == "" && pcv.FakeMode {
		c.FakeMode = "1"
	}
	if c.FakeStorePath == "" {
		c.FakeStorePath = pcv.FakeStorePath
	}
	if c.AcceptLanguage == "" {
		c.AcceptLanguage = pcv.AcceptLanguage
	}
	if c.APIRootURL == "" {
		c.APIRootURL = pcv.APIRootURL
	}
	if c.RetryMax == defaultRetryMax && pcv.RetryMax > 0 {
		c.RetryMax = pcv.RetryMax
	}
	if c.RetryWaitMax == 0 {
		c.RetryWaitMax = pcv.RetryWaitMax
	}
	if c.RetryWaitMin == 0 {
		c.RetryWaitMin = pcv.RetryWaitMin
	}
	if c.APIRequestTimeout == defaultAPIRequestTimeout && pcv.HTTPRequestTimeout > 0 {
		c.APIRequestTimeout = pcv.HTTPRequestTimeout
	}
	if c.APIRequestRateLimit == defaultAPIRequestRateLimit && pcv.HTTPRequestRateLimit > 0 {
		c.APIRequestRateLimit = pcv.HTTPRequestRateLimit
	}
	return nil
}

func (c *Config) validate() error {
	var err error
	if c.AccessToken == "" {
		err = multierror.Append(err, errors.New("AccessToken is required"))
	}
	if c.AccessTokenSecret == "" {
		err = multierror.Append(err, errors.New("AccessTokenSecret is required"))
	}
	return err
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() (*APIClient, error) {
	if err := c.loadFromProfile(); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}

	tfUserAgent := httpclient.TerraformUserAgent(c.terraformVersion)
	providerUserAgent := fmt.Sprintf("%s/v%s", "terraform-provider-sakuracloud", Version)
	ua := fmt.Sprintf("%s %s", tfUserAgent, providerUserAgent)
	if add := os.Getenv(uaEnvVar); add != "" {
		ua += " " + add
		log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Duration(c.APIRequestTimeout) * time.Second
	httpClient.Transport = &sacloud.RateLimitRoundTripper{RateLimitPerSec: c.APIRequestRateLimit}

	retryWaitMax := sacloud.APIDefaultRetryWaitMax
	retryWaitMin := sacloud.APIDefaultRetryWaitMin
	if c.RetryWaitMax > 0 {
		retryWaitMax = time.Duration(c.RetryWaitMax) * time.Second
	}
	if c.RetryWaitMin > 0 {
		retryWaitMin = time.Duration(c.RetryWaitMin) * time.Second
	}

	caller := &sacloud.Client{
		AccessToken:       c.AccessToken,
		AccessTokenSecret: c.AccessTokenSecret,
		UserAgent:         ua,
		AcceptLanguage:    c.AcceptLanguage,
		RetryMax:          c.RetryMax,
		RetryWaitMax:      retryWaitMax,
		RetryWaitMin:      retryWaitMin,
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

		deletionWaiterTimeout = 10 * time.Second
		defaultInterval := 10 * time.Millisecond
		deletionWaiterPollingInterval = defaultInterval
		databaseWaitAfterCreateDuration = defaultInterval

		// update default polling intervals: libsacloud/sacloud
		sacloud.DefaultStatePollingInterval = defaultInterval
		sacloud.DefaultDBStatusPollingInterval = defaultInterval
		defaults.DefaultDeleteWaitInterval = defaultInterval
		defaults.DefaultProvisioningWaitInterval = defaultInterval
		defaults.DefaultPollingInterval = defaultInterval
		defaults.DefaultNICUpdateWaitDuration = defaultInterval
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
	}, nil
}
