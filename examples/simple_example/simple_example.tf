/* Simple example of using metalcloud */
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
       version = ">= 2.2.7"
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
   
    infrastructure_label = "test-infra7"
    datacenter_name = "${var.datacenter}" 

    create_if_not_exists = true
}

data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "esxi-700-uefi-v2"
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

data "metalcloud_server_type" "large"{
  server_type_name = "M.16.16.1.v3"
}


resource "metalcloud_instance_array" "cluster" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    instance_array_label = "test-3"

    instance_array_instance_count = 1 //deprecated, keep equal to 1

    instance_server_type{
      instance_index=0
      server_type_id=data.metalcloud_server_type.large.server_type_id
    }

    volume_template_id = tonumber(data.metalcloud_volume_template.esxi7.id)

    instance_array_firewall_managed = false

    interface{
      interface_index = 0
      network_id = metalcloud_network.data.id
    }

    interface{
      interface_index = 1
      network_id = metalcloud_network.storage.id
    }

    instance_custom_variables {
      instance_index = 0
      custom_variables={
        "test1":"test2"
        "test3":"test4"
      }
    }

}

resource "metalcloud_shared_drive" "datastore" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  
    shared_drive_label = "test-da-1"
    shared_drive_size_mbytes = 40966
    shared_drive_storage_type = "iscsi_ssd"

    shared_drive_attached_instance_arrays = [metalcloud_instance_array.cluster.instance_array_id]  //this will create a dependency on the instance array
}


# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = true

  # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  # instead of exiting while the deploy is ongoing

  await_deploy_finished = true

  # This option disables a safety check that metalsoft performs to prevent accidental data loss
  # It is required when testing delete operations

  allow_data_loss = true

  # IMPORTANT. This is important to ensure that deploys happen after everything else. If you need to add or remove resources dynamically
  # use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
  depends_on = [
    metalcloud_instance_array.cluster,
    metalcloud_shared_drive.datastore
  ]

}
