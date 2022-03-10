---
layout: "metalcloud"
page_title: "Provider: Metalcloud"
description: |-
  The Metalcloud provider enables control over the Metalcloud's resources using Terraform.
---

# Metalcloud Provider


The Metalcloud provider provides options to provision bare metal servers, switches, iSCSI drives, control firewall etc, via Terraform.


## Argument Reference

The following arguments are supported:

* `user_email` - (Required) **User's** email address used as the login identity. This will fallback to using METALCLOUD_API_KEY environment variable.
* `api_key` - (Required) The **User's**  API_KEY. Defaults to the METALCLOUD_API_KEY environment variable.
* `endpoint` - (Required) The **API endpoint to connect to. Defaults to METALCLOUD_ENDPOINT.

## Example Usage

```hcl
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

    interface{
      interface_index = 0
      network_label = "storage-network"
    }

    interface{
      interface_index = 1
      network_label = "data-network"
    }

    instance_custom_variables {
      instance_index = 1
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
    shared_drive_storage_type = "iscsi_hdd"

    shared_drive_attached_instance_arrays = [metalcloud_instance_array.cluster.instance_array_id]  //this will create a dependency on the instance array
}

# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = false

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
    metalcloud_shared_drive.datastore
  ]
}
```
