---
page_title: "metalcloud_logical_network Resource - terraform-provider-metalcloud"
description: |-
  Logical Network resource for creating isolated network segments within MetalCloud infrastructures
---

# metalcloud_logical_network (Resource)

A **LogicalNetwork** provides network connectivity abstraction within a MetalCloud infrastructure. It creates an isolated network segment that can be shared across multiple ServerInstanceGroups while maintaining security boundaries.

## Key Features

- **Network Isolation**: Creates isolated Layer 2 network segments
- **Multi-Switch Redundancy**: Automatically spans multiple physical switches for high availability
- **Flexible Implementation**: Supports various underlying technologies (VLAN, VXLAN, etc.) based on network profile
- **Shared Connectivity**: Can be attached to multiple ServerInstanceGroups within the same infrastructure
- **Site Distribution**: Networks can span across multiple sites for geographical distribution

## Use Cases

- **Application Tiers**: Separate frontend, backend, and database networks
- **Security Zones**: Create DMZ, internal, and management networks
- **Multi-Tenant Isolation**: Isolate different customers or departments
- **Development Environments**: Separate dev, staging, and production networks

## Example Usage

### Basic Logical Network

```hcl
resource "metalcloud_logical_network" "web_tier" {
  infrastructure_id            = metalcloud_infrastructure.example.infrastructure_id
  label                       = "web-tier-network"
  name                        = "Web Tier Network"
  logical_network_profile_id  = "vlan-default"
}
```

### Multiple Networks for Application Tiers

```hcl
# Frontend network
resource "metalcloud_logical_network" "frontend" {
  infrastructure_id            = metalcloud_infrastructure.app.infrastructure_id
  label                       = "frontend-net"
  name                        = "Frontend Network"
  logical_network_profile_id  = "vlan-frontend"
}

# Backend network
resource "metalcloud_logical_network" "backend" {
  infrastructure_id            = metalcloud_infrastructure.app.infrastructure_id
  label                       = "backend-net"
  name                        = "Backend Network"
  logical_network_profile_id  = "vlan-backend"
}

# Database network
resource "metalcloud_logical_network" "database" {
  infrastructure_id            = metalcloud_infrastructure.app.infrastructure_id
  label                       = "database-net"
  name                        = "Database Network"
  logical_network_profile_id  = "vlan-database"
}
```

### Attaching to ServerInstanceGroups

```hcl
resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id     = metalcloud_infrastructure.example.infrastructure_id
  label                = "web-servers"
  instance_count       = 3
  
  logical_network {
    logical_network_id = metalcloud_logical_network.web_tier.logical_network_id
  }
  
  logical_network {
    logical_network_id = metalcloud_logical_network.backend.logical_network_id
  }
}
```

## Schema

### Required

- `infrastructure_id` (String) Infrastructure ID where the logical network will be created. The network is scoped to this infrastructure and cannot be shared across different infrastructures.
- `label` (String) Unique identifier for the logical network within the infrastructure. Used for referencing the network in API calls and Terraform configurations. Must be unique within the infrastructure.
- `logical_network_profile_id` (String) Network profile that defines the underlying network technology and configuration. Common profiles include:
  - `vlan-default`: Standard VLAN-based network
  - `vxlan-overlay`: VXLAN overlay network for larger scale deployments
  - Site-specific profiles may be available depending on network fabric

### Optional

- `name` (String) Human-readable name for the logical network. Used in the MetalCloud UI for easier identification. If not specified, defaults to the label value.

### Read-Only

- `logical_network_id` (String) Unique system-generated identifier for the logical network. Used when referencing this network in other resources like ServerInstanceGroups.

## Important Considerations

### Network Profiles

The `logical_network_profile_id` determines the underlying network implementation:

- **VLAN-based profiles**: Traditional VLAN segmentation with site-local scope
- **VXLAN-based profiles**: Overlay networks that can span multiple sites
- **Custom profiles**: Site-specific configurations for specialized requirements

Contact your MetalCloud administrator to understand available profiles for your deployment.

### Security and Isolation

- Networks within the same infrastructure can communicate by default
- Cross-infrastructure communication requires explicit firewall rules
- Each LogicalNetwork provides Layer 2 isolation
- Additional security should be implemented at the application level

### Performance Considerations

- Network performance depends on the underlying physical infrastructure
- Multiple networks on the same ServerInstanceGroup share physical network interfaces
- Consider bandwidth requirements when designing network topology

### Lifecycle Management

- LogicalNetworks can be created and destroyed independently
- Removing a network that's attached to running instances will disrupt connectivity
- Always plan network changes in design mode before deploying

## Related Resources

- [metalcloud_server_instance_group](./server_instance_group.md) - Attach logical networks to compute instances
- [metalcloud_infrastructure](./infrastructure.md) - Container for logical networks
- [metalcloud_firewall_rule](./firewall_rule.md) - Control traffic between networks

## See Also

-
