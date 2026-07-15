---
page_title: "metalcloud_network_device Resource - terraform-provider-metalcloud"
description: |-
  Network device (switch) resource for registering switches in MetalCloud and optionally attaching them to a network fabric
---

# metalcloud_network_device (Resource)

The `metalcloud_network_device` resource registers a single network device (switch) in MetalCloud and, optionally, attaches it to an existing network fabric. It is the Terraform-native equivalent of a bulk switch import: use `for_each` over a map of switches to register a whole inventory in one configuration.

## Key Characteristics

- **Per-switch lifecycle**: Each switch is an independent resource with full create/read/update/delete and import support.
- **Optional, editable fabric assignment**: The `fabric_id` attribute attaches the switch to a fabric. Adding it attaches, changing it reassigns (detach old, attach new), and removing it detaches the switch.
- **Explicit site**: Unlike the standalone import script, the switch's `site_id` is set explicitly rather than derived from the fabric, so a switch can exist independently of any fabric.
- **No automatic deploy**: Creating, attaching, or reassigning a switch records the change via the API but does not push configuration to the hardware. Deploy the fabric separately to apply changes.

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

### Bulk import with `for_each`

```hcl
locals {
  switch_defaults = {
    driver          = "cumulus_linux"
    username        = "cumulus"
    management_port = 22
  }

  switches = {
    "leaf-su00-r0" = { position = "leaf", management_address = "192.168.200.11", loopback_address = "10.253.128.1", asn = 4200000000 }
    "leaf-su00-r1" = { position = "leaf", management_address = "192.168.200.12", loopback_address = "10.253.128.2", asn = 4200000001 }
    "spine-su00-s0" = { position = "spine", management_address = "192.168.200.60" }
  }
}

resource "metalcloud_network_device" "switch" {
  for_each = local.switches

  site_id             = data.metalcloud_site.dc.site_id
  identifier_string   = each.key
  driver              = local.switch_defaults.driver
  username            = local.switch_defaults.username
  management_password = var.switch_password
  management_port     = local.switch_defaults.management_port

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

  # No fabric_id: the switch is registered but not attached to a fabric.
}
```

## Schema

### Required

- `site_id` (String) Site ID the switch belongs to. Wire this from the `metalcloud_site` or `metalcloud_fabric` data source.
- `driver` (String) Driver used to communicate with the device. Common values: `cumulus_linux`, `sonic_enterprise`, `nvidia_ufm`, `nvidia_dpu`, `arista_eos`, `nexus9000`, `junos`, `os_10`, `dell_s4048`, `dummy`.
- `position` (String) Device position in the fabric. One of `tor`, `north`, `spine`, `super_spine`, `leaf`, `dpu`, `other`.
- `username` (String) Management username.
- `management_password` (String, Sensitive) Management password. Never returned by the API, so it is not refreshed into state.

### Optional

- `management_address` (String) Management (out-of-band) IP address.
- `management_port` (Number) Management port (e.g. `22`).
- `identifier_string` (String) Stable identifier string (hostname) for the switch.
- `loopback_address` (String) IPv4 loopback address.
- `asn` (Number) BGP ASN.
- `serial_number` (String) Hardware serial number.
- `tags_map` (Map of String) Key/value tags applied to the switch. Values must be strings.
- `fabric_id` (String) Fabric to attach the switch to. Editable: setting it attaches, changing it reassigns (detach old + attach new), clearing it detaches. No deploy is triggered.

### Read-Only

- `network_device_id` (String) Unique system-generated identifier for the network device.

## Important Considerations

### Idempotency and importing existing switches

This resource always creates a new device on `terraform apply`. To bring a switch that already exists in MetalCloud under Terraform management, use `terraform import` rather than re-declaring it (a re-declare would attempt to create a duplicate):

```sh
terraform import metalcloud_network_device.leaf01 <network_device_id>
```

A bare import does not populate `fabric_id` (the device API does not report which fabric it belongs to). Add the `fabric_id` afterwards; re-attaching an already-attached switch is a no-op upstream.

### Fabric membership drift

Because the device API does not report its fabric, `Read` can only confirm that the fabric currently recorded in state still contains the switch. If the switch is detached from that fabric outside Terraform, the next refresh clears `fabric_id` to surface the drift. Detaching and re-attaching to a *different* fabric out-of-band is not detected unless the state's fabric no longer contains it.

### No deploy

Attaching, reassigning, or detaching a switch only records the intent via the API. To apply the change to the physical switches, deploy the fabric separately (for example through your fabric deploy workflow). This mirrors the behavior of the standalone switch-import script.

### Optional field handling

Unset optional attributes are left null and are only populated from the API when the device actually has a value, to avoid perpetual `null` ↔ empty/zero diffs. As a consequence, clearing a previously-set optional attribute may not be detected as drift on refresh.

## Related Resources

- [metalcloud_fabric](../data-sources/fabric.md) - Look up an existing fabric (and its site) to attach switches to.
- [metalcloud_site](../data-sources/site.md) - Look up the site a switch belongs to.
