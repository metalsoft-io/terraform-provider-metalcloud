---
layout: "metalcloud"
page_title: "Metalcloud: metalcloud_infrastructure/instance_array"
description: |-
  Controls a Metalcloud InstanceArray (collection of servers)
---

# metalcloud_infrastructure/instance_array

InstanceArrays are central to Metal Cloud. They operate on groups of indentical Instances (that have servers associated to them).
`instance_array` blocks are not a terraform resource per se. They need to be part of an [metalcloud_infrastructure](/docs/providers/metalcloud/r/infrastructure.html) resource.

In general lines an InstanceArray has the following properties:

* provision control flags & other properties
* zero or more [drive_array](/docs/providers/metalcloud/r/drive_array.html) blocks
* zero or more [interface](/docs/providers/metalcloud/r/interface.html) blocks
* zero or more [firewall_rule](/docs/providers/metalcloud/r/interface.html) blocks
* zero or more [instance_array_custom_variables](/docs/providers/metalcloud/r/instance_array_custom_variables.html) blocks
* zero or more [instance_custom_variables](/docs/providers/metalcloud/r/instance_custom_variables.html) blocks

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

        instance_array_custom_variables = {
            b = "c"
            d = "e"
            c = "f"
            r = "p"
        }

         instance_custom_variables {
              instance_index = 1
              custom_variables = {
                a = "z"

            }
        }

    
        instance_custom_variables {
            instance_index = 0
            custom_variables = {
                a = "d"

            }
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `instance_array_label` - (Required) **InstanceArray** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars. This will translate into a DNS record in the form of ```<label>.metalsoft.io``` or ```<label>.<env>.metalcloud.io``` for local deployments.
* `instance_array_instance_count` - (Required) **Instance** count. This is the number of instances in the instance array. The number of servers can be scaled up or down at any time (eg: while autoscaling). It can also be zero (or shrinked to zero) to allow stop-and-resume scenarios. 
* `instance_array_ram_gbytes` (Optional, default: 1). The minimum RAM capacity of each instance.
* `instance_array_processor_count` (Optional, default: 1). The minimum CPU count on each instance.
* `instance_array_processor_core_mhz` (Optional, default: 1000). The minimum clock speed of a CPU.
* `instance_array_processor_core_count` (Optional, default: 1). The minimum cores of a CPU.
* `instance_array_disk_count` (Optional, default: 0). The minimum number of physical disks.
* `instance_array_disk_size_mbytes` (Optional, default: 0). The minimum size of a single disk.
* `instance_array_boot_method` (Optional, default: 'pxe_iscsi'). Determines wether the server will boot from local drives or iSCSI LUNs. Possible values: 'pxe_iscsi', 'local_drives'.
* `instance_array_firewall_managed` (Optional, default: 'true'). When set to true, all firewall rules on the server are removed and the firewall rules specified in the `firewall_rule` properties are applied on the server. When set to false, the firewall rules specified in `firewall_rule` properties are ignored. The feature only works for drives that are using a supported OS template.
* `instance_array_additional_wan_ipv4_json` (Optional) This is a custom WAN configuration used in certain environments where user-provided secondary subnets and VLAN configuration is enabled. The format of this property will change in later versions of the provider. example configuration:
  ```
  instance_array_additional_wan_ipv4_json = "{\"configs\":[{\"forced_subnet_pool_id\":8,\"override_vlan_id\":100},{\"forced_subnet_pool_id\":9,\"override_vlan_id\":200},{\"forced_subnet_pool_id\":10,\"override_vlan_id\":300},{\"forced_subnet_pool_id\":11,\"override_vlan_id\":400},{\"forced_subnet_pool_id\":12,\"override_vlan_id\":500}]}"}
  ```
* `volume_template_id` (Optional, default: 0). The volume template ID (or name) to use if the servers in the InstanceArray have local disks. The template must support local install.
* `drive_array` (Optional) One or more blocks of this type define **DriveArrays** linked to this InstanceArray. Reffer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details.
* `firewall_rule` (Optional) One or more blocks of this type define firewall rules to be applied on each server of this InstanceArray. Reffer to [firewall_rule](/docs/providers/metalcloud/r/firewall_rule.html) for more details.
* `interface` (Optional) One or more blocks of this type define how the InstanceArray is connected to a Network. Reffer to [interface](/docs/providers/metalcloud/r/instance_array_interface.html) for more details.
* `instance_array_custom_variables` (Optional, default []) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. These are variables that will be applied at the `instance array` level and will override any identical ones configured at the `infrastructure` level specified via the `infrastructure_custom_variables` property. Example:
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

## Attributes

The instance array will export the following attributes:
`instance_array_id` - Which is the ID of the instance array resource. This can be accessed via `metalcloud_infrastructure.my_infra.instance_array[n].instance_array_id`

## Expanding and contracting

InstanceArrays can expand and shrink if the `instance_array_instance_count` property changes. Along with it all attached DriveArrays will shrink and contract. Reffer to [drive_array](/docs/providers/metalcloud/r/drive_array.html) for more details. 
On new instances the same FirewallRules will apply and the same server characteristics (same ServerType) will be used for new servers. If those are not available the closest match is located and used automatically.


## Hardware migrations

Instances have the ability to change hardware. If you change the characteristics of the InstanceArray (by changing the `instance_array_ram_gbytes` property for instance), the system will atempt to replace the servers associated with Instances in the Instance Array with ones that match the new requirements. This is done via a reboot.

# Outputs
The instance array will return a property of type JSON called `instances` which includes many details returned by the serverside including credentials and ips.
```
{
  "instance-258" = {
    "instance_array_id" = 255
    "instance_change_id" = 804
    "instance_created_timestamp" = "2021-08-23T14:51:43Z"
    "instance_credentials" = {
      "SharedDrives" = {
        "my-shared-drive" = {
          "storage_ip_address" = "100.98.0.6"
          "storage_port" = 3260
          "target_iqn" = "iqn.2013-01.com.redacted:storage.redacted.redacted.redacted"
        }
      }
      "idrac" = {}
      "ilo" = {
        "control_panel_url" = "https://172.18.34.34:443"
        "initial_password" = "redacted"
        "username" = "redacted"
      }
      "ip_addresses_public" = [
        {
          "instance_interface_id" = 1030
          "ip_change_id" = 1046
          "ip_hex" = "2a02cb80100000000000000000000002"
          "ip_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0002"
          "ip_id" = 764
          "ip_lease_expires" = "0000-00-00T00:00:00Z"
          "ip_operation" = {
            "instance_interface_id" = 1030
            "ip_change_id" = 1046
            "ip_deploy_status" = "finished"
            "ip_deploy_type" = "create"
            "ip_hex" = "2a02cb80100000000000000000000002"
            "ip_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0002"
            "ip_id" = 764
            "ip_label" = "ip-764"
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_subdomain" = "ip-764.subnet-362.data-network.tf-simple-test.7.us01.metalsoft.io"
            "ip_type" = "ipv6"
            "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
            "subnet_id" = 362
          }
          "ip_type" = "ipv6"
          "subnet_destination" = "wan"
          "subnet_gateway_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0001"
          "subnet_id" = 362
          "subnet_netmask_human_readable" = "ffff:ffff:ffff:ffff:0000:0000:0000:0000"
        },
        {
          "instance_interface_id" = 1030
          "ip_change_id" = 1047
          "ip_hex" = "b0dff882"
          "ip_human_readable" = "176.223.248.130"
          "ip_id" = 765
          "ip_lease_expires" = "0000-00-00T00:00:00Z"
          "ip_operation" = {
            "instance_interface_id" = 1030
            "ip_change_id" = 1047
            "ip_deploy_status" = "finished"
            "ip_deploy_type" = "create"
            "ip_hex" = "b0dff882"
            "ip_human_readable" = "176.223.248.130"
            "ip_id" = 765
            "ip_label" = "ip-765"
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_subdomain" = "ip-765.subnet-363.data-network.tf-simple-test.7.us01.metalsoft.io"
            "ip_type" = "ipv4"
            "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
            "subnet_id" = 363
          }
          "ip_type" = "ipv4"
          "subnet_destination" = "wan"
          "subnet_gateway_human_readable" = "176.223.248.129"
          "subnet_id" = 363
          "subnet_netmask_human_readable" = "255.255.255.252"
        },
        {
          "instance_interface_id" = 1030
          "ip_change_id" = 1048
          "ip_hex" = "ac010002"
          "ip_human_readable" = "172.1.0.2"
          "ip_id" = 766
          "ip_lease_expires" = "0000-00-00T00:00:00Z"
          "ip_operation" = {
            "instance_interface_id" = 1030
            "ip_change_id" = 1048
            "ip_deploy_status" = "finished"
            "ip_deploy_type" = "create"
            "ip_hex" = "ac010002"
            "ip_human_readable" = "172.1.0.2"
            "ip_id" = 766
            "ip_label" = "ip-766"
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_subdomain" = "ip-766.subnet-364.data-network.tf-simple-test.7.us01.metalsoft.io"
            "ip_type" = "ipv4"
            "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
            "subnet_id" = 364
          }
          "ip_type" = "ipv4"
          "subnet_destination" = "wan"
          "subnet_gateway_human_readable" = "172.1.0.1"
          "subnet_id" = 364
          "subnet_netmask_human_readable" = "255.255.255.252"
        },
        {
          "instance_interface_id" = 1030
          "ip_change_id" = 1049
          "ip_hex" = "ac020002"
          "ip_human_readable" = "172.2.0.2"
          "ip_id" = 767
          "ip_lease_expires" = "0000-00-00T00:00:00Z"
          "ip_operation" = {
            "instance_interface_id" = 1030
            "ip_change_id" = 1049
            "ip_deploy_status" = "finished"
            "ip_deploy_type" = "create"
            "ip_hex" = "ac020002"
            "ip_human_readable" = "172.2.0.2"
            "ip_id" = 767
            "ip_label" = "ip-767"
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_subdomain" = "ip-767.subnet-365.data-network.tf-simple-test.7.us01.metalsoft.io"
            "ip_type" = "ipv4"
            "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
            "subnet_id" = 365
          }
          "ip_type" = "ipv4"
          "subnet_destination" = "wan"
          "subnet_gateway_human_readable" = "172.2.0.1"
          "subnet_id" = 365
          "subnet_netmask_human_readable" = "255.255.255.252"
        },
        {
          "instance_interface_id" = 1030
          "ip_change_id" = 1050
          "ip_hex" = "ac030002"
          "ip_human_readable" = "172.3.0.2"
          "ip_id" = 768
          "ip_lease_expires" = "0000-00-00T00:00:00Z"
          "ip_operation" = {
            "instance_interface_id" = 1030
            "ip_change_id" = 1050
            "ip_deploy_status" = "finished"
            "ip_deploy_type" = "create"
            "ip_hex" = "ac030002"
            "ip_human_readable" = "172.3.0.2"
            "ip_id" = 768
            "ip_label" = "ip-768"
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_subdomain" = "ip-768.subnet-366.data-network.tf-simple-test.7.us01.metalsoft.io"
            "ip_type" = "ipv4"
            "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
            "subnet_id" = 366
          }
          "ip_type" = "ipv4"
          "subnet_destination" = "wan"
          "subnet_gateway_human_readable" = "172.3.0.1"
          "subnet_id" = 366
          "subnet_netmask_human_readable" = "255.255.255.252"
        },
      ]
      "ipmi" = {
        "initial_password" = "redacted"
        "ip_address" = "172.18.34.xx"
        "username" = "clientSd4bf"
        "version" = "2"
      }
      "iscsi" = {
        "gateway" = "100.64.0.1"
        "initiator_ip_address" = "100.64.0.6"
        "initiator_iqn" = "iqn.2021-08.com.redacted.redacted:instance-258"
        "netmask" = "255.255.255.248"
        "password" = "redacted"
        "username" = "redacted"
      }
      "rdp" = {}
      "remote_console" = {
        "remote_control_panel_url" = "?product=instance&id=258"
        "remote_protocol" = "ssh"
        "tunnel_path_url" = "https://us-chi-qts01-dc-api.us01.metalsoft.io/remote-console/instance-tunnel"
      }
      "ssh" = {
        "initial_password" = "redacted"
        "port" = 22
        "username" = "root"
      }
    }
    "instance_id" = 258
    "instance_interfaces" = [
      {
        "instance_id" = 258
        "instance_interface_capacity_mbps" = 10000
        "instance_interface_change_id" = 2543
        "instance_interface_id" = 1029
        "instance_interface_label" = "if0"
        "instance_interface_operation" = {
          "instance_id" = 258
          "instance_interface_capacity_mbps" = 10000
          "instance_interface_change_id" = 2543
          "instance_interface_deploy_status" = "finished"
          "instance_interface_deploy_type" = "create"
          "instance_interface_id" = 1029
          "instance_interface_label" = "if0"
          "instance_interface_subdomain" = "if0.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
          "network_id" = 159
        }
        "instance_interface_service_status" = "active"
        "instance_interface_subdomain" = "if0.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        "network_id" = 159
        "server_interface" = {
          "server_interface_mac_address" = "18:61:2f:00:02:20"
        }
      },
      {
        "instance_id" = 258
        "instance_interface_capacity_mbps" = 10000
        "instance_interface_change_id" = 2542
        "instance_interface_id" = 1030
        "instance_interface_index" = 1
        "instance_interface_ips" = [
          {
            "instance_interface_id" = 1030
            "ip_change_id" = 1046
            "ip_hex" = "2a02cb80100000000000000000000002"
            "ip_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0002"
            "ip_id" = 764
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_operation" = {
              "instance_interface_id" = 1030
              "ip_change_id" = 1046
              "ip_deploy_status" = "finished"
              "ip_deploy_type" = "create"
              "ip_hex" = "2a02cb80100000000000000000000002"
              "ip_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0002"
              "ip_id" = 764
              "ip_label" = "ip-764"
              "ip_lease_expires" = "0000-00-00T00:00:00Z"
              "ip_subdomain" = "ip-764.subnet-362.data-network.tf-simple-test.7.us01.metalsoft.io"
              "ip_type" = "ipv6"
              "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
              "subnet_id" = 362
            }
            "ip_type" = "ipv6"
            "subnet_destination" = "wan"
            "subnet_gateway_human_readable" = "2a02:cb80:1000:0000:0000:0000:0000:0001"
            "subnet_id" = 362
            "subnet_netmask_human_readable" = "ffff:ffff:ffff:ffff:0000:0000:0000:0000"
          },
          {
            "instance_interface_id" = 1030
            "ip_change_id" = 1047
            "ip_hex" = "b0dff882"
            "ip_human_readable" = "176.223.248.130"
            "ip_id" = 765
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_operation" = {
              "instance_interface_id" = 1030
              "ip_change_id" = 1047
              "ip_deploy_status" = "finished"
              "ip_deploy_type" = "create"
              "ip_hex" = "b0dff882"
              "ip_human_readable" = "176.223.248.130"
              "ip_id" = 765
              "ip_label" = "ip-765"
              "ip_lease_expires" = "0000-00-00T00:00:00Z"
              "ip_subdomain" = "ip-765.subnet-363.data-network.tf-simple-test.7.us01.metalsoft.io"
              "ip_type" = "ipv4"
              "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
              "subnet_id" = 363
            }
            "ip_type" = "ipv4"
            "subnet_destination" = "wan"
            "subnet_gateway_human_readable" = "176.223.248.129"
            "subnet_id" = 363
            "subnet_netmask_human_readable" = "255.255.255.252"
          },
          {
            "instance_interface_id" = 1030
            "ip_change_id" = 1048
            "ip_hex" = "ac010002"
            "ip_human_readable" = "172.1.0.2"
            "ip_id" = 766
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_operation" = {
              "instance_interface_id" = 1030
              "ip_change_id" = 1048
              "ip_deploy_status" = "finished"
              "ip_deploy_type" = "create"
              "ip_hex" = "ac010002"
              "ip_human_readable" = "172.1.0.2"
              "ip_id" = 766
              "ip_label" = "ip-766"
              "ip_lease_expires" = "0000-00-00T00:00:00Z"
              "ip_subdomain" = "ip-766.subnet-364.data-network.tf-simple-test.7.us01.metalsoft.io"
              "ip_type" = "ipv4"
              "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
              "subnet_id" = 364
            }
            "ip_type" = "ipv4"
            "subnet_destination" = "wan"
            "subnet_gateway_human_readable" = "172.1.0.1"
            "subnet_id" = 364
            "subnet_netmask_human_readable" = "255.255.255.252"
          },
          {
            "instance_interface_id" = 1030
            "ip_change_id" = 1049
            "ip_hex" = "ac020002"
            "ip_human_readable" = "172.2.0.2"
            "ip_id" = 767
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_operation" = {
              "instance_interface_id" = 1030
              "ip_change_id" = 1049
              "ip_deploy_status" = "finished"
              "ip_deploy_type" = "create"
              "ip_hex" = "ac020002"
              "ip_human_readable" = "172.2.0.2"
              "ip_id" = 767
              "ip_label" = "ip-767"
              "ip_lease_expires" = "0000-00-00T00:00:00Z"
              "ip_subdomain" = "ip-767.subnet-365.data-network.tf-simple-test.7.us01.metalsoft.io"
              "ip_type" = "ipv4"
              "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
              "subnet_id" = 365
            }
            "ip_type" = "ipv4"
            "subnet_destination" = "wan"
            "subnet_gateway_human_readable" = "172.2.0.1"
            "subnet_id" = 365
            "subnet_netmask_human_readable" = "255.255.255.252"
          },
          {
            "instance_interface_id" = 1030
            "ip_change_id" = 1050
            "ip_hex" = "ac030002"
            "ip_human_readable" = "172.3.0.2"
            "ip_id" = 768
            "ip_lease_expires" = "0000-00-00T00:00:00Z"
            "ip_operation" = {
              "instance_interface_id" = 1030
              "ip_change_id" = 1050
              "ip_deploy_status" = "finished"
              "ip_deploy_type" = "create"
              "ip_hex" = "ac030002"
              "ip_human_readable" = "172.3.0.2"
              "ip_id" = 768
              "ip_label" = "ip-768"
              "ip_lease_expires" = "0000-00-00T00:00:00Z"
              "ip_subdomain" = "ip-768.subnet-366.data-network.tf-simple-test.7.us01.metalsoft.io"
              "ip_type" = "ipv4"
              "ip_updated_timestamp" = "2021-08-23T14:51:43Z"
              "subnet_id" = 366
            }
            "ip_type" = "ipv4"
            "subnet_destination" = "wan"
            "subnet_gateway_human_readable" = "172.3.0.1"
            "subnet_id" = 366
            "subnet_netmask_human_readable" = "255.255.255.252"
          },
        ]
        "instance_interface_label" = "if1"
        "instance_interface_operation" = {
          "instance_id" = 258
          "instance_interface_capacity_mbps" = 10000
          "instance_interface_change_id" = 2542
          "instance_interface_deploy_status" = "finished"
          "instance_interface_deploy_type" = "create"
          "instance_interface_id" = 1030
          "instance_interface_index" = 1
          "instance_interface_label" = "if1"
          "instance_interface_subdomain" = "if1.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
          "network_id" = 158
        }
        "instance_interface_service_status" = "active"
        "instance_interface_subdomain" = "if1.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        "network_id" = 158
        "server_interface" = {
          "server_interface_mac_address" = "18:61:2f:00:02:24"
        }
      },
      {
        "instance_id" = 258
        "instance_interface_capacity_mbps" = 10000
        "instance_interface_change_id" = 2540
        "instance_interface_id" = 1031
        "instance_interface_index" = 2
        "instance_interface_label" = "if2"
        "instance_interface_operation" = {
          "instance_id" = 258
          "instance_interface_capacity_mbps" = 10000
          "instance_interface_change_id" = 2540
          "instance_interface_deploy_status" = "finished"
          "instance_interface_deploy_type" = "create"
          "instance_interface_id" = 1031
          "instance_interface_index" = 2
          "instance_interface_label" = "if2"
          "instance_interface_subdomain" = "if2.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        }
        "instance_interface_service_status" = "active"
        "instance_interface_subdomain" = "if2.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        "server_interface" = {
          "server_interface_mac_address" = "18:61:2f:00:02:28"
        }
      },
      {
        "instance_id" = 258
        "instance_interface_capacity_mbps" = 10000
        "instance_interface_change_id" = 2541
        "instance_interface_id" = 1032
        "instance_interface_index" = 3
        "instance_interface_label" = "if3"
        "instance_interface_operation" = {
          "instance_id" = 258
          "instance_interface_capacity_mbps" = 10000
          "instance_interface_change_id" = 2541
          "instance_interface_deploy_status" = "finished"
          "instance_interface_deploy_type" = "create"
          "instance_interface_id" = 1032
          "instance_interface_index" = 3
          "instance_interface_label" = "if3"
          "instance_interface_subdomain" = "if3.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        }
        "instance_interface_service_status" = "active"
        "instance_interface_subdomain" = "if3.instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
        "server_interface" = {
          "server_interface_mac_address" = "18:61:2f:00:02:2c"
        }
      },
    ]
    "instance_label" = "instance-258"
    "instance_operation" = {
      "instance_array_id" = 255
      "instance_change_id" = 804
      "instance_deploy_status" = "finished"
      "instance_deploy_type" = "create"
      "instance_id" = 258
      "instance_label" = "instance-258"
      "instance_subdomain" = "instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
      "server_id" = 2
      "server_type_id" = 1
      "template_id_origin" = 2
    }
    "instance_service_status" = "active"
    "instance_subdomain" = "instance-258.esxi-group1-1.vanilla.tf-simple-test.7.us01.metalsoft.io"
    "instance_subdomain_permanent" = "instance-258.us01.metalsoft.io"
    "instance_updated_timestamp" = "2021-08-23T15:48:08Z"
    "server_id" = 2
    "server_type_id" = 1
    "template_id_origin" = 2
  }
}
```