---
layout: "metalcloud"
page_title: "Metalcloud: metalcloud_infrastructure"
description: |-
  Controls a Metalcloud infrastructure and all it's elements such as instance arrays and others.
---

# metalcloud_infrastructure

This is the main (and only) resource that the metal cloud provider implements as all other elements are part of it. It has:

* provision control flags & other properties
* one or more [instance_array](/docs/providers/metalcloud/r/instance_array.html) blocks
* one or mare [network](/docs/providers/metalcloud/r/network.html) blocks


## Example Usage

The following example deployes 3 servers, each with a 40TB drive with Centos 7.6. The servers are groupped into "master" (with 1 server) and "slave" (with 2 servers):

```hcl
resource "metalcloud_infrastructure" "foo" {

			infrastructure_label = "my-terraform-infra-1"
			datacenter_name = "uk-reading"

			infrastructure_custom_variables  = {
				a = "b"
				b = "a"
				c = "c"
				d = "f"
			}

			prevent_deploy = true

			network{
			  network_type = "san"
			  network_label = "san"
			}
		  
			network{
			  network_type = "wan"
			  network_label = "internet"
			}
		  
			network{
			  network_type = "lan"
			  network_label = "private"
			}
		  
		  
			instance_array {
				  instance_array_label = "master"
				  instance_array_instance_count = 1
				  interface{
					  interface_index = 0
					  network_label = "san"
				  }
		  
				  interface{
					  interface_index = 1
					  network_label = "internet"
				  }
		  
				  interface{
					  interface_index = 2
					  network_label = "private"
				  }
				  
				  drive_array{
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
		  
			instance_array {
				instance_array_label = "slave"  
				instance_array_instance_count = 2

				instance_array_custom_variables = {
					b = "c"
					d = "e"
					c = "f"
					r = "p"
				}

				instance_custom_variables {
					instance_index = 0
					custom_variables = {
						aa = "00"
						bb = "00"

					}
				}
				instance_custom_variables {
					instance_index = 1
					custom_variables = {
						# aa = "11"
						bb = "11"
						cc = "11"
						# d = "11"
					}
				}
		
				drive_array{
					drive_array_label="asd2-centos"
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


			firmware_upgrade_policy {
				server_firmware_upgrade_policy_label = "test1"
				server_firmware_upgrade_policy_action = "accept"
				instance_array_label = "web-servers"
				server_firmware_upgrade_policy_rules {
					operation = "string_equal"
					property = "datacenter_name"
					value = "slavedatacenter-138"
				}
			}

			firmware_upgrade_policy {
				server_firmware_upgrade_policy_label = "test2"
				server_firmware_upgrade_policy_action = "accept"
				instance_array_label = "web-servers"
				server_firmware_upgrade_policy_rules {
					operation = "string_equal"
					property = "datacenter_name"
					value = "slavedatacenter-138"
				}

				server_firmware_upgrade_policy_rules {
					operation = "string_equal"
					property = "server_vendor"
					value = "dell"
				}
			}

		}
```

## Argument Reference

The following arguments are supported:

* `infrastructure_label` - (Required) **Infrastructure** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `datacenter_name` - (Required) The name of the **Datacenter** where the provisioning will take place. Check the MetalCloud provider for available options.
* `network` (Optional) Zero or more blocks of this type define **Networks**. If zero, the default 'WAN' network type is provisioned. In Cloud metal cloud deploymnets the deployment also includes the SAN network. In local deployments the SAN network is by default omitted to allow servers with a local drive to be deployed. Reffer to [network](/docs/providers/metalcloud/r/network.html) for more details.
* `instance_array` - (Required) One or more blocks of this type define **InstanceArrays** within this infrastructure. Reffer to [instance_array](/docs/providers/metalcloud/r/instance_array.html) for more details.
* `prevent_deploy` (Optional) If **True** provisioning will be omitted. From terraform's point of view everything would have finished successfully. This is usefull mainly during testing & development to prevent spending money on resources. The default value is **True**.
* `hard_shutdown_after_timeout` (Optional, default True) - The timeout can be configured with this object's `soft_shutdown_timeout_seconds property`. If false, the deploy will hang if at least a required server is still powered on. The servers may be powered off manually.
* `attempt_soft_shutdown` (Optional, default True) - An ACPI soft shutdown command will be sent to corresponding servers. If false, a hard shutdown is executed.
* `soft_shutdown_timeout_seconds` (Optional, default 180) - When the timeout expires, if `hard_shutdown_after_timeout` is true, then a hard power off will be attempted. Specifying a long timeout such as 1 day will block edits or deploying other new edits on infrastructure elements until the timeout expires or the servers are powered off. The servers may be powered off manually.
* `allow_data_loss` (Optional, default true) - If **true**, any operations that might cause data loss (stopping or deleting drives) will be conducted as if the "I understand that this operation is irreversible and that all snapshots will also be destroyed" checkbox in the interface has been checked. If **false** then the function will throw an error whenever an operation that might cause data loss (stopping or deleting drives) is encountered. The parameter servers to facilitate automatic infrastructure operations without risking the accidental loss of data.
* `skip_ansible`(Optional, default false) - If **true** some automatic provisioning steps will be skipped. This parameter should generally be ignored.
* `await_deploy_finished` (Optional, default true) - If **true**, the provider will wait until the deploy has finished before exiting. If **false**, the deploy will continue after the provider exited. No other operations are permitted on theis infrastructure during deploy.
* `await_delete_finished` (Optional, default false) - If **true**, the provider will wait for a deploy (involving delete) to finish before exiting. If **false**, the delete operation (really a deploy) will continue after the provider existed. This operation is generally quick.
* `keep_detaching_drives` (Optional, default true) - If **true**, the detaching Drive objects will not be deleted. If **false**, and the number of Instance objects is reduced, then the detaching Drive objects will be deleted.
* `infrastructure_custom_variables` (Optional, default []) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. 


## Attributes

This resource exports the following attributes:

* `infrastructure_id` - The id of the infrastructure is used for many operations. It is also the ID of the resource object.
