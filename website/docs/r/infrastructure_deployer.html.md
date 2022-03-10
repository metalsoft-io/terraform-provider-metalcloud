---
layout: "metalcloud"
page_title: "Metalcloud: metalcloud_infrastructure"
description: |-
  Controls a Metalcloud infrastructure and all it's elements such as instance arrays and others.
---

# metalcloud_infrastructure

This is the main (and only) resource that the metal cloud provider implements as all other elements are part of it. It has:

* provision control flags & other properties
* one or more [instance_array](./instance_array.html.md) blocks
* one or mare [network](./network.html.md) blocks


## Example Usage

The following example deploys 3 servers, each with a 40TB drive with CentOS 7.6. The servers are grouped into "master" (with 1 server) and "slave" (with 2 servers):

```hcl

data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test" 
    datacenter_name = "dc-1" 

}

resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to true to preserve the empty infrastructure after "terraform destroy"

  # Set this to false to trigger deploys.
  prevent_deploy = true

  # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  # instead of exiting while the deploy is ongoing
  await_deploy_finished = true

  # This option disables a safety check that MetalSoft performs to prevent accidental data loss.
  # It is required when testing delete operations.
  allow_data_loss = true

  # IMPORTANT: All resources that are part of this infrastructure need to be referenced here in order 
  # for the deploy operation to happen AFTER all of the create or edit operations
  depends_on = [
    module.tenancy_cluster
  ]

}
```

## Argument Reference

The following arguments are supported:
* `infrastructure_id` - (Required) The id of the infrastructure to which this object belongs to. Use the `infrastructure_reference` data source to retrieve this id. 
* `infrastructure_label` - (Required) **Infrastructure** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.

* `network` (Optional) Zero or more blocks of this type define **Networks**. If zero, the default 'WAN' network type is provisioned. In Cloud metal cloud deploymnets the deployment also includes the SAN network. In local deployments the SAN network is by default omitted to allow servers with a local drive to be deployed. Reffer to [network](./network.html.md) for more details.
* `instance_array` - (Required) One or more blocks of this type define **InstanceArrays** within this infrastructure. Reffer to [instance_array](./instance_array.html.md) for more details.
* `prevent_deploy` (Optional) If **True** provisioning will be omitted. From terraform's point of view everything would have finished successfully. This is usefull mainly during testing & development to prevent spending money on resources. The default value is **True**.
* `hard_shutdown_after_timeout` (Optional, default True) - The timeout can be configured with this object's `soft_shutdown_timeout_seconds property`. If false, the deploy will hang if at least a required server is still powered on. The servers may be powered off manually.
* `attempt_soft_shutdown` (Optional, default True) - An ACPI soft shutdown command will be sent to corresponding servers. If false, a hard shutdown is executed.
* `soft_shutdown_timeout_seconds` (Optional, default 180) - When the timeout expires, if `hard_shutdown_after_timeout` is true, then a hard power off will be attempted. Specifying a long timeout such as 1 day will block edits or deploying other new edits on infrastructure elements until the timeout expires or the servers are powered off. The servers may be powered off manually.
* `allow_data_loss` (Optional, default true) - If **true**, any operations that might cause data loss (stopping or deleting drives) will be conducted as if the "I understand that this operation is irreversible and that all snapshots will also be destroyed" checkbox in the interface has been checked. If **false** then the function will throw an error whenever an operation that might cause data loss (stopping or deleting drives) is encountered. The parameter servers to facilitate automatic infrastructure operations without risking the accidental loss of data.
* `skip_ansible`(Optional, default false) - If **true** some automatic provisioning steps will be skipped. This parameter should generally be ignored.
* `await_deploy_finished` (Optional, default true) - If **true**, the provider will wait until the deploy has finished before exiting. If **false**, the deploy will continue after the provider exited. No other operations are permitted on theis infrastructure during deploy.
* `keep_detaching_drives` (Optional, default true) - If **true**, the detaching Drive objects will not be deleted. If **false**, and the number of Instance objects is reduced, then the detaching Drive objects will be deleted.
* `infrastructure_custom_variables` (Optional, default []) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. 
* `server_allocation_policy` (DEPRECATED, Optional, default []) - Server allocation policies control how servers are allocated to instance arrays. This option allows the user to specify a particular server or a list of server types per instance array. Example:
  ```
    server_allocation_policy{
      instance_array_id  = metalcloud_instance_array.cluster1.instance_array_id
      allocation_policy {
        server_type_id = data.metalcloud_server_type.small.server_type_id
        server_count = 1
        server_ids = [10,22]
      }

      allocation_policy {
        server_type_id = data.metalcloud_server_type.small.server_type_id
        server_count = 1
        server_ids = [44,55]
      }
    }
  ```


## Attributes

This resource exports the following attributes:

* `infrastructure_id` - The id of the infrastructure is used for many operations. It is also the ID of the resource object.
* `instances` A property of type JSON which includes many details returned by the server-side including credentials and ips.
```
  "instance-258" = {
    "instance_array_id" = 255
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
  }
```
* `shared_drives` A list of shared drives belonging to the infrastructure, which includes information about the targets and WWN.
```
shared_drives = {
  shared-drive-306 = {
    shared_drive_targets_json = jsonencode(
      [
        {
            nPrefixSize  = 64
            nVLANID      = 200
            strIPAddress = "fdb6:959b:4444:1:0:0:0:2"
            strPort      = "spa_eth2"
            strPortalID  = "if_690"
            strTargetIQN = "iqn.1992-04.com.emc:cx.virt2133vnq80x.a2"
          },
        - {
            nPrefixSize  = 64
            nVLANID      = 200
            strIPAddress = "fdb6:959b:4444:1:0:0:0:1"
            strPort      = "spa_eth1"
            strPortalID  = "if_691"
            strTargetIQN = "iqn.1992-04.com.emc:cx.virt2133vnq80x.a1"
          },
      ]
    ),
    shared_drive_wwn = "60:06:01:60:B2:44:08:2B:78:6A:27:62:AB:DD:F9:9B"
  },
  shared-drive-307 = {
    shared_drive_targets_json = jsonencode(
      [
        {
            nPrefixSize  = 64
            nVLANID      = 200
            strIPAddress = "fdb6:959b:4444:1:0:0:0:2"
            strPort      = "spa_eth2"
            strPortalID  = "if_690"
            strTargetIQN = "iqn.1992-04.com.emc:cx.virt2133vnq80x.a2"
          },
        {
            nPrefixSize  = 64
            nVLANID      = 200
            strIPAddress = "fdb6:959b:4444:1:0:0:0:1"
            strPort      = "spa_eth1"
            strPortalID  = "if_691"
            strTargetIQN = "iqn.1992-04.com.emc:cx.virt2133vnq80x.a1"
        },
      ]
    )
  }
}
```