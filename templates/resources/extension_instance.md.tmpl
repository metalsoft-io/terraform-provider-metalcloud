---
page_title: "metalcloud_extension_instance Resource - terraform-provider-metalcloud"
description: |-
  Extension Instance resource for deploying and managing MetalCloud extensions within infrastructures
---

# metalcloud_extension_instance (Resource)

Extension Instance resource allows you to deploy and manage MetalCloud extensions within your infrastructure. Extensions provide additional functionality and services that can be integrated with your server instances and other resources.

## Overview

Extension instances are deployed within a specific infrastructure and provide services that can be consumed by server instance groups and other resources. They follow the same lifecycle management as other MetalCloud resources, supporting both design and deploy modes.

## Example Usage

### Basic Extension Instance

```hcl
resource "metalcloud_extension_instance" "monitoring" {
  label             = "monitoring-ext"
  extension_id      = "ext_monitoring_123"
  infrastructure_id = metalcloud_infrastructure.main.infrastructure_id
}
```

### Extension Instance with Dependencies

```hcl
resource "metalcloud_infrastructure" "main" {
  label = "production-env"
}

resource "metalcloud_extension_instance" "logging" {
  label             = "centralized-logging"
  extension_id      = "ext_elasticsearch_456"
  infrastructure_id = metalcloud_infrastructure.main.infrastructure_id
}

resource "metalcloud_server_instance_group" "web_servers" {
  label                    = "web-tier"
  infrastructure_id        = metalcloud_infrastructure.main.infrastructure_id
  instance_count          = 3
  server_type_id          = "server_web_1"
  os_template_id          = "ubuntu_22_04"
  
  # Extension instance can be referenced for configuration
  depends_on = [metalcloud_extension_instance.logging]
}
```

## Schema

### Required

- `extension_id` (String) The ID of the extension to deploy. This references a specific extension available in your MetalCloud environment.
- `infrastructure_id` (String) The ID of the infrastructure where the extension instance will be deployed. The extension will be scoped to this infrastructure.
- `label` (String) A human-readable label for the extension instance. Must be unique within the infrastructure.

### Read-Only

- `extension_instance_id` (String) The unique identifier assigned to the extension instance after creation.

## Extension Lifecycle

Extension instances follow MetalCloud's standard lifecycle management:

### Design Mode
When `prevent_deploy = true` is set on the infrastructure:
- Extension instances are planned but not provisioned
- Configuration can be validated without consuming resources
- Changes can be reviewed before deployment

### Deploy Mode
When `prevent_deploy = false`:
- Extension instances are provisioned and become active
- Services provided by the extension become available
- Integration with other infrastructure resources is established

## Important Considerations

### Resource Dependencies
- Extension instances are scoped to a specific infrastructure
- Some extensions may require specific network configurations or drives
- Consider extension requirements when planning infrastructure capacity

### Extension Availability
- Extension availability depends on your MetalCloud environment configuration
- Contact your MetalCloud administrator for available extensions
- Some extensions may have licensing or resource requirements

### Integration Points
- Extensions can provide services consumed by server instance groups
- Network connectivity between extensions and other resources is managed automatically within the infrastructure
- Configuration of extension services may require additional setup outside of Terraform

## State Management

Extension instances maintain their configuration state separately from the underlying physical resources, similar to other MetalCloud resources. This means:

- Extension instances can be stopped and restarted without losing configuration
- Physical resource allocation is managed by MetalCloud's orchestration layer
- Extension data persistence depends on the specific extension implementation

## Import

Extension instances can be imported using their ID:

```bash
terraform import metalcloud_extension_instance.example 12345
```

## Notes

> **Extension Compatibility**: Ensure that the extension is compatible with your infrastructure configuration and other deployed resources.

> **Resource Planning**: Some extensions may require significant resources. Plan infrastructure capacity accordingly.

> **Network Requirements**: Extensions may have specific network requirements. Ensure proper logical network configuration.
