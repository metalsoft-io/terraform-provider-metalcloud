package metalcloud

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccInfrastructureResourceFixture1(infrastructureLabel string, instanceArray1Count int, instanceArray2Count int) string {

	datacenter := os.Getenv("METALCLOUD_DATACENTER")
	apiKey := os.Getenv("METALCLOUD_API_KEY")
	user := os.Getenv("METALCLOUD_USER_EMAIL")
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")

	return fmt.Sprintf(
		`
		data "metalcloud_volume_template" "centos76" {
			volume_template_label = "centos7-6"
		}

		resource "metalcloud_infrastructure" "foo" {

			infrastructure_label = "my-terraform-infra-%s"
			datacenter_name = "%s"	

			prevent_deploy = true

			network{
			  network_type = "san"
			  network_label = "san"
			}
		  
			network{
			  network_type = "wan"
			  network_label = "internet"
			}
		  
			network{
			  network_type = "lan"
			  network_label = "private"
			}
		  
		  
			instance_array {
				  instance_array_label = "master"
				  instance_array_instance_count = %d
				  interface{
					  interface_index = 0
					  network_label = "san"
				  }
		  
				  interface{
					  interface_index = 1
					  network_label = "internet"
				  }
		  
				  interface{
					  interface_index = 2
					  network_label = "private"
				  }
				  
				  drive_array{
					drive_array_label = "testia2-centos"
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
				  }
		  
				  firewall_rule {
							  firewall_rule_description = "test fw rule"
							  firewall_rule_port_range_start = 22
							  firewall_rule_port_range_end = 22
							  firewall_rule_source_ip_address_range_start="0.0.0.0"
							  firewall_rule_source_ip_address_range_end="0.0.0.0"
							  firewall_rule_protocol="tcp"
							  firewall_rule_ip_address_type="ipv4"
						  }
			}
		  
			instance_array {
				  instance_array_label = "slave"  
				  instance_array_instance_count = %d
		  
				  drive_array{
					drive_array_label="asd2-centos"
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
				  }
		  
						  firewall_rule {
							  firewall_rule_description = "test fw rule"
							  firewall_rule_port_range_start = 22
							  firewall_rule_port_range_end = 22
							  firewall_rule_source_ip_address_range_start="0.0.0.0"
							  firewall_rule_source_ip_address_range_end="0.0.0.0"
							  firewall_rule_protocol="tcp"
							  firewall_rule_ip_address_type="ipv4"
						  }
			}
		}
		provider "metalcloud" {
			user_email = "%s"
			api_key = "%s"
			endpoint = "%s"
			}
		`,
		infrastructureLabel,
		datacenter,
		instanceArray1Count,
		instanceArray2Count,
		user,
		apiKey,
		endpoint)
}

func testAccInfrastructureResourceFixture2(infrastructureLabel string, instanceArray1Count int, instanceArray2Count int) string {

	datacenter := os.Getenv("METALCLOUD_DATACENTER")
	apiKey := os.Getenv("METALCLOUD_API_KEY")
	user := os.Getenv("METALCLOUD_USER_EMAIL")
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")

	return fmt.Sprintf(
		`
		data "metalcloud_volume_template" "centos76" {
			volume_template_label = "centos7-6"
		}

		resource "metalcloud_infrastructure" "foo" {

			infrastructure_label = "my-terraform-infra-%s"
			datacenter_name = "%s"	
			
			prevent_deploy = true

			network{
			  network_type = "san"
			  network_label = "san"
			}
		  
			network{
			  network_type = "wan"
			  network_label = "internet"
			}
		  
			network{
			  network_type = "lan"
			  network_label = "private"
			}
		  
		  
			instance_array {
				  instance_array_label = "master"
				  instance_array_instance_count = %d
				  interface{
					  interface_index = 0
					  network_label = "san"
				  }
		  
				  interface{
					  interface_index = 1
					  network_label = "internet"
				  }
		  
				  interface{
					  interface_index = 2
					  network_label = "private"
				  }
				  
				  drive_array{
					drive_array_label = "testia2-centos"
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
				  }
		  
				  firewall_rule {
							  firewall_rule_description = "test fw rule"
							  firewall_rule_port_range_start = 22
							  firewall_rule_port_range_end = 22
							  firewall_rule_source_ip_address_range_start="0.0.0.0"
							  firewall_rule_source_ip_address_range_end="0.0.0.0"
							  firewall_rule_protocol="tcp"
							  firewall_rule_ip_address_type="ipv4"
						  }
			}
		  
			instance_array {
				  instance_array_label = "slave"  
				  instance_array_instance_count = %d
		  
				  drive_array{
					drive_array_label="asd2-centos"
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
				  }
		  
						  firewall_rule {
							  firewall_rule_description = "test fw rule"
							  firewall_rule_port_range_start = 22
							  firewall_rule_port_range_end = 22
							  firewall_rule_source_ip_address_range_start="0.0.0.0"
							  firewall_rule_source_ip_address_range_end="0.0.0.0"
							  firewall_rule_protocol="tcp"
							  firewall_rule_ip_address_type="ipv4"
						  }
			}

		}
		provider "metalcloud" {
			user_email = "%s"
			api_key = "%s"
			endpoint = "%s"
			}
		`,
		infrastructureLabel,
		datacenter,
		instanceArray1Count,
		instanceArray2Count,
		user,
		apiKey,
		endpoint)
}

func TestAccInfrastructureResource_basic(t *testing.T) {

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	///////at create
	expectedIAsAfterCreate := []interface{}{
		map[string]interface{}{
			"instance_array_label":          "master",
			"instance_array_instance_count": 1,
		},
		map[string]interface{}{
			"instance_array_label":          "slave",
			"instance_array_instance_count": 1,
		},
	}

	///////after update
	expectedIAsAfterUpdate := []interface{}{
		map[string]interface{}{
			"instance_array_label":          "master",
			"instance_array_instance_count": 1,
		},
		map[string]interface{}{
			"instance_array_label":          "slave",
			"instance_array_instance_count": 2,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccInfrastructureResourceFixture1(rName, 1, 1),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckResourceExists("metalcloud_infrastructure.foo"),
					// verify remote values
					testAccCheckInfrastructureExists("metalcloud_infrastructure.foo"),
					testAccCheckInstanceArray("metalcloud_infrastructure.foo", expectedIAsAfterCreate),
					// verify local values
					resource.TestCheckResourceAttr("metalcloud_infrastructure.foo", "infrastructure_label", "my-terraform-infra-"+rName),
					testAccCheckVolumeTemplate("data.metalcloud_volume_template.centos76"),
				),
			},
			{
				// expand second IA
				Config: testAccInfrastructureResourceFixture1(rName, 1, 2),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckResourceExists("metalcloud_infrastructure.foo"),
					// verify remote values
					testAccCheckInfrastructureExists("metalcloud_infrastructure.foo"),
					testAccCheckInstanceArray("metalcloud_infrastructure.foo", expectedIAsAfterUpdate),

					// verify local values
					resource.TestCheckResourceAttr("metalcloud_infrastructure.foo", "infrastructure_label", "my-terraform-infra-"+rName),
				),
			},
			{
				// shrink second IA back
				Config: testAccInfrastructureResourceFixture1(rName, 1, 1),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckResourceExists("metalcloud_infrastructure.foo"),
					// verify remote values
					testAccCheckInfrastructureExists("metalcloud_infrastructure.foo"),
					testAccCheckInstanceArray("metalcloud_infrastructure.foo", expectedIAsAfterCreate),

					// verify local values
					resource.TestCheckResourceAttr("metalcloud_infrastructure.foo", "infrastructure_label", "my-terraform-infra-"+rName),
				),
			},
		},
	})
}

func testAccCheckVolumeTemplate(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in Terraform: %s", n)
		}

		volumeTemplateID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*metalcloud.Client)

		vt, err := client.VolumeTemplateGet(volumeTemplateID)
		if err != nil || vt.VolumeTemplateLabel != rs.Primary.Attributes["volume_template_label"] {
			return fmt.Errorf("Volume template data could not be verified returned vt=%v", vt)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckResourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in Terraform: %s", n)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckInfrastructureExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in Terraform: %s", n)
		}

		infraID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*metalcloud.Client)

		infra, err := client.InfrastructureGet(infraID)

		if infra.InfrastructureID != infraID {
			return fmt.Errorf("infrastructure id not correct %d", infraID)
		}

		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckInstanceArray(n string, expectedIAs []interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in Terraform: %s", n)
		}

		infraID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*metalcloud.Client)

		realIAs, err := client.InstanceArrays(infraID)
		if err != nil {
			return err
		}

		for _, e := range expectedIAs {
			expectedIA := e.(map[string]interface{})
			verified := false
			for _, r := range *realIAs {
				if expectedIA["instance_array_label"] == r.InstanceArrayLabel {

					if expectedIA["instance_array_instance_count"] != r.InstanceArrayInstanceCount {
						return fmt.Errorf("%s instance array's instance_array_instance_count filed is wrong: %d and expected %d", r.InstanceArrayLabel, r.InstanceArrayInstanceCount, expectedIA["instance_array_instance_count"])
					}

					verified = true
				}
			}
			if !verified {
				return fmt.Errorf("Instance array expected %s but was not found", expectedIA["instance_array_label"])
			}
		}

		return nil
	}
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("METALCLOUD_USER_EMAIL"); v == "" {
		t.Fatal("METALCLOUD_USER_EMAIL must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_API_KEY"); v == "" {
		t.Fatal("METALCLOUD_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_ENDPOINT"); v == "" {
		t.Fatal("METALCLOUD_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("METALCLOUD_DATACENTER"); v == "" {
		t.Fatal("METALCLOUD_DATACENTER must be set for acceptance tests")
	}
}
