---
layout: "metalcloud"
page_title: "Metalcloud: kubernetes"
description: |-
  Controls a Kubernetes deployment
---


# vmware_vsphere

This structure represents a MetalCloud Kubernetes deployment.  Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

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


resource "metalcloud_kubernetes" "k8s1" {
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

}

```
## Argument Reference

* `cluster_label` (Required) *  **Cluster** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `instance_array_instance_count_master` The number of instances in the master instance array.
* `instance_server_type_master` The id of the server type to use for master nodes for each instance (see example above)
* `instance_array_instance_count_worker` The count of instances in the worker instance array.
* `instance_server_type_worker` The id of the server type to use for worker nodes for each instance (see example above)


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
cluster_credentials = "{\"testvmware\":{\"vsphere_worker\":{\"instance-1884\":{\"instance_id\":1884,\"instance_label\":\"instance-1884\",\"instance_hostname\":\"instance-1884.us01.metalsoft.io\",\"instance_cluster_url\":\"unavailable\",\"instance_health\":\"unavailable\",\"type\":\"AppVMwarevSphereInstance\",\"esxi_username\":\"root\",\"esxi_password\":\"RndHHb8sagwLjw_8P\"},\"instance-1885\":{\"instance_id\":1885,\"instance_label\":\"instance-1885\",\"instance_hostname\":\"instance-1885.us01.metalsoft.io\",\"instance_cluster_url\":\"unavailable\",\"instance_health\":\"unavailable\",\"type\":\"AppVMwarevSphereInstance\",\"esxi_username\":\"root\",\"esxi_password\":\"ThjrsLhdNg7Jfn_6K\"},\"instance-1886\":{\"instance_id\":1886,\"instance_label\":\"instance-1886\",\"instance_hostname\":\"instance-1886.us01.metalsoft.io\",\"instance_cluster_url\":\"unavailable\",\"instance_health\":\"unavailable\",\"type\":\"AppVMwarevSphereInstance\",\"esxi_username\":\"root\",\"esxi_password\":\"FCGwswEmGXr9PM_8F\"}},\"vsphere_master\":{\"instance-1881\":{\"instance_id\":1881,\"instance_label\":\"instance-1881\",\"instance_hostname\":\"instance-1881.us01.metalsoft.io\",\"instance_cluster_url\":\"unavailable\",\"instance_health\":\"unavailable\",\"type\":\"AppVMwarevSphereInstance\",\"esxi_username\":\"root\",\"esxi_password\":\"ypdFxkL9CDjrXg_8W\"},\"instance-1883\":{\"instance_id\":1883,\"instance_label\":\"instance-1883\",\"instance_hostname\":\"instance-1883.us01.metalsoft.io\",\"instance_cluster_url\":\"unavailable\",\"instance_health\":\"unavailable\",\"type\":\"AppVMwarevSphereInstance\",\"esxi_username\":\"root\",\"esxi_password\":\"wgPpKdegmKfj9S_5J\"}},\"admin_username\":\"administrator@vsphere.local\",\"cluster_software_available_versions\":[\"7.0.0\"],\"type\":\"AppVMwarevSphere\",\"vcsa_username\":\"root\",\"vcsa_initial_password\":\"LGassNtm9BLFYP\"}}"
```