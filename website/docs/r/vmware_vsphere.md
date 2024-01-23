---
layout: "metalcloud"
page_title: "Metalcloud: vmware_vsphere"
description: |-
  Controls a Metalcloud VMWare Vsphere deployment
---


# vmware_vsphere

This structure represents a MetalCloud VMWare VSphere deployment.  Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

Don't forget to always add a depends on reference in the infrastructure deployer object.

## Example usage

```hcl
data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 
}

data "metalcloud_server_type" "large"{
     server_type_name = "M.12.8.1"
}


resource "metalcloud_network" "wan" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
    network_type = "wan"
}

data "metalcloud_network_profile" "vmware"{
    network_profile_label = "vmware-cluster"
    datacenter_name = var.datacenter
}


resource "metalcloud_vmware_vsphere" "VMWareVsphere" {
    infrastructure_id =  data.metalcloud_infrastructure.infra.infrastructure_id

    cluster_label = "testvmware"
    instance_array_instance_count_master = 2
    instance_server_type_master {
        instance_index = 0
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

    instance_server_type_master {
        instance_index = 1
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

     instance_array_network_profile_master {
        network_id = metalcloud_network.wan.id
        network_profile_id = data.metalcloud_network_profile.vmware.id
    }
    
    instance_array_instance_count_worker = 3
    instance_server_type_worker {
        instance_index = 0
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

    instance_server_type_worker {
        instance_index = 1
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

     instance_server_type_worker {
        instance_index = 2
        server_type_id = data.metalcloud_server_type.large.server_type_id
    }

    interface_worker {
      interface_index = 0
      network_id = metalcloud_network.wan.id
    }

     instance_array_network_profile_worker {
        network_id = metalcloud_network.wan.id
        network_profile_id = data.metalcloud_network_profile.vmware.id
    }

    instance_array_custom_variables_master = {
        aa = "00"
        bb = "00"
    }
 
     instance_custom_variables_master {
      instance_index = 0
      custom_variables={
        "test1":"test2"
        "test3":"test4"
      }
    }
}

```
## Argument Reference

* `cluster_label` (Required) *  **Cluster** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `instance_array_instance_count_master` The number of instances in the master instance array.
* `instance_server_type_master` The id of the server type to use for master nodes for each instance (see example above)
* `instance_array_network_profile_master` The network profile to use with this instance array. (see example above)
* `instance_array_instance_count_worker` The count of instances in the worker instance array.
* `instance_server_type_worker` The id of the server type to use for worker nodes for each instance (see example above)
* `instance_array_network_profile_worker` The network profile to use with this instance array. (see example above)
* `interface_master` The interface mapping to a network. (see example above)
* `interface_worker` The interface mapping to a network. (see example above)
* `instance_array_custom_variables_master` The instance array custom variables.
* `instance_custom_variables_master` instance level custom variables. (see example above)
* `instance_custom_variables_worker` instance level custom variables. (see example above)


## Expanding the vmware cluster

Note that it is possible to expand the cluster by editing the instance_array_instance_count_worker count but shrinking the cluster is not supported.

## Retrieving the credentials

The credentials for logging into the system can be found on the infrastructure_output object as a json object.

```hcl
data "metalcloud_infrastructure_output" "output"{
  infrastructure_id = metalcloud_infrastructure_deployer.infrastructure_deployer.infrastructure_id
  depends_on = [
      metalcloud_infrastructure_deployer.infrastructure_deployer
    ]
}

output "cluster_credentials" {
    value = data.metalcloud_infrastructure_output.output.clusters
}
```

Will output
```
cluster_credentials = "{\"testkube\":{\"admin_username\":\"admin\",\"cluster_software_available_versions\":[\"1.27.1\"],\"type\":\"AppKubernetes\"}}"
```