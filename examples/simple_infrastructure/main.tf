# Configure the metalcloud provider
provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key 
   endpoint = var.endpoint
}

# Identity the ID of the volume template we want
data "metalcloud_volume_template" "centos76" {
  volume_template_label = "centos7-6"
}


resource "metalcloud_infrastructure" "my-infra92" {
  
  infrastructure_label = "my-terraform-infra92"
  datacenter_name = var.datacenter

  # Remove this to actually deploy changes, otherwise all changes will remain in edit mode only.
  prevent_deploy = true 

  shared_drive {
    shared_drive_label = "my-shared-drive"
    shared_drive_size_mbytes = 40965
    shared_drive_storage_type = "iscsi_ssd"
    shared_drive_attached_instance_arrays = ["web-servers-2"]
  }

  instance_array {
    # Name of your cluster. Needs to obey DNS rules as it will translate into a DNS record.
    instance_array_label = "web-servers"

    instance_array_instance_count = 1
    instance_array_ram_gbytes = 2
    instance_array_processor_count = 1
    instance_array_processor_core_count = 2

    drive_array{
      drive_array_label = "web-servers-centos"
      drive_array_storage_type = "iscsi_hdd"

      # The size of the drive array in MBytes
      drive_size_mbytes_default = 49000

      # The id of the template we located earlier
      volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
    }

    #one or more FW rules. By default all traffic is denied so we need at least one entry.
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
    # Name of your cluster. Needs to obey DNS rules as it will translate into a DNS record.
    instance_array_label = "web-servers-2"

    instance_array_instance_count = 1
    instance_array_ram_gbytes = 2
    instance_array_processor_count = 1
    instance_array_processor_core_count = 2
    instance_array_additional_wan_ipv4_json = jsonencode(
                {
                   configs = [
                       {
                           forced_subnet_pool_id = 1
                           override_vlan_id      = 100
                        },
                    ]
                }
            )

    drive_array{
      drive_array_label = "web-servers-centos-2"
      drive_array_storage_type = "iscsi_hdd"

      # The size of the drive array in MBytes
      drive_size_mbytes_default = 49000

      # The id of the template we located earlier
      volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
    }

    #one or more FW rules. By default all traffic is denied so we need at least one entry.
    firewall_rule {
      firewall_rule_description = "test fw rule-2"
      firewall_rule_port_range_start = 22
      firewall_rule_port_range_end = 22
      firewall_rule_source_ip_address_range_start="0.0.0.0"
      firewall_rule_source_ip_address_range_end="0.0.0.0"
      firewall_rule_protocol="tcp"
      firewall_rule_ip_address_type="ipv4"
    }
  }

  
 
}
