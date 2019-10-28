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
	user := os.Getenv("METALCLOUD_USER")
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")

	return fmt.Sprintf(
		`
		data "metalcloud_volume_template" "centos76" {
			volume_template_label = "centos7-6"
		}

		resource "metalcloud_infrastructure" "foo" {

			infrastructure_label = "my-terraform-infra-%s"
			datacenter_name = "%s"	

			instance_array {
				instance_array_label = "as111"
				instance_array_instance_count = %d

				firewall_rule {
					firewall_rule_description = "test fw rule"
					firewall_rule_port_range_start = 22
					firewall_rule_port_range_end = 22
					firewall_rule_source_ip_address_range_start="0.0.0.0"
					firewall_rule_source_ip_address_range_end="0.0.0.0"
					firewall_rule_protocol="tcp"
					firewall_rule_ip_address_type="ipv4"
				}

				drive_array{
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = data.metalcloud_volume_template.centos76.id
				}
			}

			instance_array {
				instance_array_label = "asd2"  
				instance_array_instance_count = %d
				drive_array{
					drive_array_storage_type = "iscsi_hdd"
					drive_size_mbytes_default = 49000
					volume_template_id = data.metalcloud_volume_template.centos76.id
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
			user = "%s"
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
	var expectedIAs1 []interface{}

	var ia1 = make(map[string]interface{})
	ia1["instance_array_label"] = "as111"
	ia1["instance_array_instance_count"] = 1
	ia1["instance_array_subdomain"] = "test"

	expectedIAs1 = append(expectedIAs1, ia1)

	var ia2 = make(map[string]interface{})
	ia2["instance_array_label"] = "asd2"
	ia2["instance_array_instance_count"] = 1
	ia2["instance_array_subdomain"] = "test"

	expectedIAs1 = append(expectedIAs1, ia2)
	///////after update
	var expectedIAs2 []interface{}

	expectedIAs2 = append(expectedIAs2, ia1)

	var ia2Modif = make(map[string]interface{})
	ia2Modif["instance_array_label"] = "asd2"
	ia2Modif["instance_array_instance_count"] = 2

	expectedIAs2 = append(expectedIAs2, ia2Modif)

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
					testAccCheckInstanceArray("metalcloud_infrastructure.foo", expectedIAs1),
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
					testAccCheckInstanceArray("metalcloud_infrastructure.foo", expectedIAs2),

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

		client := testAccProvider.Meta().(*metalcloud.MetalCloudClient)

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

		client := testAccProvider.Meta().(*metalcloud.MetalCloudClient)

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

		client := testAccProvider.Meta().(*metalcloud.MetalCloudClient)

		realIAs, err := client.InstanceArrays(infraID)
		if err != nil {
			return err
		}

		for _, e := range expectedIAs {
			expectedIA := e.(map[string]interface{})
			var verified = false
			for _, r := range *realIAs {
				if expectedIA["instance_array_label"] == r.InstanceArrayLabel &&
					expectedIA["instance_array_instance_count"] == r.InstanceArrayInstanceCount {
					verified = true
					break
				}
			}
			if !verified {
				return fmt.Errorf("Instance array with label %s not provisioned correctly", expectedIA["instance_array_label"])
			}
		}

		return nil
	}
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
	if v := os.Getenv("METALCLOUD_DATACENTER"); v == "" {
		t.Fatal("METALCLOUD_DATACENTER must be set for acceptance tests")
	}
}
