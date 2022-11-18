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
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/api-client-go/profile"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
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

func unsetSakuraCloudEnvs() {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		if strings.HasPrefix(key, "SAKURACLOUD_") {
			os.Unsetenv(key) // nolint
		}
	}
}

func setupConfigTest(t *testing.T, profiles map[string]*profile.ConfigValue, envs map[string]string) func() {
	unsetSakuraCloudEnvs()
	f := initTestProfileDir()
	for profileName, profileValue := range profiles {
		if err := profile.Save(profileName, profileValue); err != nil {
			t.Fatal(err)
		}
	}
	for k, v := range envs {
		os.Setenv(k, v) // nolint
	}
	return func() {
		f()
		for k := range envs {
			os.Unsetenv(k) // nolint
		}
	}
}

var (
	defaultProfile = &profile.ConfigValue{
		AccessToken:          "token",
		AccessTokenSecret:    "secret",
		Zone:                 "dummy1",
		Zones:                []string{"dummy1", "dummy2"},
		DefaultZone:          "is1b",
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
	testProfile = &profile.ConfigValue{
		AccessToken:          "testtoken",
		AccessTokenSecret:    "testsecret",
		Zone:                 "test",
		Zones:                []string{"test1", "test2"},
		DefaultZone:          "is1b",
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
)

func TestConfig_callerOptions(t *testing.T) {
	type fields struct {
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
		terraformVersion    string
	}
	tests := []struct {
		name     string
		fields   fields
		want     *api.CallerOptions
		wantErr  bool
		profiles map[string]*profile.ConfigValue
		envs     map[string]string
	}{
		{
			name:   "minimum",
			fields: fields{},
			want: &api.CallerOptions{
				Options: &client.Options{
					AccessToken:          "token",
					AccessTokenSecret:    "secret",
					HttpClient:           httpClient,
					HttpRequestTimeout:   300,
					HttpRequestRateLimit: 5,
					RetryMax:             10,
					RetryWaitMin:         1,
					RetryWaitMax:         64,
				},
				APIRootURL:  iaas.SakuraCloudAPIRoot,
				DefaultZone: "is1a",
				Zones:       iaas.SakuraCloudZones,
			},
			wantErr: false,
			envs: map[string]string{
				"SAKURACLOUD_ACCESS_TOKEN":        "token",
				"SAKURACLOUD_ACCESS_TOKEN_SECRET": "secret",
			},
		},
		{
			name:   "from default profile",
			fields: fields{},
			want: &api.CallerOptions{
				Options: &client.Options{
					AccessToken:          defaultProfile.AccessToken,
					AccessTokenSecret:    defaultProfile.AccessTokenSecret,
					AcceptLanguage:       defaultProfile.AcceptLanguage,
					Gzip:                 defaultProfile.Gzip,
					HttpClient:           httpClient,
					HttpRequestTimeout:   defaultProfile.HTTPRequestTimeout,
					HttpRequestRateLimit: defaultProfile.HTTPRequestRateLimit,
					RetryMax:             defaultProfile.RetryMax,
					RetryWaitMax:         defaultProfile.RetryWaitMax,
					RetryWaitMin:         defaultProfile.RetryWaitMin,
					UserAgent:            defaultProfile.UserAgent,
					Trace:                defaultProfile.EnableHTTPTrace(),
				},
				APIRootURL:    defaultProfile.APIRootURL,
				DefaultZone:   defaultProfile.DefaultZone,
				Zones:         defaultProfile.Zones,
				TraceAPI:      defaultProfile.EnableAPITrace(),
				FakeMode:      defaultProfile.FakeMode,
				FakeStorePath: defaultProfile.FakeStorePath,
			},
			wantErr: false,
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
				"test":    testProfile,
			},
		},
		{
			name: "from named profile",
			fields: fields{
				Profile: "test",
			},
			want: &api.CallerOptions{
				Options: &client.Options{
					AccessToken:          testProfile.AccessToken,
					AccessTokenSecret:    testProfile.AccessTokenSecret,
					AcceptLanguage:       testProfile.AcceptLanguage,
					Gzip:                 testProfile.Gzip,
					HttpClient:           httpClient,
					HttpRequestTimeout:   testProfile.HTTPRequestTimeout,
					HttpRequestRateLimit: testProfile.HTTPRequestRateLimit,
					RetryMax:             testProfile.RetryMax,
					RetryWaitMax:         testProfile.RetryWaitMax,
					RetryWaitMin:         testProfile.RetryWaitMin,
					UserAgent:            testProfile.UserAgent,
					Trace:                testProfile.EnableHTTPTrace(),
				},
				APIRootURL:    testProfile.APIRootURL,
				DefaultZone:   testProfile.DefaultZone,
				Zones:         testProfile.Zones,
				TraceAPI:      testProfile.EnableAPITrace(),
				FakeMode:      testProfile.FakeMode,
				FakeStorePath: testProfile.FakeStorePath,
			},
			wantErr: false,
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
				"test":    testProfile,
			},
		},
		{
			name: "override profile",
			fields: fields{
				AccessToken:         "token",
				AccessTokenSecret:   "secret",
				Zone:                "zone1",
				Zones:               []string{"zone1", "zone2"},
				DefaultZone:         "zone",
				TraceMode:           "trace",
				FakeMode:            "fake",
				FakeStorePath:       "fakepath",
				AcceptLanguage:      "acceptlanguage",
				APIRootURL:          "apirooturl",
				RetryMax:            1,
				RetryWaitMin:        2,
				RetryWaitMax:        3,
				APIRequestTimeout:   4,
				APIRequestRateLimit: 5,
			},
			want: &api.CallerOptions{
				Options: &client.Options{
					AccessToken:          "token",
					AccessTokenSecret:    "secret",
					AcceptLanguage:       "acceptlanguage",
					Gzip:                 defaultProfile.Gzip,
					HttpClient:           httpClient,
					RetryMax:             1,
					RetryWaitMin:         2,
					RetryWaitMax:         3,
					HttpRequestTimeout:   4,
					HttpRequestRateLimit: 5,
					UserAgent:            defaultProfile.UserAgent,
					Trace:                true,
				},
				APIRootURL:    "apirooturl",
				DefaultZone:   "zone",
				Zones:         []string{"zone1", "zone2"},
				TraceAPI:      true,
				FakeMode:      true,
				FakeStorePath: "fakepath",
			},
			wantErr: false,
			profiles: map[string]*profile.ConfigValue{
				"default": defaultProfile,
				"test":    testProfile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupConfigTest(t, tt.profiles, tt.envs)
			defer cleanup()

			for profileName, profileValue := range tt.profiles {
				if err := profile.Save(profileName, profileValue); err != nil {
					t.Fatal(err)
				}
			}

			c := &Config{
				Profile:             tt.fields.Profile,
				AccessToken:         tt.fields.AccessToken,
				AccessTokenSecret:   tt.fields.AccessTokenSecret,
				Zone:                tt.fields.Zone,
				Zones:               tt.fields.Zones,
				DefaultZone:         tt.fields.DefaultZone,
				TraceMode:           tt.fields.TraceMode,
				FakeMode:            tt.fields.FakeMode,
				FakeStorePath:       tt.fields.FakeStorePath,
				AcceptLanguage:      tt.fields.AcceptLanguage,
				APIRootURL:          tt.fields.APIRootURL,
				RetryMax:            tt.fields.RetryMax,
				RetryWaitMin:        tt.fields.RetryWaitMin,
				RetryWaitMax:        tt.fields.RetryWaitMax,
				APIRequestTimeout:   tt.fields.APIRequestTimeout,
				APIRequestRateLimit: tt.fields.APIRequestRateLimit,
				terraformVersion:    tt.fields.terraformVersion,
			}
			got, err := c.callerOptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("callerOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Note: Equal(t, tt.want, got)としたいが、exportしていないフィールドがあるため個別にEqualを呼んでいる
			require.Equal(t, tt.want.AccessToken, got.AccessToken)
			require.Equal(t, tt.want.AccessTokenSecret, got.AccessTokenSecret)
			require.Equal(t, tt.want.AcceptLanguage, got.AcceptLanguage)
			require.Equal(t, tt.want.Gzip, got.Gzip)
			require.Equal(t, tt.want.HttpClient, got.HttpClient)
			require.Equal(t, tt.want.HttpRequestTimeout, got.HttpRequestTimeout)
			require.Equal(t, tt.want.HttpRequestRateLimit, got.HttpRequestRateLimit)
			require.Equal(t, tt.want.RetryMax, got.RetryMax)
			require.Equal(t, tt.want.RetryWaitMax, got.RetryWaitMax)
			require.Equal(t, tt.want.RetryWaitMin, got.RetryWaitMin)
			require.Equal(t, tt.want.Trace, got.Trace)
			require.Equal(t, tt.want.APIRootURL, got.APIRootURL)
			require.Equal(t, tt.want.DefaultZone, got.DefaultZone)
			require.Equal(t, tt.want.Zones, got.Zones)
			require.Equal(t, tt.want.TraceAPI, got.TraceAPI)
			require.Equal(t, tt.want.FakeMode, got.FakeMode)
			require.Equal(t, tt.want.FakeStorePath, got.FakeStorePath)
		})
	}
}
