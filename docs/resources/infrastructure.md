---
page_title: "metalcloud_infrastructure Resource - terraform-provider-metalcloud"
description: |-
  Infrastructure resource for creating and deleting a MetalCloud infrastructure
---

# metalcloud_infrastructure (Resource)

Creates an infrastructure in a site. An infrastructure is the container for every other deployable resource: server and VM instance groups, drives, logical networks, extension instances and endpoint instance groups.

Only the label is updated in place. The deploy options on this resource are used **when the resource is destroyed**: the infrastructure is marked for deletion and, unless `prevent_deploy` is `true`, a deploy is triggered to tear it down. To deploy changes to the infrastructure's contents use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md).

To manage an infrastructure that already exists, [import it](#import) instead of creating a new one.

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = var.site
}

resource "metalcloud_infrastructure" "infra" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "my-infra01"

  # Deploy options used when the infrastructure is destroyed.
  prevent_deploy      = false
  await_deploy_finish = true
  allow_data_loss     = true
}
```

## Schema

### Required

- `label` (String) Infrastructure label. Updated in place.
- `site_id` (String) Id of the site the infrastructure is created in. Changing this value forces a new resource.

### Optional

- `allow_data_loss` (Boolean) Allow the deploy to perform operations that may cause data loss. Defaults to `false`.
- `attempt_hard_shutdown` (Boolean) Force a hard power off of servers that are still running when the soft shutdown timeout expires. Defaults to `true`.
- `attempt_soft_shutdown` (Boolean) Attempt a soft (ACPI) power off of the servers before the deploy. Defaults to `true`.
- `await_deploy_finish` (Boolean) Wait for the deploy to finish before returning. Times out after 2 hours. Defaults to `false`.
- `force_shutdown` (Boolean) Force an immediate shutdown of all servers in the infrastructure, skipping the soft shutdown sequence. Defaults to `false`.
- `prevent_deploy` (Boolean) When `true` no deploy is triggered. Defaults to `true`.
- `soft_shutdown_timeout_seconds` (Number) Seconds to wait for the soft shutdown before attempting a hard shutdown. Defaults to `180`.

### Read-Only

- `infrastructure_id` (String) Id of the infrastructure.

## Import

An existing infrastructure can be brought under Terraform by importing it with its Id. Importing does not change the infrastructure in MetalCloud.

### Find the infrastructure Id

The Id is shown on the infrastructure's page in the MetalCloud UI. It can also be listed with the MetalCloud CLI:

```sh
metalcloud-cli infrastructure list
```

### Import the infrastructure

1. Add a `metalcloud_infrastructure` resource whose `site_id` and `label` match the existing infrastructure:

   ```hcl
   resource "metalcloud_infrastructure" "infra" {
     site_id = data.metalcloud_site.dc.site_id
     label   = "my-infra01"
   }
   ```

2. Import it, either with the CLI:

   ```sh
   terraform import metalcloud_infrastructure.infra 1234
   ```

   or, with Terraform 1.5 or later, with an `import` block that `terraform apply` processes:

   ```hcl
   import {
     to = metalcloud_infrastructure.infra
     id = "1234"
   }
   ```

3. Run `terraform plan` and check it before applying:
   * A different `site_id` **replaces** the infrastructure: Terraform plans to delete it and create a new one. Correct `site_id` in the configuration before applying.
   * A different `label` renames the infrastructure in place.
   * The deploy options (`prevent_deploy`, `allow_data_loss` and the others) exist only in Terraform. They are imported with their defaults, so values set in the configuration show as an in-place update that changes nothing in MetalCloud.

Once imported, the infrastructure is managed like one created by Terraform: `terraform destroy` deletes it, and tears it down when `prevent_deploy = false`.

### Migrate from `create_if_missing`

Configurations that use the [`metalcloud_infrastructure`](../data-sources/infrastructure.md) data source with `create_if_missing = true` can switch to this resource without recreating the infrastructure:

1. Replace the `data "metalcloud_infrastructure"` block with a `resource "metalcloud_infrastructure"` block with the same `site_id` and `label`.
2. Change every `data.metalcloud_infrastructure.<name>.infrastructure_id` reference to `metalcloud_infrastructure.<name>.infrastructure_id`.
3. Import the infrastructure as described above, then run `terraform plan`. It should show no changes to the resources inside the infrastructure: they are already in state and keep the same `infrastructure_id`.

!> Import the infrastructure **before** running `terraform apply` with the new configuration (an `import` block does both in one step). Otherwise Terraform creates a second infrastructure and, because changing `infrastructure_id` forces replacement, deletes and recreates every resource inside the existing one in the new infrastructure.

### Resources inside the infrastructure

The infrastructure import covers only the infrastructure itself. The resources inside it are imported separately:

* `metalcloud_server_instance_group` by `<server_instance_group_id>`
* `metalcloud_logical_network` by `<logical_network_id>`
* `metalcloud_endpoint_instance_group` by `<endpoint_instance_group_id>`
* `metalcloud_extension_instance` by `<extension_instance_id>`
* `metalcloud_drive` by `<infrastructure_id>:<drive_id>`
* `metalcloud_vm_instance_group` by `<infrastructure_id>:<vm_instance_group_id>`
