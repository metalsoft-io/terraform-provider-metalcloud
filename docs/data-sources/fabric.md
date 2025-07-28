---
page_title: "metalcloud_fabric Data Source - terraform-provider-metalcloud"
description: |-
  Retrieve information about a MetalCloud fabric for network infrastructure management.
---

# metalcloud_fabric (Data Source)

The `metalcloud_fabric` data source provides information about a specific fabric within a MetalCloud site. Fabrics represent the underlying network infrastructure that enables connectivity between server instances and external networks.

## What is a Fabric?

A **Fabric** in MetalCloud represents the physical network infrastructure within a site:

- Provides the underlying network connectivity for LogicalNetworks
- Spans multiple physical switches for redundancy and performance
- Enables communication between server instances and external networks
- Forms the foundation for VLAN, VXLAN, and other network virtualization technologies

## Example Usage

```hcl
# Retrieve fabric information by label
data "metalcloud_fabric" "primary" {
  label   = "fabric-primary"
  site_id = "us-west-1"
}

# Use fabric information in network configuration
resource "metalcloud_logical_network" "app_network" {
  label                        = "app-network"
  infrastructure_id           = metalcloud_infrastructure.example.infrastructure_id
  logical_network_type        = "lan"
  logical_network_ip_range    = "192.168.1.0/24"
  
  # Reference the fabric for network placement
  site_id = data.metalcloud_fabric.primary.site_id
}

# Display fabric information
output "fabric_details" {
  value = {
    id      = data.metalcloud_fabric.primary.fabric_id
    label   = data.metalcloud_fabric.primary.label
    site_id = data.metalcloud_fabric.primary.site_id
  }
}
```

## Schema

### Required

- `label` (String) The unique label identifier for the fabric within the site
- `site_id` (String) The identifier of the site where the fabric is located

### Read-Only

- `fabric_id` (String) The unique identifier for the fabric

## Use Cases

### Network Planning

Use fabric data sources to understand available network infrastructure before creating LogicalNetworks:

```hcl
# Get all available fabrics for network planning
data "metalcloud_fabric" "primary_fabric" {
  label   = "primary"
  site_id = var.target_site
}

data "metalcloud_fabric" "backup_fabric" {
  label   = "backup"
  site_id = var.target_site
}

# Create networks on specific fabrics
resource "metalcloud_logical_network" "primary_net" {
  label                = "primary-network"
  infrastructure_id    = metalcloud_infrastructure.main.infrastructure_id
  logical_network_type = "lan"
  site_id             = data.metalcloud_fabric.primary_fabric.site_id
}
```

### Multi-Site Deployments

Reference fabrics across different sites for distributed deployments:

```hcl
# West coast fabric
data "metalcloud_fabric" "west_fabric" {
  label   = "main-fabric"
  site_id = "us-west-1"
}

# East coast fabric
data "metalcloud_fabric" "east_fabric" {
  label   = "main-fabric"
  site_id = "us-east-1"
}

# Deploy infrastructure across both sites
resource "metalcloud_server_instance_group" "west_servers" {
  # ... configuration ...
  # Networks will use west_fabric
}

resource "metalcloud_server_instance_group" "east_servers" {
  # ... configuration ...
  # Networks will use east_fabric
}
```

## Important Notes

> **Fabric Availability**: Not all sites may have the same fabric labels. Verify fabric availability before referencing in configurations.

> **Network Dependencies**: LogicalNetworks depend on the underlying fabric infrastructure. Ensure the fabric supports your required network features.

> **Performance Considerations**: Different fabrics may have varying performance characteristics. Choose appropriate fabrics based on your workload requirements.

## Related Resources

- [`metalcloud_logical_network`](../resources/logical_network.md) - Create logical networks on fabrics
- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Deploy servers connected to fabric networks
- [`metalcloud_site`](./site.md) - Get information about sites containing fabrics

For more information about MetalCloud's network concepts, see the [Core Concepts & Terminology](../guides/concepts.html.md#logicalnetwork) guide.
