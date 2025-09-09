---
page_title: "metalcloud_extension Data Source - terraform-provider-metalcloud"
description: |-
  Use the extension data source to retrieve information about MetalCloud extensions for infrastructure automation and custom functionality.
---

# metalcloud_extension (Data Source)

The `metalcloud_extension` data source allows you to retrieve information about MetalCloud extensions. Extensions provide additional functionality and automation capabilities for your infrastructure components.

## Example Usage

### Basic Usage

```hcl
data "metalcloud_extension" "monitoring_agent" {
  label = "monitoring-agent-v1.2"
}

# Use the extension in a server instance group
resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  # ... other configuration ...
  
  extension_id = data.metalcloud_extension.monitoring_agent.extension_id
}
```

### Finding Extensions by Pattern

```hcl
# Find the latest version of a specific extension family
data "metalcloud_extension" "docker_latest" {
  label = "docker-ce-latest"
}

output "docker_extension_info" {
  value = {
    id    = data.metalcloud_extension.docker_latest.extension_id
    label = data.metalcloud_extension.docker_latest.label
  }
}
```

## Schema

### Required

- `label` (String) The unique label identifier for the extension. This should match the exact extension name as configured in MetalCloud.

### Read-Only

- `extension_id` (String) The unique identifier for the extension. This ID is used when referencing the extension in other resources.

## Important Notes

> **Extension Availability**: Extensions must be properly configured and available in your MetalCloud environment before they can be referenced.

> **Version Management**: Extension labels often include version information. Ensure you're referencing the correct version for your use case.

> **Dependency Management**: Some extensions may have dependencies on other extensions or specific OS templates. Verify compatibility before deployment.

## Common Use Cases

### 1. Infrastructure Monitoring

Extensions are commonly used to deploy monitoring agents across server instances:

```hcl
data "metalcloud_extension" "datadog_agent" {
  label = "datadog-agent-v7"
}
```

### 2. Container Runtime Installation

Deploy container runtimes like Docker or containerd:

```hcl
data "metalcloud_extension" "docker_ce" {
  label = "docker-ce-20.10"
}
```

### 3. Security Agents

Install security monitoring and compliance tools:

```hcl
data "metalcloud_extension" "security_agent" {
  label = "security-monitoring-v2.1"
}
```

### 4. Application-Specific Tools

Deploy custom applications or tools specific to your workload:

```hcl
data "metalcloud_extension" "custom_app" {
  label = "my-custom-application-v1.0"
}
```

## Best Practices

1. **Use Specific Versions**: Always reference specific extension versions rather than "latest" tags for production deployments
2. **Test Extensions**: Validate extension functionality in development environments before production use
3. **Document Dependencies**: Maintain clear documentation of extension dependencies and compatibility requirements
4. **Version Control**: Track extension versions alongside your Terraform configuration for reproducible deployments

## Related Resources

- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Server instance groups that can use extensions
- [`metalcloud_os_template`](../data-sources/os_template.md) - OS templates that may be required by certain extensions
