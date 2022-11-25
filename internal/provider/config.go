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

package provider

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/terraform-provider-sakuracloud/version"
	"github.com/sacloud/webaccel-api-go"
)

const (
	traceHTTP = "http"
	traceAPI  = "api"
)

const (
	tfUAEnvVar = "TF_APPEND_USER_AGENT"
	uaEnvVar   = "SAKURACLOUD_APPEND_USER_AGENT"
)

var (
	deletionWaiterTimeout            = 30 * time.Minute
	deletionWaiterPollingInterval    = 5 * time.Second
	databaseWaitAfterCreateDuration  = 1 * time.Minute
	vpcRouterWaitAfterCreateDuration = 2 * time.Minute
	httpClient                       = &http.Client{}
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

	apiClient *APIClient
	initOnce  sync.Once
}

func (c *Config) UserAgent() string {
	ua := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io)", c.terraformVersion)

	if add := os.Getenv(tfUAEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	providerUserAgent := fmt.Sprintf("%s/v%s", "terraform-provider-sakuracloud", version.Version)
	ua += " " + providerUserAgent
	if add := os.Getenv(uaEnvVar); add != "" {
		ua += " " + add
		log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
	}

	return ua
}

func (c *Config) enableHTTPTrace() (enableHTTPTrace bool) {
	if c.TraceMode != "" {
		enableHTTPTrace = true
		if strings.ToLower(c.TraceMode) == traceAPI { // 明示的にAPIだけのトレースを指定された場合
			enableHTTPTrace = false
		}
	}
	return
}

func (c *Config) enableAPITrace() (enableAPITrace bool) {
	if c.TraceMode != "" {
		enableAPITrace = true

		if strings.ToLower(c.TraceMode) == traceHTTP { // 明示的にHTTPだけのトレースを指定された場合
			enableAPITrace = false
		}
	}
	return
}

func (c *Config) callerOptions() (*api.CallerOptions, error) {
	opt, err := api.DefaultOptionWithProfile(c.Profile)
	if err != nil {
		return nil, err
	}
	return api.MergeOptions(opt, &api.CallerOptions{
		Options: &client.Options{
			AccessToken:          c.AccessToken,
			AccessTokenSecret:    c.AccessTokenSecret,
			AcceptLanguage:       c.AcceptLanguage,
			HttpClient:           httpClient,
			HttpRequestTimeout:   c.APIRequestTimeout,
			HttpRequestRateLimit: c.APIRequestRateLimit,
			RetryMax:             c.RetryMax,
			RetryWaitMax:         c.RetryWaitMax,
			RetryWaitMin:         c.RetryWaitMin,
			UserAgent:            c.UserAgent(),
			Trace:                c.enableHTTPTrace(),
		},
		APIRootURL:    c.APIRootURL,
		DefaultZone:   c.DefaultZone,
		Zones:         c.Zones,
		TraceAPI:      c.enableAPITrace(),
		FakeMode:      c.FakeMode != "",
		FakeStorePath: c.FakeStorePath,
	}), nil
}

func (c *Config) newClient() (*APIClient, error) {
	opt, err := c.callerOptions()
	if err != nil {
		return nil, err
	}
	caller := api.NewCallerWithOptions(opt)
	if c.Zone == "" {
		profile := opt.ProfileConfigValue()
		if profile != nil {
			c.Zone = profile.Zone
		}
	}
	if c.Zone == "" {
		c.Zone = opt.DefaultZone
	}

	// fakeモード有効時は待ち時間を短くしておく
	if opt.FakeMode || c.FakeMode != "" {
		deletionWaiterTimeout = 300 * time.Millisecond // 短すぎるとタイムアウトするため余裕を持たせておく
		deletionWaiterPollingInterval = time.Millisecond
		databaseWaitAfterCreateDuration = time.Millisecond
		vpcRouterWaitAfterCreateDuration = time.Millisecond
	}

	return &APIClient{
		iaasClient:                       caller,
		defaultZone:                      c.Zone,
		zones:                            opt.Zones,
		deletionWaiterTimeout:            deletionWaiterTimeout,
		deletionWaiterPollingInterval:    deletionWaiterPollingInterval,
		databaseWaitAfterCreateDuration:  databaseWaitAfterCreateDuration,
		vpcRouterWaitAfterCreateDuration: vpcRouterWaitAfterCreateDuration,
		webaccelClient:                   &webaccel.Client{Options: opt.Options},
	}, nil
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() (*APIClient, error) {
	var initErr error
	c.initOnce.Do(func() {
		apiClient, err := c.newClient()
		c.apiClient = apiClient
		initErr = err
	})
	return c.apiClient, initErr
}

// APIClient for SakuraCloud API
type APIClient struct {
	iaasClient     iaas.APICaller
	webaccelClient *webaccel.Client

	defaultZone                      string // 各リソースでzone未指定の場合に利用するゾーン。iaas.APIDefaultZoneとは別物。
	zones                            []string
	deletionWaiterTimeout            time.Duration
	deletionWaiterPollingInterval    time.Duration
	databaseWaitAfterCreateDuration  time.Duration
	vpcRouterWaitAfterCreateDuration time.Duration
}

func (c *APIClient) checkReferencedOption() query.CheckReferencedOption { // nolint:unused
	return query.CheckReferencedOption{
		Tick:    c.deletionWaiterPollingInterval,
		Timeout: c.deletionWaiterTimeout,
	}
}
