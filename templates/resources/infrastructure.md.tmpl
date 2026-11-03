---
page_title: "metalcloud_infrastructure Resource - terraform-provider-metalcloud"
description: |-
  Infrastructure resource for creating and deleting a MetalCloud infrastructure
---

# metalcloud_infrastructure (Resource)

Creates an infrastructure in a site. An infrastructure is the container for every other deployable resource: server and VM instance groups, drives, logical networks, extension instances and endpoint instance groups.

Only the label is updated in place. The deploy options on this resource are used **when the resource is destroyed**: the infrastructure is marked for deletion and, unless `prevent_deploy` is `true`, a deploy is triggered to tear it down. To deploy changes to the infrastructure's contents use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md).

If the infrastructure may already exist, prefer the [`metalcloud_infrastructure`](../data-sources/infrastructure.md) data source with `create_if_missing = true`.

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
- `await_deploy_finish` (Boolean) Wait for the deploy to finish before returning. Times out after 30 minutes. Defaults to `false`.
- `force_shutdown` (Boolean) Force an immediate shutdown of all servers in the infrastructure, skipping the soft shutdown sequence. Defaults to `false`.
- `prevent_deploy` (Boolean) When `true` no deploy is triggered. Defaults to `true`.
- `soft_shutdown_timeout_seconds` (Number) Seconds to wait for the soft shutdown before attempting a hard shutdown. Defaults to `180`.

### Read-Only

- `infrastructure_id` (String) Id of the infrastructure.

## Import

```sh
terraform import metalcloud_infrastructure.infra <infrastructure_id>
```
