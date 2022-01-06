// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sacloud/libsacloud/v2/helper/api"
	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/profile"
)

const (
	traceHTTP = "http"
	traceAPI  = "api"
)

const uaEnvVar = "SAKURACLOUD_APPEND_USER_AGENT"

var (
	deletionWaiterTimeout            = 30 * time.Minute
	deletionWaiterPollingInterval    = 5 * time.Second
	databaseWaitAfterCreateDuration  = 1 * time.Minute
	vpcRouterWaitAfterCreateDuration = 1 * time.Minute
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
	defaultZone                      string // 各リソースでzone未指定の場合に利用するゾーン。sacloud.APIDefaultZoneとは別物。
	zones                            []string
	deletionWaiterTimeout            time.Duration
	deletionWaiterPollingInterval    time.Duration
	databaseWaitAfterCreateDuration  time.Duration
	vpcRouterWaitAfterCreateDuration time.Duration
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

	tfUserAgent := terraformUserAgent(c.terraformVersion)
	providerUserAgent := fmt.Sprintf("%s/v%s", "terraform-provider-sakuracloud", Version)
	ua := fmt.Sprintf("%s %s", tfUserAgent, providerUserAgent)
	if add := os.Getenv(uaEnvVar); add != "" {
		ua += " " + add
		log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
	}

	enableAPITrace := false
	enableHTTPTrace := false
	if c.TraceMode != "" {
		enableAPITrace = true
		enableHTTPTrace = true
		mode := strings.ToLower(c.TraceMode)
		switch mode {
		case traceAPI:
			enableHTTPTrace = false
		case traceHTTP:
			enableAPITrace = false
		}
	}

	caller := api.NewCaller(&api.CallerOptions{
		AccessToken:          c.AccessToken,
		AccessTokenSecret:    c.AccessTokenSecret,
		APIRootURL:           c.APIRootURL,
		DefaultZone:          c.DefaultZone,
		AcceptLanguage:       c.AcceptLanguage,
		HTTPClient:           http.DefaultClient,
		HTTPRequestTimeout:   c.APIRequestTimeout,
		HTTPRequestRateLimit: c.APIRequestRateLimit,
		RetryMax:             c.RetryMax,
		RetryWaitMax:         c.RetryWaitMax,
		RetryWaitMin:         c.RetryWaitMin,
		UserAgent:            ua,
		TraceAPI:             enableAPITrace,
		TraceHTTP:            enableHTTPTrace,
		OpenTelemetry:        false,
		OpenTelemetryOptions: nil,
		FakeMode:             c.FakeMode != "",
		FakeStorePath:        c.FakeStorePath,
	})

	zones := c.Zones
	if len(zones) == 0 {
		zones = defaultZones
	}

	// fakeモード有効時は待ち時間を短くしておく
	if c.FakeMode != "" {
		deletionWaiterTimeout = 300 * time.Millisecond // 短すぎるとタイムアウトするため余裕を持たせておく
		deletionWaiterPollingInterval = time.Millisecond
		databaseWaitAfterCreateDuration = time.Millisecond
		vpcRouterWaitAfterCreateDuration = time.Millisecond
	}

	return &APIClient{
		APICaller:                        caller,
		defaultZone:                      c.Zone,
		zones:                            zones,
		deletionWaiterTimeout:            deletionWaiterTimeout,
		deletionWaiterPollingInterval:    deletionWaiterPollingInterval,
		databaseWaitAfterCreateDuration:  databaseWaitAfterCreateDuration,
		vpcRouterWaitAfterCreateDuration: vpcRouterWaitAfterCreateDuration,
	}, nil
}
