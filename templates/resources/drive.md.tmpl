---
page_title: "metalcloud_drive Resource - terraform-provider-metalcloud"
description: |-
  Drive resource for managing persistent iSCSI LUNs in MetalCloud infrastructures
---

# metalcloud_drive (Resource)

The `metalcloud_drive` resource manages iSCSI LUNs (Logical Unit Numbers) that provide persistent storage in MetalCloud. Drives are essential for stateful applications and maintain data across server instance lifecycle changes.

## Key Characteristics

- **Persistent Storage**: Data survives server instance restarts, stops, and hardware changes
- **Shared Access**: Can be attached to ServerInstanceGroups and accessed by all instances in the group
- **iSCSI Protocol**: High-performance block storage accessible over the network
- **Infrastructure Scoped**: Belongs to a specific infrastructure for security and isolation

## Schema

### Required

- `infrastructure_id` (String) The infrastructure ID where this drive will be created
- `size_mbytes` (Number) Drive size in megabytes (MB). Minimum size varies by site configuration

### Optional

- `hosts` (List of String) List of host IDs that are using this drive. When attached to a ServerInstanceGroup, all instances in the group can access the drive
- `label` (String) Human-readable label for the drive. Useful for identification and organization
- `logical_network_id` (String) Logical Network ID for network-specific drive placement. Required for drives that need specific network connectivity

### Read-Only

- `drive_id` (String) Unique identifier for the drive, assigned by MetalCloud

## Usage Examples

### Basic Drive Creation

```hcl
resource "metalcloud_drive" "app_data" {
  infrastructure_id = metalcloud_infrastructure.main.infrastructure_id
  size_mbytes      = 100000  # 100 GB
  label           = "Application Data Drive"
}
```

### Drive with Network Specification

```hcl
resource "metalcloud_drive" "database_storage" {
  infrastructure_id   = metalcloud_infrastructure.main.infrastructure_id
  size_mbytes        = 500000  # 500 GB
  label             = "Database Storage"
  logical_network_id = metalcloud_logical_network.storage_net.logical_network_id
}
```

### Drive Attached to ServerInstanceGroup

```hcl
resource "metalcloud_infrastructure" "main" {
  label           = "Production Environment"
  datacenter_name = "us-west"
}

resource "metalcloud_drive" "shared_storage" {
  infrastructure_id = metalcloud_infrastructure.main.infrastructure_id
  size_mbytes      = 1000000  # 1 TB
  label           = "Shared Application Storage"
}

resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id     = metalcloud_infrastructure.main.infrastructure_id
  label                = "Web Server Cluster"
  instance_count       = 3
  instance_server_type = "M.8"
  
  drive_id = [
    metalcloud_drive.shared_storage.drive_id
  ]
  
  # OS and network configuration...
}
```

## Important Considerations

### Data Persistence

> **Critical**: Drives provide the only persistent storage in MetalCloud. Local server storage is wiped when servers are released or reassigned. Always use drives for data that must survive infrastructure changes.

### Performance Characteristics

- **Concurrent Access**: Multiple instances can simultaneously read/write to the same drive
- **Network Dependency**: Performance depends on network connectivity between instances and storage
- **IOPS Scaling**: Performance typically scales with drive size

### Size Planning

- Consider future growth when setting `size_mbytes`
- Drive resizing may require infrastructure redeployment
- Factor in filesystem overhead (typically 5-10% of raw capacity)

### Security and Access

- Drives inherit the security context of their infrastructure
- Access is controlled at the infrastructure level
- Use separate infrastructures for security isolation between environments

## Best Practices

1. **Use Descriptive Labels**: Make drives easy to identify with meaningful labels
2. **Right-Size Storage**: Allocate appropriate capacity considering growth and performance needs
3. **Network Planning**: Use `logical_network_id` for drives requiring specific network placement
4. **Backup Strategy**: Implement backup procedures for critical drive data
5. **Monitoring**: Monitor drive usage and performance metrics

## Common Use Cases

### Database Storage
```hcl
resource "metalcloud_drive" "mysql_data" {
  infrastructure_id = metalcloud_infrastructure.production.infrastructure_id
  size_mbytes      = 200000  # 200 GB
  label           = "MySQL Data Directory"
}
```

### Shared Application Storage
```hcl
resource "metalcloud_drive" "app_shared" {
  infrastructure_id = metalcloud_infrastructure.staging.infrastructure_id
  size_mbytes      = 50000   # 50 GB
  label           = "Shared Application Files"
}
```

### Log Storage
```hcl
resource "metalcloud_drive" "log_storage" {
  infrastructure_id = metalcloud_infrastructure.monitoring.infrastructure_id
  size_mbytes      = 100000  # 100 GB
  label           = "Centralized Log Storage"
}
```

## Lifecycle Management

- **Creation**: Drives are created during infrastructure deployment
- **Attachment**: Can be attached to ServerInstanceGroups via the `drive_id` attribute
- **Modification**: Size changes typically require infrastructure redeployment
- **Deletion**: Removing a drive from configuration will delete it and all data permanently

## Integration with Other Resources

Drives are commonly used with:
- `metalcloud_server_instance_group`: For attaching persistent storage to compute instances
- `metalcloud_logical_network`: For network-specific storage placement
- `metalcloud_infrastructure`: As the containing scope for drive resources

## Troubleshooting

### Common Issues

1. **Drive Not Visible**: Ensure the ServerInstanceGroup includes the drive in its `drive_id` list
2. **Performance Issues**: Check network connectivity and consider drive placement on appropriate logical networks
3. **Size Limitations**: Verify site-specific minimum and maximum drive size limits
4. **Access Problems**: Confirm infrastructure-level permissions and network connectivity

### Validation

Always validate drive configuration in design mode (`prevent_deploy = true`) before applying changes to
