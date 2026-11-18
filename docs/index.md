---
page_title: "metalcloud Provider"
description: |-
  The MetalCloud provider enables control over the MetalCloud's resources using Terraform.
---

# MetalCloud Provider

The MetalCloud provider provides options to provision bare metal servers, switches, drives, etc., via Terraform.

Provider version 8.x is built for the MetalCloud v8.0 API. Pin the provider version to match your MetalCloud controller, as shown in the example below.

## Argument Reference

The following arguments are supported:

* `endpoint` - (Required) Base URL of the MetalCloud controller, for example `https://metal.mycompany.com`. Do not add a path: the provider adds `/api/v2/...` itself.
* `api_key` - (Required) The user's API key.
* `user_email` - (Optional) The user's email address.
* `logging` - (Optional) Set to `"true"` to enable debug logging of API requests.
* `insecure` - (Optional, Boolean) Allow insecure TLS connections. Defaults to `false`.
* `timeout` - (Optional, Number) Timeout of each API request, in seconds. Defaults to `300`. This does not limit how long the provider waits for a deploy.

## Configuration

The example below reads its settings from input variables, which can be set with environment variables or `-var`:

```shell
export TF_VAR_endpoint="https://metal.mycompany.com"
export TF_VAR_api_key="<your-api-key>"
export TF_VAR_user_email="user@mycompany.com"
export TF_VAR_site="uk-reading"

terraform init
terraform plan
terraform apply
```

## Deploys and destroy

Creating, changing or deleting resources inside an infrastructure only changes its pending configuration. Nothing changes on the servers until the infrastructure is **deployed**:

* [`metalcloud_infrastructure_deployer`](resources/infrastructure_deployer.md) deploys the infrastructure on every `terraform apply` when `prevent_deploy = false`. Make it `depends_on` the other resources so the deploy runs after all of them. With `prevent_deploy = true`, `terraform apply` only records the pending changes.
* With `await_deploy_finish = true`, `terraform apply` waits for the deploy to finish, for up to 2 hours, and fails if any deploy job fails. Otherwise it returns as soon as the deploy has started.
* `terraform destroy` deletes the infrastructure through [`metalcloud_infrastructure`](resources/infrastructure.md). The servers are only torn down if that resource has `prevent_deploy = false`; otherwise the infrastructure is only marked for deletion and stays active until it is deployed in MetalCloud. `prevent_deploy` is read from state during destroy, so run `terraform apply` after changing it and before `terraform destroy`. Destroying `metalcloud_infrastructure_deployer` does nothing.
* `allow_data_loss = true` is required for deploys that delete servers or drives, including the one triggered by `terraform destroy`.

## Example Usage

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
