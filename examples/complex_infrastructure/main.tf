provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key 
   endpoint = var.endpoint
}

data "metalcloud_volume_template" "centos76" {
  volume_template_label = "centos7-6"
}

resource "metalcloud_infrastructure" "my-infra216" {
  
  infrastructure_label = "my-terraform-infra216"
  datacenter_name = var.datacenter
  
  //remove this to actually deploy changes, otherwise all changes will remain in edit mode only.
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
    instance_array_label = "exmaple-master"
    instance_array_instance_count = 2
    instance_array_ram_gbytes = 8
    instance_array_processor_count = 1
    instance_array_processor_core_count = 8

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
      drive_array_label = "example-master-os-drive"
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
    instance_array_label = "example-slave"  
    instance_array_instance_count = 1
    instance_array_ram_gbytes = 8
    instance_array_processor_count = 1
    instance_array_processor_core_count = 8

    drive_array{
      drive_array_label = "example-slave-os-drive"
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