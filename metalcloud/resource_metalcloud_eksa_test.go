package metalcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccEKS_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEKSConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metalcloud_eksa.cluster01", "cluster_label", "test-eksa"),
				),
			},
			{
				ResourceName:      "metalcloud_eksa.cluster01",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEKSConfig(t *testing.T) string {
	dc := os.Getenv("METALCLOUD_DATACENTER")
	serverType := os.Getenv("METALCLOUD_SERVER_TYPE")
	networkProfileWan := os.Getenv("METALCLOUD_NETWORK_PROFILE_WAN")

	return fmt.Sprintf(`
	data "metalcloud_infrastructure" "infra" {
   
		infrastructure_label = "infra-test-eks"
		datacenter_name = "%[1]s"
	 
		create_if_not_exists = true
	}
	 
	data "metalcloud_server_type" "large"{
		 server_type_name = "%[2]s"
	}
	 
	resource "metalcloud_network" "wan" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_type = "wan"
	}
	 
	resource "metalcloud_network" "san" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_type = "lan"
	}
	 
	 
	data "metalcloud_network_profile" "eksa-mgmt"{
		network_profile_label = "%[3]s"
		datacenter_name = "%[1]s"
	}
	 

	data "metalcloud_network_profile" "eksa-control-plane"{
		network_profile_label = "%[3]s"
		datacenter_name = "%[1]s"
	}

	data "metalcloud_network_profile" "eksa-workload"{
		network_profile_label = "%[3]s"
		datacenter_name = "%[1]s"
	}
	
	resource "metalcloud_eksa" "cluster01" {
		infrastructure_id =  data.metalcloud_infrastructure.infra.infrastructure_id
	 
		cluster_label = "test-eksa"

		
		instance_array_instance_count_eksa_mgmt = 1
		instance_array_instance_count_mgmt = 1
		instance_array_instance_count_worker = 1
	 
		instance_server_type_eksa_mgmt {
			instance_index = 0
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
		instance_server_type_mgmt {
			instance_index = 0
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
		instance_server_type_worker {
			instance_index = 0
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
	 
		interface_eksa_mgmt{
		  interface_index = 0
		  network_id = metalcloud_network.wan.id
		}
	 
		interface_eksa_mgmt{
		  interface_index = 1
		  network_id = metalcloud_network.san.id
		}
	 
	 
		interface_mgmt{
		  interface_index = 0
		  network_id = metalcloud_network.wan.id
		}
	 
		interface_mgmt {
		  interface_index = 1
		  network_id = metalcloud_network.san.id
		}
	 
		interface_worker {
		  interface_index = 0
		  network_id = metalcloud_network.wan.id
		}
	 
		interface_worker {
		  interface_index = 1
		  network_id = metalcloud_network.san.id
		}
	 
		instance_array_network_profile_eksa_mgmt {
			network_id = metalcloud_network.wan.id
			network_profile_id = data.metalcloud_network_profile.eksa-mgmt.id
		}

		instance_array_network_profile_eksa_mgmt {
			network_id = metalcloud_network.wan.id
			network_profile_id = data.metalcloud_network_profile.eksa-control-plane.id
		}

		instance_array_network_profile_worker {
			network_id = metalcloud_network.wan.id
			network_profile_id = data.metalcloud_network_profile.eksa-workload.id
		}
	}
	 
`, dc, serverType, networkProfileWan)
}
