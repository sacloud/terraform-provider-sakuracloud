package sakuracloud

import (
	"log"
	"net/http"
	"strings"
	"time"

	v1 "github.com/sacloud/libsacloud/api"
	v2 "github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/fake"
	"github.com/sacloud/libsacloud/v2/sacloud/trace"
)

const (
	traceHTTP = "http"
	traceAPI  = "api"
)

const defaultSearchLimit = 10000

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
}

// APIClient for SakuraCloud API
type APIClient struct {
	*v1.Client
	v2.APICaller
	defaultZone string
}

// NewClient returns new API Client for SakuraCloud
func (c *Config) NewClient() *APIClient {
	client := v1.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)

	if c.AcceptLanguage != "" {
		client.AcceptLanguage = c.AcceptLanguage
	}
	if c.APIRootURL != "" {
		v1.SakuraCloudAPIRoot = c.APIRootURL
	}
	if c.RetryMax > 0 {
		client.RetryMax = c.RetryMax
	}
	if c.RetryInterval > 0 {
		client.RetryInterval = time.Duration(c.RetryInterval) * time.Second
	}
	if c.TimeoutMinute > 0 {
		client.DefaultTimeoutDuration = time.Duration(c.TimeoutMinute) * time.Minute
	}

	httpClient := &http.Client{}
	if c.APIRequestTimeout > 0 {
		httpClient.Timeout = time.Duration(c.APIRequestTimeout) * time.Second
	}
	if c.APIRequestRateLimit > 0 {
		httpClient.Transport = &v1.RateLimitRoundTripper{RateLimitPerSec: c.APIRequestRateLimit}
	}
	client.HTTPClient = httpClient

	if c.TraceMode != "" {
		client.TraceMode = true
		log.SetPrefix("[DEBUG] ")
	}
	client.UserAgent = "Terraform for SakuraCloud/v" + Version

	if c.FakeMode != "" {
		if c.FakeStorePath != "" {
			fake.DataStore = fake.NewJSONFileStore(c.FakeStorePath)
		}
		fake.SwitchFactoryFuncToFake()
	}

	v2Client := c.newClientV2()

	return &APIClient{
		Client:      client,
		APICaller:   v2Client,
		defaultZone: c.Zone,
	}
}

func (c *Config) newClientV2() v2.APICaller {
	httpClient := &http.Client{
		Timeout:   time.Duration(c.APIRequestTimeout) * time.Second,
		Transport: &v2.RateLimitRoundTripper{RateLimitPerSec: c.APIRequestRateLimit},
	}
	caller := &v2.Client{
		AccessToken:            c.AccessToken,
		AccessTokenSecret:      c.AccessTokenSecret,
		DefaultTimeoutDuration: time.Duration(c.TimeoutMinute) * time.Minute,
		UserAgent:              "Terraform for SakuraCloud/v" + Version,
		AcceptLanguage:         c.AcceptLanguage,
		RetryMax:               c.RetryMax,
		RetryInterval:          time.Duration(c.RetryInterval) * time.Second,
		HTTPClient:             httpClient,
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
			trace.AddClientFactoryHooks()
		}
		if enableHTTPTrace {
			caller.HTTPClient.Transport = &v2.TracingRoundTripper{
				Transport: caller.HTTPClient.Transport,
			}
		}
	}
	return caller
}
