---
page_title: "metalcloud_server_instance_group Resource - terraform-provider-metalcloud"
description: |-
  A ServerInstanceGroup is a collection of identical server instances managed as a single unit, providing horizontal scaling and consistent configuration across all instances.
---

# metalcloud_server_instance_group (Resource)

A **ServerInstanceGroup** is a fundamental resource in MetalCloud that represents a collection of identical server instances managed as a single unit. All instances within the group share the same configuration including OS template, storage drives, and network connections.

## Key Features

- **Horizontal Scaling**: Dynamically adjust the number of instances using `instance_count`
- **Consistent Configuration**: All instances share identical OS templates, drives, and network settings
- **Group-Level Operations**: Management operations are performed at the group level, not individual instances
- **High Availability**: Instances are distributed across available physical servers for redundancy

## Use Cases

- **Application Clusters**: Deploy identical application servers that can be load-balanced
- **Container Orchestration**: Provision worker nodes for Kubernetes or other container platforms
- **Web Server Farms**: Scale web servers horizontally based on demand
- **Distributed Databases**: Deploy database cluster nodes with shared storage
- **Microservices**: Run multiple instances of stateless services

## Example Usage

### Basic Server Instance Group

```hcl
resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "web-cluster"
  name             = "Web Server Cluster"
  instance_count   = 3
  server_type_id   = "server_type_small"
  os_template_id   = "ubuntu_22_04"

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.public.logical_network_id
      access_mode       = "public"
      tagged           = false
    },
    {
      logical_network_id = metalcloud_logical_network.private.logical_network_id
      access_mode       = "private"
      tagged           = true
      mtu              = 9000
    }
  ]

  custom_variables = [
    {
      name  = "APP_ENV"
      value = "production"
    },
    {
      name  = "CLUSTER_SIZE"
      value = "3"
    }
  ]
}
```

### Server Instance Group with Shared Storage

```hcl
resource "metalcloud_server_instance_group" "database_cluster" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "db-cluster"
  name             = "Database Cluster"
  instance_count   = 2
  server_type_id   = "server_type_large"
  os_template_id   = "centos_8"

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.database.logical_network_id
      access_mode       = "private"
      tagged           = false
    }
  ]

  custom_variables = [
    {
      name  = "DB_CLUSTER_MODE"
      value = "primary-replica"
    }
  ]
}

# Attach shared storage to the database cluster
resource "metalcloud_drive_attachment" "db_storage" {
  drive_id                  = metalcloud_drive.shared_storage.drive_id
  server_instance_group_id  = metalcloud_server_instance_group.database_cluster.server_instance_group_id
}
```

## Schema

### Required

- `infrastructure_id` (String) The ID of the infrastructure that will contain this server instance group
- `instance_count` (Number) Number of server instances to provision in this group. Can be scaled up or down dynamically
- `label` (String) Unique identifier for the server instance group within the infrastructure. Used for internal references
- `os_template_id` (String) ID of the OS template that will be applied to all instances in the group
- `server_type_id` (String) ID of the server type that defines the hardware specifications for all instances

### Optional

- `name` (String) Human-readable name for the server instance group. If not specified, defaults to the label value
- `storage_controllers` (Attributes Set) Storage controllers configuration for the server instances (see [below for nested schema](#nestedatt--storage_controllers))
- `custom_variables` (Attributes Set) Environment variables and configuration parameters passed to all instances (see [below for nested schema](#nestedatt--custom_variables))
- `network_connections` (Attributes Set) Network interfaces and connectivity configuration for all instances (see [below for nested schema](#nestedatt--network_connections))

### Read-Only

- `server_instance_group_id` (String) Unique identifier assigned by MetalCloud after the group is created

## Nested Schema Reference

<a id="nestedatt--storage_controllers"></a>
### Nested Schema for `storage_controllers`

Required:

- `mode` (String) Storage controller mode
- `storage_controller_id` (String) Storage controller Id
- `volumes` (Attributes Set) Storage volumes configuration (see [below for nested schema](#nestedatt--storage_controllers--volumes))

<a id="nestedatt--storage_controllers--volumes"></a>
### Nested Schema for `storage_controllers.volumes`

Required:

- `controller_name` (String) Storage controller name
- `disk_count` (Number) Volume disk count
- `disk_size_gb` (Number) Volume disk size in GB
- `disk_type` (String) Volume disk type
- `raid_type` (String) Volume RAID type
- `volume_name` (String) Storage volume name

<a id="nestedatt--custom_variables"></a>
### Nested Schema for `custom_variables`

Custom variables are environment variables or configuration parameters that are passed to all instances during provisioning. These can be used by OS templates for configuration automation.

**Required:**

- `name` (String) Name of the custom variable. Must be a valid environment variable name
- `value` (String) Value of the custom variable. Will be available to all instances in the group

**Example:**
```hcl
custom_variables = [
  {
    name  = "APPLICATION_PORT"
    value = "8080"
  },
  {
    name  = "LOG_LEVEL"
    value = "INFO"
  }
]
```

<a id="nestedatt--network_connections"></a>
### Nested Schema for `network_connections`

Network connections define how instances in the group connect to logical networks. Each connection creates a network interface on all instances.

**Required:**

- `access_mode` (String) Determines the type of network access. Valid values:
  - `"public"` - Internet-accessible network interface
  - `"private"` - Internal network interface for inter-service communication
  - `"storage"` - Dedicated interface for storage traffic (iSCSI, NFS)
- `logical_network_id` (String) ID of the logical network to connect to
- `tagged` (Boolean) Whether to use VLAN tagging on this connection:
  - `true` - Use VLAN tagging (802.1Q)
  - `false` - Untagged/native VLAN

**Optional:**

- `mtu` (Number) Maximum Transmission Unit size for this network connection. Default is typically 1500. Common values:
  - `1500` - Standard Ethernet MTU
  - `9000` - Jumbo frames for high-performance applications

**Example:**
```hcl
network_connections = [
  {
    logical_network_id = metalcloud_logical_network.public.logical_network_id
    access_mode       = "public"
    tagged           = false
  },
  {
    logical_network_id = metalcloud_logical_network.storage.logical_network_id
    access_mode       = "storage"
    tagged           = true
    mtu              = 9000
  }
]
```

## Important Notes

### Instance Management

- **Group Operations**: All management operations (start, stop, scale) are performed at the group level
- **Identical Configuration**: All instances in a group are identical and cannot be individually customized
- **Hardware Distribution**: MetalCloud automatically distributes instances across available physical servers for high availability

### Scaling Considerations

- **Dynamic Scaling**: The `instance_count` can be modified to scale the group up or down
- **Zero Downtime**: Scaling operations are performed without affecting existing instances
- **Resource Limits**: Scaling is subject to available hardware resources in the infrastructure's site

### Network Behavior

- **Consistent Connectivity**: All instances receive the same network connections
- **Load Balancing**: External load balancers should be used to distribute traffic across instances
- **Internal Communication**: Instances can communicate with each other through private networks

### Storage Considerations

- **Shared Drives**: Drives attached to a ServerInstanceGroup are accessible by all instances
- **Local Storage**: Each instance has local storage that is ephemeral and wiped when the instance is released
- **Persistent Data**: Use attached drives for data that must persist across instance lifecycle changes

## Related Resources

- [`metalcloud_infrastructure`](infrastructure.md) - Container for the server instance group
- [`metalcloud_drive`](drive.md) - Persistent storage that can be attached to the group
- [`metalcloud_drive_attachment`](drive_attachment.md) - Connects drives to server instance groups
- [`metalcloud_logical_network`](logical_network.md) - Networks that instances can connect to
- [`metalcloud_server_type`](../data-sources/server_type.md) - Hardware specifications for instances
- [`metalcloud_os_template`](../data-sources/os_template.md) - Operating system configuration for instances

## Import

Server Instance Groups can be imported using their ID:

```shell
terraform import metalcloud_server_instance_group.example 12345
```
