---
page_title: "metalcloud_storage_pool Data Source - terraform-provider-metalcloud"
description: |-
  Storage pool data source for looking up the storage pool that backs a drive
---

# metalcloud_storage_pool (Data Source)

Resolves the Id of a **storage pool** within a site. Storage pools represent the storage systems from which drives are provisioned. Use it to obtain the `storage_pool_id` required by [`metalcloud_drive`](../resources/drive.md) without hard-coding numeric Ids.

Selection works as follows:

- No pool matches the filters: an error is returned.
- Exactly one pool matches: it is used.
- Several pools match and `name` is not set: an "Ambiguous storage pool" error lists the candidate names.
- `name` is set: the pool with that name is selected, or an error is returned if none matches.

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = "uk-reading"
}

data "metalcloud_storage_pool" "iscsi" {
  site_id    = data.metalcloud_site.dc.site_id
  technology = "iscsi"
}

# Disambiguate when a site exposes several pools of the same technology.
data "metalcloud_storage_pool" "fast" {
  site_id    = data.metalcloud_site.dc.site_id
  technology = "nvme"
  name       = "nvme-tier-1"
}
```

## Schema

### Required

- `site_id` (String) Id of the site the storage pool belongs to.

### Optional

- `name` (String) Name of the storage pool. Required to disambiguate when the other filters match more than one pool.
- `technology` (String) Storage technology to filter by, for example `iscsi`, `fc` or `nvme`. When omitted, pools of any technology are considered.

### Read-Only

- `storage_pool_id` (String) Id of the storage pool.
