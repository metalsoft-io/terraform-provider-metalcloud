---
page_title: "metalcloud_vm_type Data Source - terraform-provider-metalcloud"
description: |-
  Retrieve information about VM types available in MetalCloud for configuring server instances.
---

# metalcloud_vm_type (Data Source)

The `metalcloud_vm_type` data source allows you to retrieve information about VM types available in your MetalCloud environment. VM types define the compute specifications (CPU, RAM, storage) that can be allocated to server instances within a ServerInstanceGroup.

## Overview

VM types in MetalCloud represent standardized compute configurations that define:

- **CPU specifications**: Number of cores, processor type, and performance characteristics
- **Memory allocation**: Available RAM for the instance
- **Storage capacity**: Local storage specifications
- **Network capabilities**: Available network interfaces and bandwidth

VM types are predefined by the MetalCloud platform and vary based on the physical hardware available at each site.

## Example Usage

### Basic VM Type Lookup

```hcl
data "metalcloud_vm_type" "standard_compute" {
  label = "standard.2xlarge"
}

resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id    = metalcloud_infrastructure.main.infrastructure_id
  server_type_id      = data.metalcloud_vm_type.standard_compute.vm_type_id
  instance_count      = 3
  # ... other configuration
}
```

### Multiple VM Types for Different Workloads

```hcl
# High-memory VM for database workloads
data "metalcloud_vm_type" "database" {
  label = "memory.8xlarge"
}

# Compute-optimized VM for processing workloads
data "metalcloud_vm_type" "compute" {
  label = "compute.4xlarge"
}

# General purpose VM for web services
data "metalcloud_vm_type" "general" {
  label = "general.large"
}
```

## Schema

### Required

- `label` (String) The VM type label. This is a human-readable identifier that describes the VM type's specifications (e.g., "standard.2xlarge", "memory.4xlarge", "compute.large").

### Read-Only

- `vm_type_id` (String) The unique identifier for the VM type. This ID is used when configuring ServerInstanceGroups and other resources that require VM type specification.

## Common VM Type Categories

VM types are typically categorized by their primary use case:

### General Purpose
- **standard.small**: Basic workloads, development environments
- **standard.medium**: Web applications, small databases
- **standard.large**: Production web services, application servers
- **standard.xlarge**: High-traffic applications, medium databases

### Memory Optimized
- **memory.large**: In-memory caching, real-time analytics
- **memory.xlarge**: Large databases, big data processing
- **memory.2xlarge**: High-memory applications, data warehousing

### Compute Optimized
- **compute.large**: CPU-intensive applications, batch processing
- **compute.xlarge**: High-performance computing, scientific workloads
- **compute.2xlarge**: Parallel processing, computational modeling

> **Note**: Available VM types depend on the physical hardware at each MetalCloud site. Use the MetalCloud CLI or API to list available VM types for your specific environment.

## Important Considerations

### Site-Specific Availability
- VM types may vary between different MetalCloud sites
- Always verify VM type availability at your target deployment site
- Consider multiple VM type options for multi-site deployments

### Resource Planning
- Larger VM types may have limited availability during peak usage
- Plan resource requirements in advance for production workloads
- Consider using multiple smaller instances instead of fewer large instances for better availability

### Cost Optimization
- Right-size your VM types based on actual workload requirements
- Monitor resource utilization to identify optimization opportunities
- Use design mode to test different VM type configurations before deployment

## Related Resources

- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Configure server instances using VM types
- [`metalcloud_infrastructure`](../resources/infrastructure.md) - Create infrastructures to contain your server instances
- [Core Concepts Guide](../guides/concepts.html.md) - Understanding MetalCloud's infrastructure model

## Error Handling

If the specified VM type label is not found or not available, Terraform will return an error during the plan phase. Common issues include:

- **Invalid label**: The specified label doesn't exist in the system
- **Site restrictions**: The VM type is not available at the target site
- **Capacity constraints**: The VM type is temporarily unavailable due to resource constraints

Always validate VM type availability using the MetalCloud CLI before using in production configurations:

```bash
metalcloud-cli vm-types list
```
