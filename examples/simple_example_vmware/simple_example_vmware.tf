
variable "user_email" {}
variable "api_key" {}
variable "endpoint" {}
variable "datacenter" {}

terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}

provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key
   endpoint = var.endpoint

}

# This is an infrastructure reference. It is needed to avoid a cyclic dependency where the 
# infrastructure depends on the resources and vice-versa. This will create the infrastructure if it does not exist
# if the create_if_not_exists flag is set to true
data "metalcloud_infrastructure" "infra" {

    infrastructure_label = "test-infra-vmware"
    datacenter_name = "${var.datacenter}" 

    create_if_not_exists = true
}

data "metalcloud_server_type" "large"{
     server_type_name = "M.64.512.10"
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
    network_profile_label = "vmware-cluster"
    datacenter_name = var.datacenter
}

data "metalcloud_network_profile" "vmware_lan"{
    network_profile_label = "vmware-cluster-lan"
    datacenter_name = var.datacenter
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

# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = true

  #these options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  #instead of exiting while the deploy is ongoing

  await_deploy_finished = false

  #this option disables a safety check that MetalSoft performs to prevent accidental data loss
  #it is required when testing delete operations

  allow_data_loss = true

  # IMPORTANT. This is important to ensure that deploys happen after everything else. If you need to add or remove resources dynamically
  # use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
  depends_on = [
    metalcloud_vmware_vsphere.VMWareVsphere
  ]
}

data "metalcloud_infrastructure_output" "output"{
    infrastructure_id = data.metalcloud_infrastructure.infra.id
    depends_on = [ resource.metalcloud_infrastructure_deployer.infrastructure_deployer ]
}
output "cluster_credentials" {
    value = jsondecode(data.metalcloud_infrastructure_output.output.clusters)
}
