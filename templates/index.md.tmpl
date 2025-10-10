---
page_title: "metalcloud Provider"
description: |-
  The MetalCloud provider enables control over the MetalCloud's resources using Terraform.
---

# MetalCloud Provider

The MetalCloud provider provides options to provision bare metal servers, switches, drives, etc., via Terraform.

## Argument Reference

The following arguments are supported:

* `endpoint` - (Required) The **API endpoint to connect to.
* `api_key` - (Required) The **User's** API_KEY.
* `user_email` - **User's** email address used as the login identity.
* `logging` - Set the logging level.
* `insecure` - (Boolean) Allow insecure connections. Default is false.

## Example Usage

```terraform
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}

provider "metalcloud" {
  endpoint = var.endpoint
  user_email = var.user_email
  api_key = var.api_key
}

data "metalcloud_site" "dc" {
  label = "${var.site}"
}

data "metalcloud_fabric" "wan" {
  site_id = data.metalcloud_site.dc.site_id
  label = "wan-fabric"
}

data "metalcloud_logical_network_profile" "np01" {
  label = "np-01"
  fabric_id = data.metalcloud_fabric.wan.fabric_id
}

data "metalcloud_server_type" "srv1" {
  label = "M.16.64.2"
}

data "metalcloud_os_template" "os1" {
  label = "ubuntu-22-04"
}

resource "metalcloud_infrastructure" "infra" {
    site_id = data.metalcloud_site.dc.site_id
    label = "my-infra01"

    # Set this to false to actually trigger deploy when the infrastructure is deleted.
    prevent_deploy = true

    # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
    # instead of exiting while the deploy is ongoing
    await_deploy_finish = false

    # This option disables a safety check that metalsoft performs to prevent accidental data loss
    # It is required when testing delete operations
    allow_data_loss = true
}

resource "metalcloud_logical_network" "net1" {
  infrastructure_id = metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.np01.logical_network_profile_id

  name = "net01"
  label = "net01"

  depends_on = [
    metalcloud_infrastructure.infra,
  ]
}

resource "metalcloud_server_instance_group" "inst01" {
  infrastructure_id = metalcloud_infrastructure.infra.infrastructure_id

  name = "inst01"
  label = "inst01"

  instance_count = 1
  server_type_id = data.metalcloud_server_type.srv1.server_type_id
  os_template_id = data.metalcloud_os_template.os1.os_template_id

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.net1.logical_network_id
      tagged = true
      access_mode = "l2"
      mtu = 1500
    }
  ]

  custom_variables = [
    {
      name = "key1"
      value = "test1"
    },
    {
      name = "key2"
      value = "test2"
    }
  ]

  depends_on = [
    metalcloud_logical_network.net1,
    metalcloud_infrastructure.infra,
  ]
}

# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = true

  # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  # instead of exiting while the deploy is ongoing
  await_deploy_finish = false

  # This option disables a safety check that metalsoft performs to prevent accidental data loss
  # It is required when testing delete operations
  allow_data_loss = true

  # IMPORTANT. This is important to ensure that deploys happen after everything else. If you need to add or remove resources dynamically
  # use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
  depends_on = [
      metalcloud_server_instance_group.inst01,
  ]
}

variable "endpoint" {
  default =""
}

variable "api_key" {
  default = ""
}

variable "logging" {
  default="false"
}

variable "site" {
  default=""
}
```
