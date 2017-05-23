package sakuracloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testTargetZone = "is1b"

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
	if v := os.Getenv("SAKURACLOUD_ACCESS_TOKEN"); v == "" {
		t.Fatal("SAKURACLOUD_ACCESS_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"); v == "" {
		t.Fatal("SAKURACLOUD_ACCESS_TOKEN_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("SAKURACLOUD_ZONE"); v == "" {
		os.Setenv("SAKURACLOUD_ZONE", testTargetZone)
	}
}
