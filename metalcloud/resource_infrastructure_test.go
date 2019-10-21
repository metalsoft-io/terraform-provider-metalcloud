package metalcloud

import "testing"
import "fmt"
import "os"
import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccInfrastructureResource_basic(t *testing.T) {

	var infrastructure Infrastructure

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccInfrastructureResource(rName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckInfrastructureResourceExists("metalcloud_infrastructure.foo", &infrastructure),
					// verify remote values
					testAccCheckInfrastructureResourceValues(infrastructure, rName),
					// verify local values
					resource.TestCheckResourceAttr("metalcloud_infrastructure.foo", "infrastructure_label", rName),
					resource.TestCheckResourceAttr("metalcloud_infrastructure.foo", "name", rName),
				),
			},
		},
	})
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckInfrastructureResourceExists(n string, infrastructure *Infrastructure) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in Terraform: %s", n)
		}

		// retrieve the configured client from the test setup
		conn := testAccProvider.Meta().(*MetalCloudClient)

		resp, err := conn.getCompleteInfrastructureByLabel(infrastructure.Infrastructure_label)

		if err != nil {
			return err
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*infrastructure = *resp

		return fmt.Errorf("Infrastructure (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckInfrastructureResourceValues(infrastructure Infrastructure, infrastructure_label string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		got := infrastructure.Infrastructure_label
		if got != infrastructure_label {
			return fmt.Errorf("bad infrastructure_label, expected \"%s\", got: %#v", infrastructure_label, got)
		}
		return nil
	}
}

func testAccInfrastructureResource(infrastructure_label string) string {

	user := os.Getenv("METALCLOUD_USER")
	api_key := os.Getenv("METALCLOUD_API_KEY")
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")

	return fmt.Sprintf(`resource "metalcloud_infrastructure" "foo" {
		  infrastructure_label = "%s"
		  datacenter_name = "uk-reading"

		  instance_array {
		      instance_array_label = "test1"
		      instance_array_instance_count = 1
		  }

		  instance_array {
		      instance_array_label = "test1"
		      instance_array_instance_count = 2
		  }

		  user = "%s"
		  api_key = "%s"
		  endpoint = "%s"

		}`, infrastructure_label, user, api_key, endpoint)
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("METALCLOUD_USER"); v == "" {
		t.Fatal("METALCLOUD_USER must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		t.Fatal("METALCLOUD_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		t.Fatal("METALCLOUD_ENDPOINT must be set for acceptance tests")
	}
}
