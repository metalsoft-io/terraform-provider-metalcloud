---
page_title: "metalcloud_vm_instance_group Resource - terraform-provider-metalcloud"
description: |-
  VM Instance Group resource for managing collections of identical virtual machine instances in MetalCloud
---

# metalcloud_vm_instance_group (Resource)

A VM Instance Group is a collection of identical virtual machine instances that are managed as a single unit within MetalCloud. This resource allows you to provision and scale groups of VMs with shared configurations including OS templates, storage, and network connectivity.

## Key Features

- **Horizontal Scaling**: Dynamically adjust the number of VM instances using `instance_count`
- **Shared Configuration**: All instances in the group inherit identical configurations
- **Network Connectivity**: Connect to multiple logical networks with configurable access modes
- **Custom Variables**: Pass environment-specific variables to VM instances
- **Persistent Storage**: Attach shared drives for stateful applications

## Example Usage

### Basic VM Instance Group

```hcl
resource "metalcloud_vm_instance_group" "web_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "web-servers"
  instance_count   = 3
  vm_type_id       = "standard.large"
  os_template_id   = "ubuntu-20.04-lts"
  disk_size_gbytes = 50
}
```

### VM Instance Group with Network Connections

```hcl
resource "metalcloud_vm_instance_group" "app_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "app-servers"
  instance_count   = 2
  vm_type_id       = "compute.xlarge"
  os_template_id   = "centos-8-stream"
  disk_size_gbytes = 100

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.private.logical_network_id
      access_mode       = "private"
      tagged           = false
      mtu              = 1500
    },
    {
      logical_network_id = metalcloud_logical_network.public.logical_network_id
      access_mode       = "public"
      tagged           = true
    }
  ]
}
```

### VM Instance Group with Custom Variables

```hcl
resource "metalcloud_vm_instance_group" "database_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "database-servers"
  instance_count   = 1
  vm_type_id       = "memory.large"
  os_template_id   = "postgres-14"
  disk_size_gbytes = 200

  custom_variables = [
    {
      name  = "DB_NAME"
      value = "production_db"
    },
    {
      name  = "DB_MAX_CONNECTIONS"
      value = "200"
    },
    {
      name  = "BACKUP_SCHEDULE"
      value = "0 2 * * *"
    }
  ]

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.database.logical_network_id
      access_mode       = "private"
      tagged           = false
    }
  ]
}
```

### Scaling Example

```hcl
# Scale up the web servers based on load
resource "metalcloud_vm_instance_group" "auto_scaling_web" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "auto-scaling-web"
  instance_count   = var.web_server_count  # Can be dynamically adjusted
  vm_type_id       = "standard.medium"
  os_template_id   = "nginx-alpine"
  disk_size_gbytes = 30

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.load_balancer.logical_network_id
      access_mode       = "private"
      tagged           = false
    }
  ]
}
```

## Schema

### Required

- `disk_size_gbytes` (Number) Disk size in GB for each VM instance. This is the local storage allocated to each instance.
- `infrastructure_id` (String) Infrastructure ID where the VM instance group will be created. Must reference an existing infrastructure.
- `instance_count` (Number) Number of VM instances in the group. Must be a positive integer. Can be modified to scale the group up or down.
- `label` (String) Human-readable label for the VM instance group. Must be unique within the infrastructure.
- `os_template_id` (String) OS template ID that defines the operating system and initial configuration for all instances in the group.
- `vm_type_id` (String) VM type ID that specifies the compute resources (CPU, RAM) for each instance.

### Optional

- `custom_variables` (Attributes Set) Custom environment variables passed to all VM instances during provisioning. These can be used for application configuration, environment setup, or integration with configuration management tools. (see [below for nested schema](#nestedatt--custom_variables))
- `network_connections` (Attributes Set) Network connections that define how the VM instances connect to logical networks. Each connection specifies access mode, VLAN tagging, and other network parameters. (see [below for nested schema](#nestedatt--network_connections))

### Read-Only

- `vm_instance_group_id` (String) Unique identifier for the VM instance group, automatically assigned by MetalCloud.

<a id="nestedatt--custom_variables"></a>
### Nested Schema for `custom_variables`

Custom variables allow you to pass configuration data to VM instances during provisioning. These variables are typically used by OS templates for environment-specific configuration.

#### Required

- `name` (String) Name of the custom variable. Should follow standard environment variable naming conventions (uppercase letters, numbers, and underscores).
- `value` (String) Value of the custom variable. Can contain any string data including JSON, configuration parameters, or simple values.

#### Usage Notes

- Variables are available to the OS template during instance provisioning
- Common use cases include database configuration, application settings, and service discovery
- Variables are inherited by all instances in the group
- Consider using Terraform variables or data sources for dynamic values

<a id="nestedatt--network_connections"></a>
### Nested Schema for `network_connections`

Network connections define how VM instances connect to logical networks within the infrastructure.

#### Required

- `access_mode` (String) Access mode for the network connection. Valid values:
  - `"private"` - Internal network access only
  - `"public"` - Internet-accessible network connection
  - `"management"` - Management network for administrative access
- `logical_network_id` (String) ID of the logical network to connect to. Must reference an existing logical network within the same infrastructure.
- `tagged` (Boolean) Whether the network connection uses VLAN tagging:
  - `true` - Uses VLAN tagging for network isolation
  - `false` - Untagged connection (typically for primary networks)

#### Optional

- `mtu` (Number) Maximum Transmission Unit (MTU) size for the network connection. Default is typically 1500 bytes. Higher values (up to 9000) may improve performance for specific workloads but must be supported by the underlying network infrastructure.

#### Usage Notes

- VM instances can have multiple network connections for different purposes
- Network connections are applied to all instances in the group
- Ensure logical networks are properly configured before referencing them
- Consider network security and isolation requirements when designing connections

## Import

VM Instance Groups can be imported using their ID:

```bash
terraform import metalcloud_vm_instance_group.example 12345
```

## Important Considerations

### Scaling Operations

- **Scale Up**: Increasing `instance_count` provisions additional VM instances with identical configuration
- **Scale Down**: Decreasing `instance_count` terminates excess instances (data on local disks will be lost)
- **Zero Downtime**: Scaling operations can be performed without affecting existing instances

### Data Persistence

- **Local Storage**: Data stored on the local disk (`disk_size_gbytes`) is not persistent across instance lifecycle changes
- **Shared Storage**: For persistent data, attach drives or use network-attached storage
- **Backup Strategy**: Implement appropriate backup procedures for critical data

### Network Security

- Configure appropriate firewall rules and network ACLs
- Use private networks for internal communication
- Limit public network access to necessary services only

### Resource Planning

- Consider the total resource requirements when scaling (CPU, memory, storage, network bandwidth)
- Monitor resource utilization across the infrastructure
- Plan for peak load scenarios and scaling requirements

## Related Resources

- [`metalcloud_infrastructure`](infrastructure.md) - Container for VM instance groups
- [`metalcloud_logical_network`](logical_network.md) - Network connectivity
- [`metalcloud_os_template`](os_template.md) - Operating system configuration
- [`metalcloud_drive`](drive.md) - Persistent storage attachment
