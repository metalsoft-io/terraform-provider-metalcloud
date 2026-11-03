---
page_title: "metalcloud_drive Resource - terraform-provider-metalcloud"
description: |-
  Drive resource for managing shared block-storage drives provisioned from a storage pool
---

# metalcloud_drive (Resource)

A drive is a block-storage volume provisioned from a storage pool in the infrastructure's site. It is attached to server instance groups through the `hosts` attribute and is then visible to every instance in those groups.

Creating, resizing, attaching or deleting a drive only records the change in the infrastructure design. Run a deploy with [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) to apply it.

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = var.site
}

data "metalcloud_storage_pool" "iscsi" {
  site_id    = data.metalcloud_site.dc.site_id
  technology = "iscsi"
}

data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "my-infra01"
  create_if_missing = true
}

resource "metalcloud_drive" "data01" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  storage_pool_id   = data.metalcloud_storage_pool.iscsi.storage_pool_id
  label             = "data01"
  size_mbytes       = 40960

  # Server instance groups that mount the drive.
  hosts = [
    metalcloud_server_instance_group.inst01.server_instance_group_id,
  ]
}
```

## Schema

### Required

- `infrastructure_id` (String) Id of the infrastructure the drive belongs to. Changing this value forces a new resource.
- `size_mbytes` (Number) Drive size in megabytes. Can be changed in place.
- `storage_pool_id` (String) Id of the storage pool the drive is provisioned from. Use the `metalcloud_storage_pool` data source to look it up. Changing this value forces a new resource.

### Optional

- `hosts` (List of String) Ids of the server instance groups the drive is attached to. Entries added or removed here attach or detach the drive.
- `label` (String) Drive label.
- `logical_network_id` (String) Id of the logical network the drive is exposed on.

### Read-Only

- `drive_id` (String) Id of the drive.

## Import

```sh
terraform import metalcloud_drive.data01 <drive_id>
```

## Notes

- Deleting the resource deletes the drive and its data at the next deploy.
- `size_mbytes`, `label`, `logical_network_id` and `hosts` are updated in place.
