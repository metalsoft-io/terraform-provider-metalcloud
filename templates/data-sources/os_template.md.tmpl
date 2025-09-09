---
page_title: "metalcloud_os_template Data Source - terraform-provider-metalcloud"
description: |-
  Retrieve information about OS templates available in MetalCloud for server provisioning
---

# metalcloud_os_template (Data Source)

The `metalcloud_os_template` data source allows you to retrieve information about OS templates available in your MetalCloud environment. OS templates define the operating system and initial configuration that will be applied to server instances in a ServerInstanceGroup.

## What is an OS Template?

An **OS Template** in MetalCloud is a pre-configured operating system image that includes:

- Base operating system (Linux distributions, Windows, VMware ESXi, etc.)
- Initial system configurations and settings
- Pre-installed software packages and applications
- Custom automation scripts and configurations
- Network and security configurations

OS templates are applied to all instances within a ServerInstanceGroup, ensuring consistent deployment across your infrastructure.

## Example Usage

### Basic Usage

```hcl
data "metalcloud_os_template" "ubuntu" {
  label = "ubuntu-22-04-lts"
}

resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id = metalcloud_infrastructure.example.infrastructure_id
  label            = "web-servers"
  instance_count   = 3
  os_template_id   = data.metalcloud_os_template.ubuntu.os_template_id
  
  # ... other configuration
}
```

### Finding Available Templates

```hcl
# Use a partial label to find templates
data "metalcloud_os_template" "centos" {
  label = "centos-stream-9"
}

data "metalcloud_os_template" "windows" {
  label = "windows-server-2022"
}

data "metalcloud_os_template" "vmware" {
  label = "vmware-esxi-8"
}
```

## Schema

### Required

- `label` (String) The exact label of the OS template. This must match the template name as configured in your MetalCloud environment.

### Read-Only

- `os_template_id` (String) The unique identifier for the OS template. Use this value when configuring ServerInstanceGroups.

## Common OS Template Types

MetalCloud typically provides templates for:

- **Linux Distributions**: Ubuntu LTS, CentOS/RHEL, Debian, SUSE
- **Windows Server**: Various Windows Server versions
- **Virtualization Platforms**: VMware ESXi, Proxmox
- **Container Platforms**: Kubernetes-ready distributions
- **Custom Templates**: Organization-specific configurations

## Important Notes

> **Template Availability**: OS template availability depends on your MetalCloud site configuration. Contact your administrator if required templates are not available.

> **Instance Group Consistency**: All instances in a ServerInstanceGroup use the same OS template. You cannot mix different templates within a single group.

> **Template Updates**: OS templates are typically versioned. Ensure you're using the correct version for your deployment requirements.

## Related Resources

- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Uses OS templates for instance provisioning
- [`metalcloud_infrastructure`](../resources/infrastructure.md) - Contains ServerInstanceGroups that use OS templates

For more information about MetalCloud concepts, see the [Core Concepts & Terminology](../guides/concepts.html.md) guide.
