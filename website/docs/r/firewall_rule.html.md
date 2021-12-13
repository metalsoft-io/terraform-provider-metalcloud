---
layout: "metalcloud"
page_title: "Metalcloud: firewall rule"
description: |-
  Represents a firewall rule that is applied on an instance array.
---

# metalcloud_infrastructure/instance_array/firewall_rule

FirewallRules are a ACL-like rules that are applied on all instances of an instance array if the InstanceArray's `instance_array_firewall_managed` is set to **True**. It is part of an [instance_array](./instance_array.html.md) block.

The default rule is **Deny all** thus the FirewallRule system is a whitelist.

## Example usage

The following rules allow access to SSH access only from the HQ and HTTPs traffic from anywhere:

```hcl
data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 

    create_if_not_exists = true
}
resource "metalcloud_instance_array" "instance" {

         firewall_rule {
            firewall_rule_description = "allow ssh from HQ"
            firewall_rule_port_range_start = 22
            firewall_rule_port_range_end = 22
            firewall_rule_source_ip_address_range_start="84.84.12.0"
            firewall_rule_source_ip_address_range_end="84.84.12.255"
            firewall_rule_protocol="tcp"
            firewall_rule_ip_address_type="ipv4"
		      }

        firewall_rule {
            firewall_rule_description = "allow https traffic from anywhere"
            firewall_rule_port_range_start = 443
            firewall_rule_port_range_end = 443
            firewall_rule_source_ip_address_range_start="0.0.0.0"
            firewall_rule_source_ip_address_range_end="0.0.0.0"
            firewall_rule_protocol="tcp"
            firewall_rule_ip_address_type="ipv4"
		    }
    }
}
```

## Arguments

`firewall_rule_description` (Optional, default null) - A human readable description of the rule
`firewall_rule_port_range_start` (Optional, default null) The port range start of the firewall rule. When null, no ports are being taken into consideration when applying the firewall rule.
`firewall_rule_port_range_end` (Optional, default null) The port range end of the firewall rule. When null, no ports are being taken into consideration when applying the firewall rule.
`firewall_rule_source_ip_address_range_start` (Optional, default null) The IP address range start of the firewall rule. When null, no source IP address is taken into consideration when applying the firewall rule.
`firewall_rule_source_ip_address_range_end` (Optional, default null) The IP address range end of the firewall rule. When null, no source IP address is taken into consideration when applying the firewall rule.
`firewall_rule_destination_ip_address_range_start` (Optional, default null) The IP address range start of the firewall rule. When null, no destination IP address is taken into consideration when applying the firewall rule.
`firewall_rule_destination_ip_address_range_end` (Optional, default null) The IP address range end of the firewall rule. When null, no destination IP address is taken into consideration when applying the firewall rule.
> Setting destination rules is best avoided except in situations where the Instance is acting as a router.
`firewall_rule_protocol` (Optional, default tcp ) The protocol of the firewall rule. Possible values: *all*, *icmp*, *tcp*, *udp*.
`firewall_rule_ip_address_type` (Optional, default "ipv4") The IP address type of the firewall rule. Possible values: ipv4, ipv6
`firewall_rule_enabled` (Optional, default true) Specifies if the firewall rule will be applied or not.
