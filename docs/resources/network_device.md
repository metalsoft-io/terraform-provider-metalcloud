---
page_title: "metalcloud_network_device Resource - terraform-provider-metalcloud"
description: |-
  Network device (switch) resource for registering switches in MetalCloud and optionally attaching them to a network fabric
---

# metalcloud_network_device (Resource)

Registers a single network device (switch) in MetalCloud and, optionally, attaches it to an existing network fabric. Use `for_each` over a map of switches to register a whole inventory in one configuration.

- **Optional, editable fabric assignment**: `fabric_id` attaches the switch to a fabric. Adding it attaches, changing it reassigns (detach old, attach new), and removing it detaches the switch.
- **Explicit site**: `site_id` is set explicitly rather than derived from the fabric, so a switch can exist independently of any fabric.
- **No automatic deploy**: Creating, attaching or reassigning a switch records the change via the API but does not push configuration to the hardware. Deploy the fabric separately.

## Example Usage

### Single switch attached to a fabric

```hcl
data "metalcloud_site" "dc" {
  label = "uk-reading"
}

data "metalcloud_fabric" "target" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "my-fabric"
}

resource "metalcloud_network_device" "leaf01" {
  site_id             = data.metalcloud_site.dc.site_id
  identifier_string   = "leaf-01"
  driver              = "cumulus_linux"
  position            = "leaf"
  username            = "cumulus"
  management_password = var.switch_password
  management_address  = "192.168.200.11"
  management_port     = 22
  loopback_address    = "10.253.128.1"
  asn                 = 4200000000

  fabric_id = data.metalcloud_fabric.target.fabric_id
}
```

### Bulk registration with `for_each`

```hcl
locals {
  switches = {
    "leaf-su00-r0"  = { position = "leaf", management_address = "192.168.200.11", loopback_address = "10.253.128.1", asn = 4200000000 }
    "leaf-su00-r1"  = { position = "leaf", management_address = "192.168.200.12", loopback_address = "10.253.128.2", asn = 4200000001 }
    "spine-su00-s0" = { position = "spine", management_address = "192.168.200.60" }
  }
}

resource "metalcloud_network_device" "switch" {
  for_each = local.switches

  site_id             = data.metalcloud_site.dc.site_id
  identifier_string   = each.key
  driver              = "cumulus_linux"
  username            = "cumulus"
  management_password = var.switch_password
  management_port     = 22

  position           = each.value.position
  management_address = each.value.management_address
  loopback_address   = try(each.value.loopback_address, null)
  asn                = try(each.value.asn, null)

  fabric_id = data.metalcloud_fabric.target.fabric_id
}
```

### Switch with tags, not attached to any fabric

```hcl
resource "metalcloud_network_device" "standalone" {
  site_id             = data.metalcloud_site.dc.site_id
  identifier_string   = "tor-99"
  driver              = "sonic_enterprise"
  position            = "tor"
  username            = "admin"
  management_password = var.switch_password
  management_address  = "10.20.30.50"

  tags_map = {
    environment = "production"
    managed_by  = "terraform"
  }
}
```

## Schema

### Required

- `driver` (String) Driver used to communicate with the device. Accepted values include `cumulus_linux`, `sonic_enterprise`, `nvidia_ufm`, `nvidia_dpu`, `arista_eos`, `nexus9000`, `junos`, `os_10`, `dell_s4048`, `cisco_aci51`, `cisco_ndfc`, `vmware_vds` and `dummy`.
- `management_password` (String, Sensitive) Management password. Never returned by the API, so it is not refreshed into state.
- `position` (String) Device position in the fabric, for example `leaf`, `spine`, `super_spine`, `tor`, `north`, `dpu` or `other`.
- `site_id` (String) Id of the site the switch belongs to.
- `username` (String) Management username.

### Optional

- `asn` (Number) BGP ASN.
- `fabric_id` (String) Fabric to attach the switch to. Setting it attaches, changing it reassigns, clearing it detaches. No deploy is triggered.
- `identifier_string` (String) Stable identifier string (hostname) for the switch.
- `loopback_address` (String) IPv4 loopback address.
- `management_address` (String) Management (out-of-band) IP address.
- `management_port` (Number) Management port, for example `22`.
- `serial_number` (String) Hardware serial number.
- `tags_map` (Map of String) Key/value tags applied to the switch. Values must be strings.

### Read-Only

- `network_device_id` (String) Id of the network device.

## Import

```sh
terraform import metalcloud_network_device.leaf01 <network_device_id>
```

Import does not populate `fabric_id`, because the device API does not report which fabric a device belongs to. Add `fabric_id` to the configuration afterwards; re-attaching an already attached switch is a no-op upstream.

## Notes

- **Fabric membership drift**: `Read` can only confirm that the fabric recorded in state still contains the switch. If the switch is detached outside Terraform, the next refresh clears `fabric_id`. A move to a different fabric done out-of-band is not detected unless the recorded fabric no longer contains the switch.
- **Optional fields**: unset optional attributes stay null and are only populated from the API when the device has a value, to avoid perpetual null versus empty diffs. Clearing a previously set optional attribute out-of-band may therefore not show up as drift.

## Related Resources

- [`metalcloud_fabric`](../data-sources/fabric.md) - Look up the fabric to attach switches to
- [`metalcloud_site`](../data-sources/site.md) - Look up the site a switch belongs to
