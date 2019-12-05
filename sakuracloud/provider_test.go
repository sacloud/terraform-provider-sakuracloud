package sakuracloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var (
	testDefaultTargetZone   = "is1b"
	testDefaultAPIRetryMax  = "30"
	testDefaultAPIRateLimit = "5"
)

func init() {
	if v := os.Getenv("SAKURACLOUD_TEST_ZONE"); v != "" {
		os.Setenv("SAKURACLOUD_ZONE", v)
	}
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sakuracloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {

	requiredEnvs := []string{
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
		"SACLOUD_OJS_ACCESS_KEY_ID",
		"SACLOUD_OJS_SECRET_ACCESS_KEY",
	}

	for _, env := range requiredEnvs {
		if v := os.Getenv(env); v == "" {
			t.Fatal(fmt.Sprintf("%s must be set for acceptance tests", env))
		}
	}

	if v := os.Getenv("SAKURACLOUD_ZONE"); v == "" {
		os.Setenv("SAKURACLOUD_ZONE", testDefaultTargetZone)
	}

	if v := os.Getenv("SAKURACLOUD_RETRY_MAX"); v == "" {
		os.Setenv("SAKURACLOUD_RETRY_MAX", testDefaultAPIRetryMax)
	}

	if v := os.Getenv("SAKURACLOUD_RATE_LIMIT"); v == "" {
		os.Setenv("SAKURACLOUD_RATE_LIMIT", testDefaultAPIRateLimit)
	}
}

func skipIfFakeModeEnabled(t *testing.T) {
	if isFakeModeEnabled() {
		t.Skip("This test runs only without FAKE_MODE environment variables")
	}
}

func isFakeModeEnabled() bool {
	fakeMode := os.Getenv("FAKE_MODE")
	return fakeMode != ""
}
