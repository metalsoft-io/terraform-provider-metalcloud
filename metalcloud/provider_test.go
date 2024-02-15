package metalcloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"metalcloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		t.Fatal("METALCLOUD_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		t.Fatal("METALCLOUD_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_USER_EMAIL"); v == "" {
		t.Fatal("METALCLOUD_USER_EMAIL must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_USER_ID"); v == "" {
		t.Fatal("METALCLOUD_USER_ID must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_DATACENTER"); v == "" {
		t.Fatal("METALCLOUD_DATACENTER must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_SERVER_TYPE"); v == "" {
		t.Fatal("METALCLOUD_SERVER_TYPE must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_NETWORK_PROFILE_WAN"); v == "" {
		t.Fatal("METALCLOUD_NETWORK_PROFILE_WAN must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_NETWORK_PROFILE_LAN"); v == "" {
		t.Fatal("METALCLOUD_NETWORK_PROFILE_LAN must be set for acceptance tests")
	}

}
