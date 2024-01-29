---
layout: "metalcloud"
page_title: "Metalcloud: Getting started"
description: |-
  Getting started guide with the metalcloud provider.
---

# Metalcloud: Getting started

The Metalcloud provider allows users to provision bare metal resources such as physical servers, switch configurations, iSCSI drives etc.


## Provisioning a server

To provision a server:

```hcl
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
   
    infrastructure_label = "test-infra"
    datacenter_name = var.datacenter

    create_if_not_exists = true
}

data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "esxi-700-uefi-v2"
}

resource "metalcloud_network" "data" {
  network_label = "data-network"
  network_type = "wan"
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
}

resource "metalcloud_network" "lan" {
  network_label = "lan-network"
  network_type = "lan" #wan, lan or san. note that san is not available in all environments
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
}

data "metalcloud_server_type" "large"{
  server_type_name = "M.16.16.1.v3" #needs to match an existing server type in your environment
}

resource "metalcloud_instance_array" "cluster" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    instance_array_label = "test-3"

    instance_array_instance_count = 1 //deprecated, keep equal to 1
    instance_array_boot_method = "local_drives"

    instance_server_type{
      instance_index=0
      server_type_id=data.metalcloud_server_type.large.server_type_id
    }

    volume_template_id = tonumber(data.metalcloud_volume_template.esxi7.id)

    instance_array_firewall_managed = false

    interface{
      interface_index = 0
      network_id = metalcloud_network.lan.id
    }

    interface{
      interface_index = 1
      network_id = metalcloud_network.data.id
    }

    instance_custom_variables {
      instance_index = 1
      custom_variables={
        "test1":"test2"
        "test3":"test4"
      }
    }

    depends_on = [
      metalcloud_network.data,
      metalcloud_network.lan,
    ]

}

# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = false

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
    metalcloud_network.data,
    metalcloud_network.lan,
  ]

}
```

## Authentication

Getting the API Key is typically done via the  Metal Cloud's API key section. Use it with a -var or as an env variable:

```bash
export TF_VAR_api_key="<yourkey>"
export TF_VAR_user_email="test@test.com"
export TF_VAR_endpoint="https://api.poc.metalsoft.io"
export TF_VAR_datacenter="uk-reading"

terraform plan
```

!> Warning: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system. 

