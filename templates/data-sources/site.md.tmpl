---
page_title: "metalcloud_site Data Source - terraform-provider-metalcloud"
description: |-
  Site data source for retrieving information about MetalCloud sites and their capabilities
---

# metalcloud_site (Data Source)

The `metalcloud_site` data source allows you to retrieve information about a specific MetalCloud site. Sites represent physical locations or resource pools that contain the actual hardware infrastructure including servers, storage systems, and network equipment.

## Understanding Sites

Sites in MetalCloud serve several important purposes:

- **Physical Location**: Represent geographical locations or data centers
- **Resource Pool**: Contain physical servers, storage, and network equipment
- **Capacity Planning**: Enable resource allocation and availability planning
- **Geographical Distribution**: Allow workload placement for latency or compliance requirements

When designing infrastructures, you'll need to specify which site should host your resources based on:
- Physical proximity to users
- Available capacity
- Compliance requirements
- Network connectivity

## Example Usage

```hcl
# Retrieve information about a specific site
data "metalcloud_site" "us_west" {
  label = "us-west-datacenter"
}

# Use site information in infrastructure configuration
resource "metalcloud_infrastructure" "example" {
  infrastructure_label = "my-infrastructure"
  datacenter_name      = data.metalcloud_site.us_west.label
}

# Reference site details in outputs
output "site_information" {
  description = "Information about the selected site"
  value = {
    site_id = data.metalcloud_site.us_west.site_id
    label   = data.metalcloud_site.us_west.label
  }
}
```

## Schema

### Required

- `label` (String) The unique label identifier for the site. This is typically a human-readable name that describes the site location or purpose (e.g., "us-west-datacenter", "europe-primary").

### Read-Only

- `site_id` (String) The unique numeric identifier for the site within MetalCloud. This ID is used internally by the platform for resource allocation and management.

## Related Resources

Sites work in conjunction with other MetalCloud resources:

- **Infrastructures**: Must specify a site (via `datacenter_name`) where resources will be provisioned
- **ServerInstanceGroups**: Are allocated to physical servers within the specified site
- **Drives**: Storage resources are provisioned from the site's storage pool
- **LogicalNetworks**: Network resources span the site's network infrastructure

## Important Considerations

- Site selection affects resource availability and performance
- Different sites may have varying hardware configurations and capabilities
- Network latency between sites should be considered for multi-site deployments
- Some resources may not be available in all sites due to hardware differences

For more information about MetalCloud's core concepts and how sites fit into the overall architecture, see the [Core Concepts & Terminology guide](../guides/concepts.html.md).
