// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/api-client-go/profile"
	"github.com/sacloud/apprun-api-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/defaults"
	"github.com/sacloud/terraform-provider-sakuracloud/version"
	"github.com/sacloud/webaccel-api-go"
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
	vpcRouterWaitAfterCreateDuration = 2 * time.Minute
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
	iaas.APICaller
	defaultZone                      string // 各リソースでzone未指定の場合に利用するゾーン。iaas.APIDefaultZoneとは別物。
	zones                            []string
	deletionWaiterTimeout            time.Duration
	deletionWaiterPollingInterval    time.Duration
	databaseWaitAfterCreateDuration  time.Duration
	vpcRouterWaitAfterCreateDuration time.Duration

	webaccelClient *webaccel.Client
	apprunClient   *apprun.Client
	simplemqClient *queue.Client
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
	if (c.Zone == "" || c.Zone == defaults.Zone) && pcv.Zone != "" {
		c.Zone = pcv.Zone
	}

	defaultZones := iaas.SakuraCloudZones
	sort.Strings(c.Zones)
	sort.Strings(defaultZones)
	if (len(c.Zones) == 0 || reflect.DeepEqual(defaultZones, c.Zones)) && len(pcv.Zones) > 0 {
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
	if (c.RetryMax == 0 || c.RetryMax == defaults.RetryMax) && pcv.RetryMax > 0 {
		c.RetryMax = pcv.RetryMax
	}
	if c.RetryWaitMax == 0 {
		c.RetryWaitMax = pcv.RetryWaitMax
	}
	if c.RetryWaitMin == 0 {
		c.RetryWaitMin = pcv.RetryWaitMin
	}
	if (c.APIRequestTimeout == 0 || c.APIRequestTimeout == defaults.APIRequestTimeout) && pcv.HTTPRequestTimeout > 0 {
		c.APIRequestTimeout = pcv.HTTPRequestTimeout
	}
	if (c.APIRequestRateLimit == 0 || c.APIRequestRateLimit == defaults.APIRequestRateLimit) && pcv.HTTPRequestRateLimit > 0 {
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
	providerUserAgent := fmt.Sprintf("%s/v%s", "terraform-provider-sakuracloud", version.Version)
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
	callerOptions := &client.Options{
		AccessToken:          c.AccessToken,
		AccessTokenSecret:    c.AccessTokenSecret,
		AcceptLanguage:       c.AcceptLanguage,
		HttpClient:           http.DefaultClient,
		HttpRequestTimeout:   c.APIRequestTimeout,
		HttpRequestRateLimit: c.APIRequestRateLimit,
		RetryMax:             c.RetryMax,
		RetryWaitMax:         c.RetryWaitMax,
		RetryWaitMin:         c.RetryWaitMin,
		UserAgent:            ua,
		Trace:                enableHTTPTrace,
	}
	caller := api.NewCallerWithOptions(&api.CallerOptions{
		Options:       callerOptions,
		APIRootURL:    c.APIRootURL,
		DefaultZone:   c.DefaultZone,
		TraceAPI:      enableAPITrace,
		FakeMode:      c.FakeMode != "",
		FakeStorePath: c.FakeStorePath,
	})

	zones := c.Zones
	if len(zones) == 0 {
		zones = iaas.SakuraCloudZones
	}

	// fakeモード有効時は待ち時間を短くしておく
	if c.FakeMode != "" {
		deletionWaiterTimeout = 300 * time.Millisecond // 短すぎるとタイムアウトするため余裕を持たせておく
		deletionWaiterPollingInterval = time.Millisecond
		databaseWaitAfterCreateDuration = time.Millisecond
		vpcRouterWaitAfterCreateDuration = time.Millisecond
	}

	simplemqClient, err := simplemq.NewQueueClient(client.WithOptions(callerOptions))
	if err != nil {
		return nil, err
	}

	return &APIClient{
		APICaller:                        caller,
		defaultZone:                      c.Zone,
		zones:                            zones,
		deletionWaiterTimeout:            deletionWaiterTimeout,
		deletionWaiterPollingInterval:    deletionWaiterPollingInterval,
		databaseWaitAfterCreateDuration:  databaseWaitAfterCreateDuration,
		vpcRouterWaitAfterCreateDuration: vpcRouterWaitAfterCreateDuration,
		webaccelClient:                   &webaccel.Client{Options: callerOptions},
		apprunClient:                     &apprun.Client{Options: callerOptions},
		simplemqClient:                   simplemqClient,
	}, nil
}
