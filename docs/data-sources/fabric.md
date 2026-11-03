---
page_title: "metalcloud_fabric Data Source - terraform-provider-metalcloud"
description: |-
  Fabric data source for looking up a MetalCloud network fabric by site and label
---

# metalcloud_fabric (Data Source)

Resolves the Id of a **network fabric** in a site. Fabrics are needed to look up logical network profiles and to attach switches with [`metalcloud_network_device`](../resources/network_device.md).

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = "uk-reading"
}

data "metalcloud_fabric" "wan" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "wan-fabric"
}
```

## Schema

### Required

- `label` (String) Label of the fabric.
- `site_id` (String) Id of the site the fabric belongs to.

### Read-Only

- `fabric_id` (String) Id of the fabric.
