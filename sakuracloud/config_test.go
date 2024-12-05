// Copyright 2016-2023 terraform-provider-sakuracloud authors
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

	"github.com/sacloud/api-client-go/profile"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/defaults"
	"github.com/stretchr/testify/require"
)

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
			scenario: "ProfileName is not specified and Profile does not exist",
			in: &Config{
				Profile:             "",
				Zone:                defaults.Zone,
				Zones:               iaas.SakuraCloudZones,
				RetryMax:            defaults.RetryMax,
				APIRequestTimeout:   defaults.APIRequestTimeout,
				APIRequestRateLimit: defaults.APIRequestRateLimit,
			},
			profiles: map[string]*profile.ConfigValue{},
			expect: &Config{
				Profile:             "default",
				Zone:                defaults.Zone,
				Zones:               iaas.SakuraCloudZones,
				RetryMax:            defaults.RetryMax,
				APIRequestTimeout:   defaults.APIRequestTimeout,
				APIRequestRateLimit: defaults.APIRequestRateLimit,
			},
		},
		{
			scenario: "ProfileName is not specified but Profile does exist",
			in: &Config{
				Profile:             "",
				Zone:                defaults.Zone,
				Zones:               iaas.SakuraCloudZones,
				RetryMax:            defaults.RetryMax,
				APIRequestTimeout:   defaults.APIRequestTimeout,
				APIRequestRateLimit: defaults.APIRequestRateLimit,
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
			scenario: "Config has some values although Profile Name is not specified while Profile does exist",
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
			scenario: "Profile Name is specified although Profile does not exist",
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
			scenario: "Profile Name is specified with a regular Profile with only zone being set explicitly",
			in: &Config{
				Profile: "test",
				Zone:    defaults.Zone,
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
		{
			scenario: "Profile Name is specified with a regular Profile with other values being default",
			in: &Config{
				Profile:             "test",
				Zone:                defaults.Zone,
				Zones:               iaas.SakuraCloudZones,
				RetryMax:            defaults.RetryMax,
				APIRequestTimeout:   defaults.APIRequestTimeout,
				APIRequestRateLimit: defaults.APIRequestRateLimit,
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
		{
			scenario: "Profile Name is specified with a regular Profile with other values being default",
			in: &Config{
				Profile:             "test",
				Zone:                defaults.Zone,
				Zones:               iaas.SakuraCloudZones,
				RetryMax:            defaults.RetryMax,
				APIRequestTimeout:   defaults.APIRequestTimeout,
				APIRequestRateLimit: defaults.APIRequestRateLimit,
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
		{
			scenario: "Profile Name is specified with a regular Profile with no other values being set",
			in: &Config{
				Profile: "test",
			},
			profiles: map[string]*profile.ConfigValue{
				"test": testProfile,
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
			require.EqualValues(t, tt.expect, tt.in)
		})
	}
}
