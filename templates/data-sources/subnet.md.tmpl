---
page_title: "metalcloud_subnet Data Source - terraform-provider-metalcloud"
description: |-
  Use this data source to retrieve information about a MetalCloud subnet for use in network configurations.
---

# metalcloud_subnet (Data Source)

Use this data source to retrieve information about a MetalCloud subnet. Subnets define IP address ranges and network configurations that can be used by ServerInstanceGroups and other network resources.

## Example Usage

```hcl
# Retrieve subnet information by label
data "metalcloud_subnet" "example" {
  label = "production-subnet"
}

# Use subnet in a ServerInstanceGroup
resource "metalcloud_server_instance_group" "web_servers" {
  infrastructure_id = var.infrastructure_id
  # ... other configuration ...
  
  network_interfaces {
    logical_network_id = data.metalcloud_subnet.example.logical_network_id
    subnet_id         = data.metalcloud_subnet.example.subnet_id
  }
}

# Reference subnet properties
output "subnet_cidr" {
  value = data.metalcloud_subnet.example.subnet_cidr
}
```

## Argument Reference

### Required

- `label` (String) The label of the subnet to retrieve. Labels are unique identifiers for subnets within MetalCloud.

### Optional

- `infrastructure_id` (Number) The ID of the infrastructure containing the subnet. If not specified, searches across all accessible infrastructures.

## Notes

- Subnets are associated with LogicalNetworks and define the IP addressing scheme for network interfaces
- The subnet must exist and be accessible to the user's account
- Subnet configurations are applied during infrastructure deployment
- Changes to subnet properties may require infrastructure redeployment

## Related Resources

- [`metalcloud_logical_network`](logical_network.md) - Manage logical networks
- [`metalcloud_server_instance_group`](../resources/server_instance_group.md) - Server instances that use subnets
- [`metalcloud_infrastructure`](../resources/infrastructure.md) - Infrastructure containing subnets
