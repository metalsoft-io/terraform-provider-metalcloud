---
page_title: "metalcloud_infrastructure Data Source - terraform-provider-metalcloud"
description: |-
  Infrastructure data source for retrieving MetalCloud infrastructure information
---

# metalcloud_infrastructure (Data Source)

The `metalcloud_infrastructure` data source allows you to retrieve information about an existing MetalCloud infrastructure or create one if it doesn't exist.

## About MetalCloud Infrastructures

An **Infrastructure** is the top-level organizational unit in MetalCloud that serves as:

- A logical grouping for related resources (servers, drives, networks)
- A security boundary for access control
- An isolated environment for workloads
- A container for ServerInstanceGroups, Drives, and LogicalNetworks

Infrastructures enable you to separate different projects, environments, or teams while maintaining clear resource boundaries.

## Example Usage

### Basic Usage - Retrieve Existing Infrastructure

```hcl
data "metalcloud_infrastructure" "example" {
  label   = "my-infrastructure"
  site_id = "us-west-1"
}

output "infrastructure_id" {
  value = data.metalcloud_infrastructure.example.infrastructure_id
}
```

### Auto-Create Infrastructure

```hcl
data "metalcloud_infrastructure" "dev_env" {
  label             = "development-environment"
  site_id           = "us-east-1"
  create_if_missing = true
}

# Use the infrastructure in other resources
resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id = data.metalcloud_infrastructure.dev_env.infrastructure_id
  # ... other configuration
}
```

### Multiple Infrastructures for Different Environments

```hcl
data "metalcloud_infrastructure" "staging" {
  label             = "staging-env"
  site_id           = "us-west-1"
  create_if_missing = true
}

data "metalcloud_infrastructure" "production" {
  label             = "production-env"
  site_id           = "us-east-1"
  create_if_missing = true
}
```

## Schema

### Required

- `label` (String) Infrastructure label. Must be unique within the site. Used to identify and reference the infrastructure.
- `site_id` (String) Site identifier where the infrastructure will be located. Determines the physical location and available resources.

### Optional

- `create_if_missing` (Boolean) If `true`, creates the infrastructure if it doesn't exist. If `false` (default), the data source will fail if the infrastructure is not found.

### Read-Only

- `infrastructure_id` (String) Unique identifier for the infrastructure. Used by other resources to reference this infrastructure.

## Important Considerations

### Infrastructure Isolation

Each infrastructure provides complete isolation for:

- **Resource Access**: Resources in one infrastructure cannot directly access resources in another
- **Network Traffic**: LogicalNetworks are scoped to a single infrastructure
- **Security Policies**: Access controls and firewall rules are infrastructure-specific

### Site Selection

When choosing a `site_id`, consider:

- **Geographic Location**: Physical proximity to users or other systems
- **Resource Availability**: Different sites may have different server types and capacity
- **Compliance Requirements**: Data residency or regulatory requirements
- **Network Connectivity**: Network latency and bandwidth considerations

### Best Practices

1. **Naming Convention**: Use descriptive labels that indicate purpose and environment (e.g., "web-app-production", "db-cluster-staging")

2. **Environment Separation**: Create separate infrastructures for different environments:

   ```hcl
   data "metalcloud_infrastructure" "dev" {
     label   = "myapp-development"
     site_id = "us-west-1"
   }
   
   data "metalcloud_infrastructure" "prod" {
     label   = "myapp-production"
     site_id = "us-east-1"
   }
   ```

3. **Team Isolation**: Use separate infrastructures for different teams or projects to maintain clear boundaries

4. **Resource Planning**: Consider the scale and resource requirements when selecting sites

## Related Resources

Infrastructures work with other MetalCloud resources:

- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Compute resources within the infrastructure
- [`metalcloud_drive`](../resources/drive.md) - Persistent storage attached to the infrastructure
- [`metalcloud_logical_network`](../resources/logical_network.md) - Network connectivity within the infrastructure

## Error Handling

If the infrastructure doesn't exist and `create_if_missing` is `false`, the data source will return an error:

```text
Error: Infrastructure with label "non-existent-infra" not found in site "us-west-1"
```

To handle this gracefully, either:

- Set `create_if_missing = true` to auto-create
- Ensure the infrastructure exists before referencing it
- Use conditional logic in your Terraform configuration

## Import

Existing infrastructures can be referenced by their label and site without needing to import them into Terraform state, as this is a data source.
