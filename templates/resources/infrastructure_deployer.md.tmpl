---
page_title: "metalcloud_infrastructure_deployer Resource - terraform-provider-metalcloud"
description: |-
  Triggers deploys of a MetalCloud infrastructure
---

# metalcloud_infrastructure_deployer (Resource)

The infrastructure deployer has no API object of its own. It triggers a **deploy** of the given infrastructure on create and on every update, applying all pending changes made by the other resources in the same infrastructure. Destroying it does nothing.

Because Terraform applies resources in dependency order, list every other resource of the infrastructure in `depends_on` so the deploy runs last. When resources are created with `count` or `for_each`, or come from a module, depend on the whole collection or module.

## Design vs. deploy

- `prevent_deploy = true` (the default) records design changes without touching hardware. Use it while iterating on a configuration.
- `prevent_deploy = false` triggers the deploy. With `await_deploy_finish = true` the apply blocks until the deploy finishes or 30 minutes pass.

## Example Usage

```hcl
data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "my-infra01"
  create_if_missing = true
}

resource "metalcloud_server_instance_group" "inst01" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  label             = "inst01"
  instance_count    = 1
  server_type_id    = data.metalcloud_server_type.srv1.server_type_id
  os_template_id    = data.metalcloud_os_template.os1.os_template_id
}

resource "metalcloud_infrastructure_deployer" "deploy" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set to false to trigger a deploy on every apply that changes the infrastructure.
  prevent_deploy      = true
  await_deploy_finish = false
  allow_data_loss     = true

  # Deploys must run after every other resource in the infrastructure.
  depends_on = [
    metalcloud_server_instance_group.inst01,
  ]
}
```

### Deploy with a longer soft shutdown window

```hcl
resource "metalcloud_infrastructure_deployer" "deploy" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  prevent_deploy      = false
  await_deploy_finish = true
  allow_data_loss     = false

  attempt_soft_shutdown         = true
  soft_shutdown_timeout_seconds = 300
  attempt_hard_shutdown         = true
  force_shutdown                = false

  depends_on = [metalcloud_server_instance_group.inst01]
}
```

## Schema

### Required

- `infrastructure_id` (String) Id of the infrastructure to deploy.

### Optional

- `allow_data_loss` (Boolean) Allow the deploy to perform operations that may cause data loss. Defaults to `false`.
- `attempt_hard_shutdown` (Boolean) Force a hard power off of servers that are still running when the soft shutdown timeout expires. Defaults to `true`.
- `attempt_soft_shutdown` (Boolean) Attempt a soft (ACPI) power off of the servers before the deploy. Defaults to `true`.
- `await_deploy_finish` (Boolean) Wait for the deploy to finish before returning. Times out after 30 minutes. Defaults to `false`.
- `force_shutdown` (Boolean) Force an immediate shutdown of all servers in the infrastructure, skipping the soft shutdown sequence. Defaults to `false`.
- `prevent_deploy` (Boolean) When `true` no deploy is triggered. Defaults to `true`.
- `soft_shutdown_timeout_seconds` (Number) Seconds to wait for the soft shutdown before attempting a hard shutdown. Defaults to `180`.

## Import

```sh
terraform import metalcloud_infrastructure_deployer.deploy <infrastructure_id>
```

Importing only records the infrastructure Id. The deploy options take their defaults until set in configuration.

## Related Resources

- [`metalcloud_infrastructure`](infrastructure.md) - Create the infrastructure
- [`metalcloud_infrastructure`](../data-sources/infrastructure.md) - Look up or create an infrastructure by label
