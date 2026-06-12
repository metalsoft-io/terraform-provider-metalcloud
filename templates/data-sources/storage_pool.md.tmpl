---
page_title: "metalcloud_storage_pool Data Source - terraform-provider-metalcloud"
description: |-
  Storage pool data source for looking up the storage pool that backs a drive
---

# metalcloud_storage_pool (Data Source)

The `metalcloud_storage_pool` data source resolves the ID of a storage pool within a site. Storage pools represent the backing storage (a configured storage system) from which drives are provisioned. Use this data source to obtain the `storage_pool_id` required by the `metalcloud_drive` resource without hard-coding numeric IDs.

## Understanding Storage Pools

A storage pool maps to a storage system registered in a site and is characterized by:

- **Site**: The site the storage system belongs to
- **Technology**: The storage technology it provides (e.g. `iscsi`, `fc`, `nvme`)
- **Name**: A human-readable name, used to disambiguate when a site exposes several pools

When you create a drive, MetalCloud needs to know which storage pool should provision it. This data source lets you select that pool by site and technology, falling back to an explicit name when the filters match more than one pool.

## Example Usage

```hcl
# Resolve the site first
data "metalcloud_site" "uk_reading" {
  label = "uk-reading"
}

# Look up a storage pool by site and technology
data "metalcloud_storage_pool" "iscsi" {
  site_id    = data.metalcloud_site.uk_reading.site_id
  technology = "iscsi"
}

# Use the resolved pool when creating a drive
resource "metalcloud_drive" "app_data" {
  infrastructure_id = metalcloud_infrastructure.main.infrastructure_id
  size_mbytes       = 100000
  storage_pool_id   = data.metalcloud_storage_pool.iscsi.storage_pool_id
  label             = "Application Data Drive"
}
```

When a site exposes more than one pool for the same technology, disambiguate with `name`:

```hcl
data "metalcloud_storage_pool" "fast" {
  site_id    = data.metalcloud_site.uk_reading.site_id
  technology = "nvme"
  name       = "nvme-tier-1"
}
```

## Schema

### Required

- `site_id` (String) Id of the site the storage pool belongs to. Typically sourced from the `metalcloud_site` data source.

### Optional

- `technology` (String) Storage technology to filter by (e.g. `iscsi`, `fc`, `nvme`). When omitted, pools of any technology in the site are considered.
- `name` (String) Name of the storage pool. Required to disambiguate when the `site_id`/`technology` filters match more than one pool. When set, the matching pool must exist or an error is returned.

### Read-Only

- `storage_pool_id` (String) The unique identifier of the resolved storage pool, suitable for use as the `storage_pool_id` attribute of a `metalcloud_drive` resource.

## Selection Behavior

- **No match**: Returns an error indicating no storage pool was found for the given site and technology.
- **Exactly one match**: That pool is used.
- **Multiple matches without `name`**: Returns an "Ambiguous storage pool" error listing the candidate names so you can set `name`.
- **`name` set**: The pool with the matching name is selected; an error is returned if no pool with that name exists in the filtered set.

## Related Resources

- **`metalcloud_drive`**: Consumes `storage_pool_id` to determine where the drive is provisioned.
- **`metalcloud_site`**: Provides the `site_id` used to scope the lookup.

For more information about MetalCloud's core concepts, see the [Core Concepts & Terminology guide](../guides/concepts.html.md).
