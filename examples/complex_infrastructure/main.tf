provider "metalcloud" {
   user = var.user
   api_key = var.api_key 
   endpoint = var.endpoint
}

data "metalcloud_volume_template" "centos76" {
  volume_template_label = "centos7-6"
}

resource "metalcloud_infrastructure" "my-infra97" {
  
  infrastructure_label = "my-terraform-infra97"
  datacenter_name = "us-santaclara"

  instance_array {
        instance_array_label = "testia2"
        instance_array_instance_count = 2

        drive_array{
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
        instance_array_label = "asd2"  
        instance_array_instance_count = 1

        drive_array{
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
