---
layout: "metalcloud"
page_title: "Metalcloud: instance_array"
description: |-
  Controls a Metalcloud InstanceArray
---

# instance_array

InstanceArrays are central to Metal Cloud. They operate on groups of identical Instances (that have servers associated to them). They need to be part of an Infrastructure.  Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

In general lines an InstanceArray has the following properties:

* provision control flags & other properties
* zero or more [interface](./instance_array_interface.html.md) blocks
* zero or more [firewall_rule](./firewall_rule.html.md) blocks
* zero or more [instance_array_custom_variables](./instance_array_custom_variables.html.md) blocks
* zero or more [instance_custom_variables](./instance_custom_variables.html.md) blocks
* zero or more [network_profile](./network_profile.html.md) blocks

## Example usage

The following example constructs an instance array:

```hcl
data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "centos-8"
}

data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 
}

resource "metalcloud_instance_array" "cluster" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    instance_array_label = "test-cluster"

    instance_array_ram_gbytes = 1
    instance_array_processor_count = 1
    instance_array_processor_core_count = 1
    instance_array_boot_method = "local_drives"
    instance_array_disk_count = 1

    volume_template_id = tonumber(data.metalcloud_volume_template.esxi7.id)

    instance_array_firewall_managed = false

    interface{
      interface_index = 0
      network_label = "storage-network"
    }

    interface{
      interface_index = 1
      network_label = "data-network"
    }

    instance_custom_variables {
      instance_index = 1
      custom_variables={
        "test1":"test2"
        "test3":"test4"
      }
    }

     network_profile {
        network_id= metalcloud_network.wan.network_id
        network_profile_id = metalcloud_network_profile.myprofile.network_profile_id
  }

}
```

## Argument Reference

The following arguments are supported:

* `instance_array_label` - (Required) **InstanceArray** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars. This will translate into a DNS record in the form of ```<label>.metalsoft.io``` or ```<label>.<env>.metalcloud.io``` for local deployments.
* `instance_array_instance_count` - (Optional) DEPRECATED**Instance** count.  This is the number of instances in the instance array. The number of servers can be scaled up or down at any time (eg: while autoscaling). It can also be zero (or reduced to zero) to allow stop-and-resume scenarios. 
* `instance_array_ram_gbytes` (Optional, default: 1). The minimum RAM capacity of each instance.
* `instance_array_processor_count` (Optional, default: 1). The minimum CPU count on each instance.
* `instance_array_processor_core_mhz` (Optional, default: 1000). The minimum clock speed of a CPU.
* `instance_array_processor_core_count` (Optional, default: 1). The minimum cores of a CPU.
* `instance_array_disk_count` (Optional, default: 0). The minimum number of physical disks.
* `instance_array_disk_size_mbytes` (Optional, default: 0). The minimum size of a single disk.
* `instance_array_boot_method` (Optional, default: 'pxe_iscsi'). Determines wether the server will boot from local drives or iSCSI LUNs. Possible values: 'pxe_iscsi', 'local_drives'.
* `instance_array_firewall_managed` (Optional, default: `true`). When set to true, all firewall rules on the server are removed and the firewall rules specified in the `firewall_rule` properties are applied on the server. When set to false, the firewall rules specified in `firewall_rule` properties are ignored. The feature only works for drives that are using a supported OS template.
* `instance_array_additional_wan_ipv4_json` (Optional) DEPRECATED This is a custom WAN configuration used in certain environments where user-provided secondary subnets and VLAN configuration is enabled. The format of this property will change in later versions of the provider. example configuration:
  ```
  instance_array_additional_wan_ipv4_json = "{\"configs\":[{\"forced_subnet_pool_id\":8,\"override_vlan_id\":100},{\"forced_subnet_pool_id\":9,\"override_vlan_id\":200},{\"forced_subnet_pool_id\":10,\"override_vlan_id\":300},{\"forced_subnet_pool_id\":11,\"override_vlan_id\":400},{\"forced_subnet_pool_id\":12,\"override_vlan_id\":500}]}"}
  ```
* `volume_template_id` (Optional, default: `0`). The volume template ID (or name) to use if the servers in the InstanceArray have local disks. The template must support local install.
* `drive_array` (Optional, default: `none`) One or more blocks of this type define **DriveArrays** linked to this InstanceArray. Refer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details.
* `firewall_rule` (Optional, default BLOCK ALL) One or more blocks of this type define firewall rules to be applied on each server of this InstanceArray. Reffer to [firewall_rule](/docs/providers/metalcloud/r/firewall_rule.html) for more details.
* `interface` (Optional) One or more blocks of this type define how the InstanceArray is connected to a Network. Refer to [interface](/docs/providers/metalcloud/r/instance_array_interface.html) for more details.
* `instance_array_custom_variables` (Optional, default: `[]`) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. These are variables that will be applied at the `instance array` level and will override any identical ones configured at the `infrastructure` level specified via the `infrastructure_custom_variables` property. Example:
  ```
  instance_array_custom_variables = {
              b = "c"
              d = "e"
              c = "f"
              r = "p"
  }
  ```
* `instance_custom_variables` (Optional, default []) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. These are variables that will be applied at the **instance** level and will override any identical ones configured at the **infrastructure** and **instance_array** level via the `infrastructure_custom_variables` and `instance_array_custom_variables` properties. Use the `instance_index` property to specify which from the instance array's instances this set of variables applies to. For example the variables for the second instance of an array would be:
  ```
  instance_custom_variables {
      instance_index = 1
      custom_variables = {
          aa = "00"
          bb = "00"
      }
  }
  ```
* `network_profile` (Optional, default []) - Configures the  network connections that the instance array has by applying profiles to them. See [network_profile](/docs/providers/metalcloud/r/network_profile.html) for more details. Example:
  ```
   network_profile {
      network_id= metalcloud_network.wan.network_id
      netowork_profile_id = metalcloud_network_profile.myprofile.network_profile_id
   }
  ```
* `instance_server_type` (Optional, default []) - Configures the  server_types of instances part of this instance array. This is an alternative method to using `instance_array_ram_gbytes` and the other "minimums" and if set will take precedence. Example:
  ```
    data "metalcloud_server_type" "large"{
      server_type_name = "M.16.16.1.v3"
    }

    instance_server_type{
      instance_index=0
      server_type_id=data.metalcloud_server_type.large.server_type_id
    }
  ```


## Attributes

The instance array will export the following attributes:
`instance_array_id` - Which is the ID of the instance array resource.

## Creating multiple identical instance arrays
The `instance_array_instance_count` property is deprecated. Please use the `count` terraform keyword to create multiple identical instances.
```
resource "metalcloud_instance_array" "cluster" {
  ...
  count =  length(var.compute_nodes)
  ...
  # use  the count.index property to distinguish between instance_arrays:
  instance_array_label = "node - ${count.index}"
}
```

## Expanding and contracting

InstanceArrays can expand and shrink if the `instance_array_instance_count` property changes. Along with it all attached DriveArrays will shrink and contract. Refer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details. 
On new instances the same FirewallRules will apply and the same server characteristics (same ServerType) will be used for new servers. If those are not available the closest match is located and used automatically.


## Hardware migrations

Instances booted from SAN have the ability to change hardware. If you change the characteristics of the InstanceArray (by changing the `instance_array_ram_gbytes` property for instance), the system will attempt to replace the servers associated with Instances in the Instance Array with ones that match the new requirements. This is done via a reboot.
