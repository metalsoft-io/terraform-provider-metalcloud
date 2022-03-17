---
layout: "metalcloud"
page_title: "Infrastructure output: infrastructure_output"
description: |-
  Provides a mechanism to export information about an infrastructure.
---

# infrastructure_output

This data source provides a mechanism to export instance and drive information about an infrastructure using terraform outputs, based on the infrastructure id.


## Example usage

The following example locates instance credentials, drive and shared drive target information.

```hcl
data "metalcloud_infrastructure_output" "output1" {
    infrastructure_id = metalcloud_infrastructure_deployer.infrastructure_deployer.infrastructure_id
    depends_on = [
      metalcloud_infrastructure_deployer.infrastructure_deployer
    ]
}

resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  keep_detaching_drives = false
  prevent_deploy = false
  soft_shutdown_timeout_seconds = 86400
  hard_shutdown_after_timeout = false
  allow_data_loss = true

  depends_on = [
    metalcloud_drive_array.drive1,
    metalcloud_instance_array.cluster,
    metalcloud_shared_drive.datastore2,
    metalcloud_drive_array.drives2,
  ]

}

output "deployer_drives_output_data_source" { 
    value = data.metalcloud_infrastructure_output.output1
}

```

## Arguments

`infrastructure_id` (Required) The id of the infrastructure.

## Attributes

This resource exports the following attributes:

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
* `drives` A list of drives belonging to the infrastructure, which includes the WWN.
```
drives = jsonencode(
  {
    test-da = {
      drive-265 =  {
          "drive_wwn": "60:06:01:60:B2:44:08:2B:C7:04:2A:62:28:8F:91:DE"
      }
    }
  }
)
```