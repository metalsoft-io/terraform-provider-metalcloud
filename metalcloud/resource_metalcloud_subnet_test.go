package metalcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccSubnet_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metalcloud_subnet.subnet01", "subnet_prefix_size", "27"),
				),
			},
			{
				ResourceName:      "metalcloud_subnet.subnet01",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSubnetConfig(t *testing.T) string {
	dc := os.Getenv("METALCLOUD_DATACENTER")

	return fmt.Sprintf(`
	data "metalcloud_infrastructure" "infra" {
   
		infrastructure_label = "vmware-infra-test-test"
		datacenter_name = "%[1]s"
	 
		create_if_not_exists = true
	}

	resource "metalcloud_network" "wan" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_label = "data-network"
		network_type = "wan"
	}
	 
	resource metalcloud_subnet subnet01 {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_id = metalcloud_network.wan.network_id
		subnet_is_ip_range = false
		subnet_prefix_size = 27
		subnet_type = "ipv4"
	}
	 
`, dc)
}
