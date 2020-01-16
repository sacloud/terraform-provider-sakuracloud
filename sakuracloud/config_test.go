// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/profile"
)

func TestConfig_NewClient_useDefaultHTTPClient(t *testing.T) {
	config := &Config{
		AccessToken:       "dummy",
		AccessTokenSecret: "dummy",
	}

	c1, err := config.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	c2, err := config.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if c1 == c2 {
		t.Errorf("Config.NewClient() should return fresh instance: instance1: %p instance2: %p", c1, c2)
	}

	hc1 := c1.APICaller.(*sacloud.Client).HTTPClient
	hc2 := c2.APICaller.(*sacloud.Client).HTTPClient
	if hc1 != hc2 {
		t.Errorf("APIClient.HTTPClient should use same instance: instance1: %p instance2: %p", hc1, hc2)
	}
}

func initTestProfileDir() func() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("SAKURACLOUD_PROFILE_DIR", wd) // nolint
	profileDir := filepath.Join(wd, ".usacloud")
	if _, err := os.Stat(profileDir); err == nil {
		os.RemoveAll(profileDir) // nolint
	}

	return func() {
		os.RemoveAll(profileDir) // nolint
	}
}

func TestConfig_NewClient_loadFromProfile(t *testing.T) {
	defer initTestProfileDir()()

	defaultProfile := &profile.ConfigValue{
		AccessToken:          "token",
		AccessTokenSecret:    "secret",
		Zone:                 "dummy1",
		Zones:                []string{"dummy1", "dummy2"},
		UserAgent:            "dummy-ua",
		AcceptLanguage:       "ja-JP",
		RetryMax:             1,
		RetryWaitMin:         2,
		RetryWaitMax:         3,
		StatePollingTimeout:  4,
		StatePollingInterval: 5,
		HTTPRequestTimeout:   6,
		HTTPRequestRateLimit: 7,
		APIRootURL:           "dummy",
		TraceMode:            "dummy",
		FakeMode:             true,
		FakeStorePath:        "dummy",
	}
	testProfile := &profile.ConfigValue{
		AccessToken:          "testtoken",
		AccessTokenSecret:    "testsecret",
		Zone:                 "test",
		Zones:                []string{"test1", "test2"},
		UserAgent:            "test-ua",
		AcceptLanguage:       "ja-JP",
		RetryMax:             7,
		RetryWaitMin:         6,
		RetryWaitMax:         5,
		StatePollingTimeout:  4,
		StatePollingInterval: 3,
		HTTPRequestTimeout:   2,
		HTTPRequestRateLimit: 1,
		APIRootURL:           "test",
		TraceMode:            "test",
		FakeMode:             false,
		FakeStorePath:        "test",
	}

	// プロファイル指定なし & デフォルトプロファイルなし
	// プロファイル指定なし & デフォルトプロファイルあり
	// プロファイル指定あり & 指定プロファイルが存在しない
	// プロファイル指定あり 通常

	cases := []struct {
		scenario string
		in       *Config
		profiles map[string]*profile.ConfigValue
		expect   *Config
		err      error
	}{
		{
			scenario: "ProfileName is not specified and Profile is not exists",
			in: &Config{
				Profile:             "",
				Zone:                defaultZone,
				Zones:               defaultZones,
				RetryMax:            defaultRetryMax,
				APIRequestTimeout:   defaultAPIRequestTimeout,
				APIRequestRateLimit: defaultAPIRequestRateLimit,
			},
			profiles: map[string]*profile.ConfigValue{},
			expect: &Config{
				Profile:             "default",
				Zone:                defaultZone,
				Zones:               defaultZones,
				RetryMax:            defaultRetryMax,
				APIRequestTimeout:   defaultAPIRequestTimeout,
				APIRequestRateLimit: defaultAPIRequestRateLimit,
			},
		},
		{
			scenario: "ProfileName is not specified and Profile is exists",
			in: &Config{
				Profile:             "",
				Zone:                defaultZone,
				Zones:               defaultZones,
				RetryMax:            defaultRetryMax,
				APIRequestTimeout:   defaultAPIRequestTimeout,
				APIRequestRateLimit: defaultAPIRequestRateLimit,
			},
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
			},
			expect: &Config{
				Profile:             "default",
				AccessToken:         defaultProfile.AccessToken,
				AccessTokenSecret:   defaultProfile.AccessTokenSecret,
				Zone:                defaultProfile.Zone,
				Zones:               defaultProfile.Zones,
				TraceMode:           defaultProfile.TraceMode,
				FakeMode:            "1",
				FakeStorePath:       defaultProfile.FakeStorePath,
				AcceptLanguage:      defaultProfile.AcceptLanguage,
				APIRootURL:          defaultProfile.APIRootURL,
				RetryMax:            defaultProfile.RetryMax,
				RetryWaitMin:        defaultProfile.RetryWaitMin,
				RetryWaitMax:        defaultProfile.RetryWaitMax,
				APIRequestTimeout:   defaultProfile.HTTPRequestTimeout,
				APIRequestRateLimit: defaultProfile.HTTPRequestRateLimit,
			},
		},
		{
			scenario: "ProfileName is not specified with some values and Profile is exists",
			in: &Config{
				Profile:             "",
				AccessToken:         "from config",
				AccessTokenSecret:   "from config",
				Zone:                "from config",
				Zones:               []string{"zone1", "zone2"},
				TraceMode:           "from config",
				FakeMode:            "from config",
				FakeStorePath:       "from config",
				AcceptLanguage:      "from config",
				APIRootURL:          "from config",
				RetryMax:            8080,
				RetryWaitMin:        8080,
				RetryWaitMax:        8080,
				APIRequestTimeout:   8080,
				APIRequestRateLimit: 8080,
			},
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
			},
			expect: &Config{
				Profile:             "default",
				AccessToken:         "from config",
				AccessTokenSecret:   "from config",
				Zone:                "from config",
				Zones:               []string{"zone1", "zone2"},
				TraceMode:           "from config",
				FakeMode:            "from config",
				FakeStorePath:       "from config",
				AcceptLanguage:      "from config",
				APIRootURL:          "from config",
				RetryMax:            8080,
				RetryWaitMin:        8080,
				RetryWaitMax:        8080,
				APIRequestTimeout:   8080,
				APIRequestRateLimit: 8080,
			},
		},
		{
			scenario: "Profile name specified but not exists",
			in: &Config{
				Profile: "test",
			},
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
			},
			expect: &Config{
				Profile: "test",
			},
			err: errors.New(`loading profile "test" is failed: profile "test" is not exists`),
		},
		{
			scenario: "Profile name specified with normal profile",
			in: &Config{
				Profile:             "test",
				Zone:                defaultZone,
				Zones:               defaultZones,
				RetryMax:            defaultRetryMax,
				APIRequestTimeout:   defaultAPIRequestTimeout,
				APIRequestRateLimit: defaultAPIRequestRateLimit,
			},
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
				"test":    testProfile,
			},
			expect: &Config{
				Profile:             "test",
				AccessToken:         testProfile.AccessToken,
				AccessTokenSecret:   testProfile.AccessTokenSecret,
				Zone:                testProfile.Zone,
				Zones:               testProfile.Zones,
				TraceMode:           testProfile.TraceMode,
				FakeMode:            "",
				FakeStorePath:       testProfile.FakeStorePath,
				AcceptLanguage:      testProfile.AcceptLanguage,
				APIRootURL:          testProfile.APIRootURL,
				RetryMax:            testProfile.RetryMax,
				RetryWaitMin:        testProfile.RetryWaitMin,
				RetryWaitMax:        testProfile.RetryWaitMax,
				APIRequestTimeout:   testProfile.HTTPRequestTimeout,
				APIRequestRateLimit: testProfile.HTTPRequestRateLimit,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.scenario, func(t *testing.T) {
			initTestProfileDir()
			for profileName, profileValue := range tt.profiles {
				if err := profile.Save(profileName, profileValue); err != nil {
					t.Fatal(err)
				}
			}

			err := tt.in.loadFromProfile()
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("got unexpected error: expected: %s got: %s", tt.err, err)
			}
			if !reflect.DeepEqual(tt.expect, tt.in) {
				t.Errorf("got unexpected state: expected: %+v got: %+v", tt.expect, tt.in)
			}
		})
	}
}
