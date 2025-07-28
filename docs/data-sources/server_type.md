---
page_title: "metalcloud_server_type Data Source - terraform-provider-metalcloud"
description: |-
  Server type data source for identifying available hardware configurations in MetalCloud
---

# metalcloud_server_type (Data Source)

The `metalcloud_server_type` data source provides a mechanism to identify and retrieve information about available server hardware configurations in MetalCloud. Server types define the physical specifications (CPU, RAM, storage) that will be allocated to instances in a ServerInstanceGroup.

## Understanding Server Types

Server types in MetalCloud represent standardized hardware configurations available at specific sites. Each server type has:

- **Fixed hardware specifications**: CPU cores, RAM capacity, local storage
- **Site availability**: Not all server types are available at every site

## Example Usage

### Basic Server Type Lookup

```terraform
# Locate a specific server type by name
data "metalcloud_server_type" "large" {
  server_type_name = "M.16.16.1.v3"
}

# Use the server type in a ServerInstanceGroup
resource "metalcloud_instance_array" "web_cluster" {
  infrastructure_id    = metalcloud_infrastructure.example.infrastructure_id
  instance_array_label = "web-servers"
  
  instance_array_instance_count = 3
  instance_array_ram_gbytes     = data.metalcloud_server_type.large.server_type_ram_gbytes
  instance_array_processor_count = data.metalcloud_server_type.large.server_type_processor_count
  
  # Additional configuration...
}
```

### Multiple Server Type Comparison

```terraform
# Compare different server types for workload sizing
data "metalcloud_server_type" "small" {
  server_type_name = "M.8.8.1.v3"
}

data "metalcloud_server_type" "medium" {
  server_type_name = "M.16.16.1.v3"
}

data "metalcloud_server_type" "large" {
  server_type_name = "M.32.32.2.v3"
}

# Use locals to select appropriate type based on environment
locals {
  server_types = {
    development = data.metalcloud_server_type.small
    staging     = data.metalcloud_server_type.medium
    production  = data.metalcloud_server_type.large
  }
  
  selected_server_type = local.server_types[var.environment]
}
```

### Site-Specific Server Type Usage

```terraform
# Ensure server type is available at the target site
data "metalcloud_server_type" "compute" {
  server_type_name = "M.24.24.1.v3"
}

resource "metalcloud_instance_array" "compute_cluster" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  
  # Verify compatibility with site requirements
  instance_array_instance_count = var.instance_count
  
  # Reference server type specifications
  instance_array_ram_gbytes     = data.metalcloud_server_type.compute.server_type_ram_gbytes
  instance_array_processor_count = data.metalcloud_server_type.compute.server_type_processor_count
  
  # Additional configuration...
}
```

## Schema

### Required

- `server_type_name` (String) The name of the server type to look up. Must match exactly with available server types at the target site.

### Read-Only

- `server_type_id` (String) The unique identifier for the server type. Used by instance array resources for hardware allocation.
- `server_type_processor_count` (Number) Number of CPU cores available in this server type.
- `server_type_ram_gbytes` (Number) Amount of RAM in gigabytes available in this server type.
- `server_type_disk_count` (Number) Number of local storage devices in this server type.
- `server_type_disk_size_mbytes` (Number) Size of local storage devices in megabytes.
- `server_type_description` (String) Human-readable description of the server type specifications.

## Important Considerations

### Site Availability
- Server types are site-specific and may not be available at all locations
- Use the MetalCloud CLI or UI to verify server type availability at your target site
- Consider fallback server types for multi-site deployments

### Resource Planning
- Local storage is ephemeral and will be wiped when instances are deallocated
- For persistent storage, use attached drives rather than relying on local storage
- Consider memory and CPU requirements carefully as these cannot be changed after deployment

### Performance Implications
- Higher-spec server types may have longer provisioning times
- Network performance may vary between server type generations
- Local storage performance is tied to the physical hardware configuration

## Best Practices

1. **Standardize on server types** across environments when possible for consistency
2. **Document server type choices** and rationale for future reference
3. **Test workloads** on target server types before production deployment
4. **Monitor resource utilization** to optimize server type selection
5. **Plan for growth** by selecting server types that can accommodate future scaling needs

## Troubleshooting

### Common Issues

**Server type not found**: Verify the exact spelling and availability at your target site
```bash
metalcloud-cli server-types list --site-id <site_id>
```

**Insufficient resources**: Check site capacity for the requested server type
```bash
metalcloud-cli site show --id <site_id>
```

**Configuration mismatch**: Ensure ServerInstanceGroup configuration matches server type capabilities
