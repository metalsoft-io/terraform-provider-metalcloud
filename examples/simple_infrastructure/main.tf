provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key 
   endpoint = var.endpoint
}

data "metalcloud_volume_template" "centos76" {
  volume_template_label = "centos7-6"
}

resource "metalcloud_infrastructure" "my-infra102" {
  
  infrastructure_label = "my-terraform-infra102"
  datacenter_name = var.datacenter

  
  instance_array {
        instance_array_label = "testia"
        instance_array_instance_count = 1
        instance_array_ram_gbytes = 8
        instance_array_processor_count = 1
        instance_array_processor_core_count = 8

        drive_array{
          drive_array_label = "testia-os-drive"
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
