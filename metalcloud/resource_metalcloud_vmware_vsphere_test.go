package metalcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccVMWareVSphere_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMVsphereConfig(t),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("metalcloud_vmware_vsphere.VMWareVsphere", "cluster_label", "testvmware"),
					resource.TestCheckResourceAttr("metalcloud_vmware_vsphere.VMWareVsphere", "instance_server_type_master.#", "1"),
				),
			},
			{
				ResourceName:      "metalcloud_vmware_vsphere.VMWareVsphere",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVMVsphereConfig(t *testing.T) string {
	dc := os.Getenv("METALCLOUD_DATACENTER")
	serverType := os.Getenv("METALCLOUD_SERVER_TYPE")
	networkProfileWan := os.Getenv("METALCLOUD_NETWORK_PROFILE_WAN")
	networkProfileLan := os.Getenv("METALCLOUD_NETWORK_PROFILE_LAN")

	return fmt.Sprintf(`
	data "metalcloud_infrastructure" "infra" {
   
		infrastructure_label = "vmware-infra-test-test"
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
	 
	resource "metalcloud_network" "lan1" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_type = "lan"
	}
	 
	resource "metalcloud_network" "lan2" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_type = "lan"
	}
	 
	resource "metalcloud_network" "lan3" {
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		network_type = "lan"
	}
	 
	data "metalcloud_network_profile" "vmware_wan"{
		network_profile_label = "%[3]s"
		datacenter_name = "%[1]s"
	}
	 

	data "metalcloud_network_profile" "vmware_lan"{
		network_profile_label = "%[4]s"
		datacenter_name = "%[1]s"
	}
	
	resource "metalcloud_vmware_vsphere" "VMWareVsphere" {
		infrastructure_id =  data.metalcloud_infrastructure.infra.infrastructure_id
	 
		cluster_label = "testvmware"
		instance_array_instance_count_master = 1
		instance_array_instance_count_worker = 2
	 
		instance_server_type_master {
			instance_index = 0
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
		instance_server_type_worker {
			instance_index = 0
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
		instance_server_type_worker {
			instance_index = 1
			server_type_id = data.metalcloud_server_type.large.server_type_id
		}
	 
	 
		interface_master{
		  interface_index = 0
		  network_id = metalcloud_network.wan.id
		}
	 
		interface_master{
		  interface_index = 1
		  network_id = metalcloud_network.lan1.id
		}
	 
		interface_master {
		  interface_index = 2
		  network_id = metalcloud_network.lan2.id
		}
	 
		interface_master {
		  interface_index = 3
		  network_id = metalcloud_network.lan3.id
		}
	 
		interface_worker{
		  interface_index = 0
		  network_id = metalcloud_network.wan.id
		}
	 
		interface_worker {
		  interface_index = 1
		  network_id = metalcloud_network.lan1.id
		}
	 
		interface_worker {
		  interface_index = 2
		  network_id = metalcloud_network.lan2.id
		}
	 
		interface_worker {
		  interface_index = 3
		  network_id = metalcloud_network.lan3.id
		}
	 
		instance_array_network_profile_master {
			network_id = metalcloud_network.wan.id
			network_profile_id = data.metalcloud_network_profile.vmware_wan.id
		}
	 
		instance_array_network_profile_master {
			network_id = metalcloud_network.lan1.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_network_profile_master {
			network_id = metalcloud_network.lan2.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_network_profile_master {
			network_id = metalcloud_network.lan3.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_network_profile_worker {
			network_id = metalcloud_network.wan.id
			network_profile_id = data.metalcloud_network_profile.vmware_wan.id
		}
	 
		instance_array_network_profile_worker {
			network_id = metalcloud_network.lan1.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_network_profile_worker {
			network_id = metalcloud_network.lan2.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_network_profile_worker {
			network_id = metalcloud_network.lan3.id
			network_profile_id = data.metalcloud_network_profile.vmware_lan.id
		}
	 
		instance_array_custom_variables_master = {      
			"vcsa_ip"= "192.168.177.2",
			"vcsa_gateway"= "192.168.177.1",
			"vcsa_netmask"= "255.255.255.0"
		}
	}
	 
`, dc, serverType, networkProfileWan, networkProfileLan)
}
