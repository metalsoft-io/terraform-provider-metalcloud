
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

   logging = true

}

# This is an infrastructure reference. It is needed to avoid a cyclic dependency where the 
# infrastructure depends on the resources and vice-versa. This will create the infrastructure if it does not exist
# if the create_if_not_exists flag is set to true
data "metalcloud_infrastructure" "infra" {

    infrastructure_label = "test-infra"
    datacenter_name = "${var.datacenter}" 

    create_if_not_exists = true
}

data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "esxi-700-uefi-v2"
}

resource "metalcloud_instance_array" "cluster" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    instance_array_label = "test-3"

    instance_array_instance_count = 1 //deprecated, keep equal to 1
    instance_array_ram_gbytes = "16"
    instance_array_processor_count = 1
    instance_array_processor_core_count = 1
    instance_array_boot_method = "local_drives"

    volume_template_id = tonumber(data.metalcloud_volume_template.esxi7.id)

    instance_array_firewall_managed = false
/*
    interface{
      interface_index = 0
      network_label = "storage-network"
    }

    interface{
      interface_index = 1
      network_label = "data-network"
    }
*/

}

data "metalcloud_server_type" "large"{
     server_type_name = "M.12.8.1"
}


resource "metalcloud_kubernetes" "kube01" {
    infrastructure_id =  data.metalcloud_infrastructure.infra.infrastructure_id

    cluster_label = "testkube"
    instance_array_instance_count_master = 2
    instance_server_type_master {
        instance_index = 0
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

    instance_server_type_master {
        instance_index = 1
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }
    
    instance_array_instance_count_worker = 3
    instance_server_type_worker {
        instance_index = 0
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

    instance_server_type_worker {
        instance_index = 1
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

     instance_server_type_worker {
        instance_index = 2
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

}


# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = true

  #these options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  #instead of exiting while the deploy is ongoing

  await_deploy_finished = true

  #this option disables a safety check that MetalSoft performs to prevent accidental data loss
  #it is required when testing delete operations

  allow_data_loss = true

  # IMPORTANT. This is important to ensure that deploys happen after everything else. If you need to add or remove resources dynamically
  # use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
  depends_on = [
    metalcloud_instance_array.cluster,
    metalcloud_kubernetes.kube01
  ]
}

data "metalcloud_infrastructure_output" "output"{
  infrastructure_id = metalcloud_infrastructure_deployer.infrastructure_deployer.infrastructure_id
  depends_on = [
      metalcloud_infrastructure_deployer.infrastructure_deployer
    ]
}

output "cluster_credentials" {
    value = data.metalcloud_infrastructure_output.output.clusters
}