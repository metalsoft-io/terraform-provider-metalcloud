---
page_title: "metalcloud_logical_network_profile Data Source - terraform-provider-metalcloud"
description: |-
  Logical Network Profile data source for retrieving network configuration templates
---

# metalcloud_logical_network_profile (Data Source)

The `metalcloud_logical_network_profile` data source allows you to retrieve information about a Logical Network Profile, which defines network configuration templates for logical networks within a specific fabric.

## Overview

Logical Network Profiles serve as templates that define:

- Network type and protocols (VLAN, VXLAN, etc.)
- IP addressing schemes and subnets
- Quality of Service (QoS) parameters
- Security policies and access controls
- DHCP and DNS configurations

These profiles are associated with network fabrics and provide consistent network configurations across multiple logical networks.

## Example Usage

```hcl
# Retrieve a logical network profile by label
data "metalcloud_logical_network_profile" "web_tier" {
  fabric_id = "fabric-001"
  label     = "web-tier-profile"
}

# Use the profile information in a logical network
resource "metalcloud_logical_network" "web_network" {
  infrastructure_id         = metalcloud_infrastructure.example.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.web_tier.logical_network_profile_id
  label                     = "web-tier-network"
}
```

## Schema

### Required

- `fabric_id` (String) The unique identifier of the network fabric containing the profile
- `label` (String) The human-readable label of the Logical Network Profile

### Read-Only

- `logical_network_profile_id` (String) The unique identifier of the Logical Network Profile

## Related Resources

- [`metalcloud_logical_network`](../resources/logical_network.md) - Create logical networks using profiles
- [`metalcloud_fabric`](../data-sources/fabric.md) - Retrieve fabric information
