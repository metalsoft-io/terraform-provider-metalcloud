
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}


# ############################################
# define local vars
# ############################################

locals{
  BOOT_METHOD = "local_drives"

  ISCSI_LUN_TYPE = "iscsi_hdd"
  instance_array_instance_count = "1"
}



data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "esxi-700-uefi-v2"
}

resource "metalcloud_network" "data" {
    infrastructure_id = var.infrastructure_id
    network_label = "data-network"
    network_type = "wan"
}

resource "metalcloud_network" "storage" {
    infrastructure_id = var.infrastructure_id
    network_label = "storage-network"
    network_type = "san"
}

resource "metalcloud_instance_array" "cluster" {

    #this will create a series of instances
    count =  length(var.compute_nodes)

    infrastructure_id = var.infrastructure_id

    instance_array_label = "${var.compute_nodes[count.index].compute_node_name}"

    instance_array_instance_count = local.instance_array_instance_count
    instance_array_ram_gbytes = "${var.instance_array_ram_gbytes}"
    instance_array_processor_count = "${var.instance_array_processor_count}"
    instance_array_processor_core_count = "${var.instance_array_processor_core_count}"
    instance_array_boot_method = local.BOOT_METHOD

    volume_template_id = tonumber(data.metalcloud_volume_template.esxi7.id)

    instance_array_firewall_managed = false

  interface{     
    interface_index = 0  
    network_id = metalcloud_network.data.id
  }

  interface{     
    interface_index = 1
    network_id = metalcloud_network.data.id
  }

  interface{     
    interface_index = 2 
    network_id = metalcloud_network.storage.id
  }

  interface{     
    interface_index = 3 
    network_id = metalcloud_network.storage.id
  }

  network_profile {
    network_id = metalcloud_network.data.id
    network_profile_id = metalcloud_network_profile.profile.id
  }

  instance_array_custom_variables = {
    mgmt_vlan_id = 500
  }

  depends_on = [
    metalcloud_network.data,
    metalcloud_network.storage
  ]
}


resource "metalcloud_drive_array" "drives" {

    count =  length(var.compute_nodes)

    infrastructure_id = var.infrastructure_id

    drive_array_label = "${metalcloud_instance_array.cluster[count.index].instance_array_label}-da"

    //to which instance array is this drive array attached
    instance_array_id = metalcloud_instance_array.cluster[count.index].instance_array_id
    
    drive_array_storage_type = "iscsi_ssd"
    drive_size_mbytes_default = 40960    
}

resource "metalcloud_shared_drive" "datastore" {

    count =  length(var.datastores)

    infrastructure_id = var.infrastructure_id
  
    shared_drive_label = "${var.datastores[count.index].datastore_name}"
    shared_drive_size_mbytes = 40966

    shared_drive_storage_type = "iscsi_hdd"
   

    shared_drive_attached_instance_arrays = metalcloud_instance_array.cluster[*].instance_array_id
}

resource "metalcloud_network_profile" "profile" {
    network_profile_label = "profile-network-${var.clustername}-${var.infrastructure_id}"
    datacenter_name = var.datacenter_name
    network_type = "wan"

    network_profile_vlan {
      vlan_id = var.esxi_vlan_id
      port_mode = "trunk"
      provision_subnet_gateways = false
    }
}