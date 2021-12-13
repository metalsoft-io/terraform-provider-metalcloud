terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}


data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "${var.tenancy_config.customer_name}" 
    datacenter_name = "${var.tenancy_config.datacenter}" 

}

resource "metalcloud_network" "data" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
    network_label = "data-network"
    network_type = "wan"
}

resource "metalcloud_network" "storage" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
    network_label = "storage-network"
    network_type = "san"
}

resource "metalcloud_network_profile" "profile" {
    network_profile_label = "network-profile-${var.tenancy_config.customer_name}"
    datacenter_name = "${var.tenancy_config.datacenter}" 
    network_type = "wan"

    network_profile_vlan {
      vlan_id = "${var.tenancy_config.esxi_vlan_id}"
      port_mode = "trunk"
      provision_subnet_gateways = false
    }
}


module "tenancy_cluster" {
  source = "./tenancy_cluster"
  
  count = length(var.tenancy_config.clusters)
  
  clustername = var.tenancy_config.clusters[count.index].clustername
  customer_prefix = var.tenancy_config.clusters[count.index].customer_prefix
  compute_nodes = var.tenancy_config.clusters[count.index].compute_nodes
  datastores= var.tenancy_config.clusters[count.index].datastores

  instance_array_ram_gbytes = var.tenancy_config.clusters[count.index].instance_array_ram_gbytes
  instance_array_processor_count = var.tenancy_config.clusters[count.index].instance_array_processor_count
  instance_array_processor_core_count = var.tenancy_config.clusters[count.index].instance_array_processor_core_count

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  datacenter_name = var.tenancy_config.datacenter
  
  wan_network_id = metalcloud_network.data.id
  san_network_id = metalcloud_network.storage.id
  wan_network_profile_id = metalcloud_network_profile.profile.id
}

resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to trigger deploys.
  prevent_deploy = true

  #these options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  #instead of exiting while the deploy is ongoing

  await_deploy_finished = true
  await_delete_finished = true

  #this option disables a safety check that metalsoft performs to prevent accidental data loss
  #it is required when testing delete operations

  allow_data_loss = true

  depends_on = [
    module.tenancy_cluster
  ]


}