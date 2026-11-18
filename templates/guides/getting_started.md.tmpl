---
page_title: "MetalCloud: Getting started"
description: |-
  Getting started guide with the metalcloud provider.
---

# MetalCloud: Getting started

The MetalCloud provider allows users to provision bare metal resources such as physical servers, switch configurations, iSCSI drives etc.

## Provisioning a server

The configuration below creates an infrastructure with one logical network and one server, and deploys it on `terraform apply`. `terraform destroy` deletes the infrastructure and deploys the deletion, which tears the server down.

```terraform
terraform {
  required_providers {
    metalcloud = {
      source  = "metalsoft-io/metalcloud"
      version = "~> 8.0"
    }
  }
}

provider "metalcloud" {
  endpoint   = var.endpoint
  user_email = var.user_email
  api_key    = var.api_key
  logging    = var.logging
}

data "metalcloud_site" "dc" {
  label = var.site
}

data "metalcloud_fabric" "wan" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "wan-fabric"
}

data "metalcloud_logical_network_profile" "np01" {
  label     = "np-01"
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
  label   = "my-infra01"

  # Deploy the deletion on terraform destroy, so that the servers are torn down.
  # With true, terraform destroy only marks the infrastructure for deletion.
  prevent_deploy = false

  # Make terraform destroy wait for the teardown deploy to finish.
  await_deploy_finish = true

  # Disables the MetalCloud safety check that prevents accidental data loss.
  # Required for the teardown deploy.
  allow_data_loss = true
}

resource "metalcloud_logical_network" "net1" {
  infrastructure_id          = metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.np01.logical_network_profile_id

  name  = "net01"
  label = "net01"
}

resource "metalcloud_server_instance_group" "inst01" {
  infrastructure_id = metalcloud_infrastructure.infra.infrastructure_id

  name  = "inst01"
  label = "inst01"

  instance_count = 1
  server_type_id = data.metalcloud_server_type.srv1.server_type_id
  os_template_id = data.metalcloud_os_template.os1.os_template_id

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.net1.logical_network_id
      tagged             = true
      access_mode        = "l2"
      mtu                = 1500
    }
  ]

  custom_variables = [
    {
      name  = "key1"
      value = "test1"
    },
    {
      name  = "key2"
      value = "test2"
    }
  ]
}

# Deploys the pending changes of the resources above on every terraform apply.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {
  infrastructure_id = metalcloud_infrastructure.infra.infrastructure_id

  # Set to true to only record pending changes without deploying them.
  prevent_deploy = false

  # Make terraform apply wait for the deploy to finish (up to 2 hours)
  # instead of returning while the deploy is ongoing.
  await_deploy_finish = true

  # Disables the MetalCloud safety check that prevents accidental data loss.
  # Required for deploys that delete servers or drives.
  allow_data_loss = true

  # IMPORTANT: this makes the deploy run after every other change. If you add or remove
  # resources dynamically, use count or for_each in those resources, or move everything
  # dynamic into a module and make this resource depend on the module.
  depends_on = [
    metalcloud_logical_network.net1,
    metalcloud_server_instance_group.inst01,
  ]
}

variable "endpoint" {
  description = "Base URL of the MetalCloud controller, for example https://metal.mycompany.com"
  type        = string
}

variable "api_key" {
  description = "MetalCloud API key"
  type        = string
  sensitive   = true
}

variable "user_email" {
  description = "Email address of the MetalCloud user"
  type        = string
  default     = ""
}

variable "logging" {
  description = "Set to \"true\" to enable debug logging"
  type        = string
  default     = "false"
}

variable "site" {
  description = "Label of the MetalCloud site"
  type        = string
}
```

## Deploys and destroy

Changes to resources inside an infrastructure stay pending until the infrastructure is deployed:

* `metalcloud_infrastructure_deployer` with `prevent_deploy = false` deploys them on every `terraform apply`. Set it to `true` to only record the pending changes.
* `await_deploy_finish = true` makes Terraform wait for the deploy to finish, for up to 2 hours, and fail if a deploy job fails.
* `metalcloud_infrastructure` with `prevent_deploy = false` deploys the deletion on `terraform destroy`. With `true`, the infrastructure is only marked for deletion and stays active. This setting is read from state, so run `terraform apply` after changing it and before `terraform destroy`.

To manage an infrastructure that already exists, [import it](../resources/infrastructure.md#import) into the `metalcloud_infrastructure` resource. To only add resources to an existing infrastructure without Terraform ever deleting it, look it up with the [`metalcloud_infrastructure`](../data-sources/infrastructure.md) data source instead; the deletion of the resources inside it then stays pending until the infrastructure is deployed in MetalCloud. Do not use the data source's deprecated `create_if_missing` option: it creates the infrastructure during `terraform plan`.

See [Deploys and destroy](../index.md#deploys-and-destroy) on the provider page for details.

## Authentication

Getting the API Key is typically done via the MetalCloud API key section. Set `endpoint` to the base URL of the MetalCloud controller, without a path. Pass the values with `-var` or as environment variables:

```bash
export TF_VAR_endpoint="https://api.poc.metalsoft.io"
export TF_VAR_user_email="test@test.com"
export TF_VAR_api_key="<yourkey>"
export TF_VAR_site="uk-reading"

terraform init
terraform plan
terraform apply
```

!> Warning: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.
