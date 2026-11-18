---
page_title: "metalcloud_infrastructure Data Source - terraform-provider-metalcloud"
description: |-
  Infrastructure data source for looking up a MetalCloud infrastructure by label
---

# metalcloud_infrastructure (Data Source)

Resolves the Id of an existing **infrastructure** by site and label. Use it to add resources to an infrastructure that Terraform does not manage; `terraform destroy` never deletes an infrastructure looked up this way.

To create an infrastructure, or to have Terraform manage and destroy one that already exists, use the [`metalcloud_infrastructure`](../resources/infrastructure.md) resource instead, and [import](../resources/infrastructure.md#import) the existing infrastructure.

~> **`create_if_missing` is deprecated.** Terraform reads data sources during `terraform plan`, so with `create_if_missing = true` the infrastructure is created by `plan` (and by `terraform destroy`), and it is never tracked or destroyed by Terraform. Use the `metalcloud_infrastructure` resource instead. See [Migrate from `create_if_missing`](../resources/infrastructure.md#migrate-from-create_if_missing).

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = "uk-reading"
}

data "metalcloud_infrastructure" "infra" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "my-infra01"
}
```

## Schema

### Required

- `label` (String) Label of the infrastructure.
- `site_id` (String) Id of the site the infrastructure belongs to.

### Optional

- `create_if_missing` (Boolean, Deprecated) Create the infrastructure when it does not exist. The infrastructure is created during `terraform plan` and is never tracked or destroyed by Terraform; use the `metalcloud_infrastructure` resource instead. Defaults to `false`, in which case a missing infrastructure is an error.

### Read-Only

- `infrastructure_id` (String) Id of the infrastructure.
