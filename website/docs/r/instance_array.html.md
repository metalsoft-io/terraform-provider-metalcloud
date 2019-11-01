---
layout: "metalcloud"
page_title: "Metalcloud: metalcloud_infrastructure/instance_array"
description: |-
  Controls a Bigstep Metalcloud InstanceArray (collection of servers)
---

# metalcloud_infrastructure/instance_array

InstanceArrays are central to Metal Cloud. They operate on groups of indentical Instances (that have servers associated to them).
`instance_array` blocks are not a terraform resource per se. They need to be part of an [metalcloud_infrastructure](/docs/providers/metalcloud/r/infrastructure.html) resource.

In general lines an InstanceArray has the following properties:

* provision control flags & other properties
* one or more [drive_array](/docs/providers/metalcloud/r/drive_array.html) blocks
* one or mare [interface](/docs/providers/metalcloud/r/interface.html) blocks
* one or mare [firewall_rule](/docs/providers/metalcloud/r/interface.html) blocks

## Example usage

The following example constructs an instance array with 2 instances, each of which have 3 network interfaces are connected 3 distinct networks (WAN, internet and a private LAN), and have a 40GB CentOS 7.6 iSCSI volume attached:

```hcl
resource "metalcloud_infrastructure" "my-infra"{
    ...

    instance_array {
        instance_array_label = "master"
        instance_array_instance_count = 2
        interface{
            interface_index = 0
            network_label = "san"
        }

        interface {
            interface_index = 1
            network_label = "internet"
        }

        interface {
            interface_index = 2
            network_label = "private"
        }
        
        drive_array {
            drive_array_label = "testia2-centos"
            drive_array_storage_type = "iscsi_hdd"
            drive_size_mbytes_default = 49000
            volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
        }

        firewall_rule {
                    firewall_rule_description = "test fw rule"
                    firewall_rule_port_range_start = 22
                    firewall_rule_port_range_end = 22
                    firewall_rule_source_ip_address_range_start="0.0.0.0"
                    firewall_rule_source_ip_address_range_end="0.0.0.0"
                    firewall_rule_protocol="tcp"
                    firewall_rule_ip_address_type="ipv4"
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `instance_array_label` - (Required) **InstanceArray** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars. This will translate into a DNS record in the form of ```<label>.bigstep.io``` or ```<label>.<env>.metalcloud.io``` for local deployments.
* `instance_array_instance_count` - (Required) **Instance** count. This is the number of instances in the instance array. The number of servers can be scaled up or down at any time (eg: while autoscaling). It can also be zero (or shrinked to zero) to allow stop-and-resume scenarios. 
* `instance_array_ram_gbytes` (Optional, default: 1). The minimum RAM capacity of each instance.
* `instance_array_processor_count` (Optional, default: 1). The minimum CPU count on each instance.
* `instance_array_processor_core_mhz` (Optional, default: 1000). The minimum clock speed of a CPU.
* `instance_array_processor_core_count` (Optional, default: 1). The minimum cores of a CPU.
* `instance_array_disk_count` (Optional, default: 0). The minimum number of physical disks.
* `instance_array_disk_size_mbytes` (Optional, default: 0). The minimum size of a single disk.
* `instance_array_boot_method` (Optional, default: 'pxe_iscsi'). Determines wether the server will boot from local drives or iSCSI LUNs. Possible values: 'pxe_iscsi', 'local_drives'.
* `instance_array_firewall_managed` (Optional, default: 'true'). When set to true, all firewall rules on the server are removed and the firewall rules specified in the `firewall_rule` properties are applied on the server. When set to false, the firewall rules specified in `firewall_rule` properties are ignored. The feature only works for drives that are using a supported OS template.
* `volume_template_id` (Optional, default: 0). The volume template ID (or name) to use if the servers in the InstanceArray have local disks. The template must support local install.
* `drive_array` (Optional) One or more blocks of this type define **DriveArrays** linked to this InstanceArray. Reffer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details.
* `firewall_rule` (Optional) One or more blocks of this type define firewall rules to be applied on each server of this InstanceArray. Reffer to [firewall_rule](/docs/providers/metalcloud/r/firewall_rule.html) for more details.
* `interface` (Optional) One or more blocks of this type define how the InstanceArray is connected to a Network. Reffer to [interface](/docs/providers/metalcloud/r/instance_array_interface.html) for more details.

## Attributes

The instance array will export the following attributes:
`instance_array_id` - Which is the ID of the instance array resource. This can be accessed via `metalcloud_infrastructure.my_infra.instance_array[n].instance_array_id`

## Expanding and contracting

InstanceArrays can expand and shrink if the `instance_array_instance_count` property changes. Along with it all attached DriveArrays will shrink and contract. Reffer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details. 
On new instances the same FirewallRules will apply and the same server characteristics (same ServerType) will be used for new servers. If those are not available the closest match is located and used automatically.


## Hardware migrations

Instances have the ability to change hardware. If you change the characteristics of the InstanceArray (by changing the `instance_array_ram_gbytes` property for instance), the system will atempt to replace the servers associated with Instances in the Instance Array with ones that match the new requirements. This is done via a reboot.

