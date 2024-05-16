---
layout: "metalcloud"
page_title: "Metalcloud: kubernetes"
description: |-
  Controls a Kubernetes deployment
---


# EKS-A

This structure represents a MetalCloud AWS EKS-A deployment.  Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

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


data "metalcloud_subnet_pool" "wan" {
	subnet_pool_label = "wan"
}

resource "metalcloud_subnet" "kube_boot_network" {
		subnet_type = "ipv4"
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		subnet_label="kube_boot_network"
		cluster_id = metalcloud_eksa.cluster01.cluster_id
		network_id = metalcloud_network.wan.network_id
		subnet_pool_id = data.metalcloud_subnet_pool.wan
		subnet_automatic_allocation = false
		subnet_is_ip_range = true
		subnet_ip_range_ip_count = 5
		subnet_override_vlan_id=1003
		
		
}

resource "metalcloud_subnet" "kube_vip_network" {
		subnet_type = "ipv4"
		subnet_label="kube_vip_network"
		infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
		cluster_id = metalcloud_eksa.cluster01.cluster_id
		network_id = metalcloud_network.wan.network_id
		subnet_pool_id = data.metalcloud_subnet_pool.wan
		subnet_automatic_allocation = false
		subnet_is_ip_range = true
		subnet_ip_range_ip_count = 5
		subnet_override_vlan_id=1003
		
}

resource "metalcloud_eksa" "cluster01" {
	infrastructure_id =  data.metalcloud_infrastructure.infra.infrastructure_id
	
	cluster_label = "test-eksa"

	
	instance_array_instance_count_eksa_mgmt = 1
	instance_array_instance_count_mgmt = 1
	instance_array_instance_count_worker = 1
	
	instance_server_type_eksa_mgmt {
		instance_index = 0
		server_type_id = data.metalcloud_server_type.large.server_type_id
	}
	
	instance_server_type_mgmt {
		instance_index = 0
		server_type_id = data.metalcloud_server_type.large.server_type_id
	}
	
	instance_server_type_worker {
		instance_index = 0
		server_type_id = data.metalcloud_server_type.large.server_type_id
	}
	
	
	interface_eksa_mgmt{
		interface_index = 0
		network_id = metalcloud_network.wan.id
	}
	
	interface_eksa_mgmt{
		interface_index = 1
		network_id = metalcloud_network.san.id
	}
	
	
	interface_mgmt{
		interface_index = 0
		network_id = metalcloud_network.wan.id
	}
	
	interface_mgmt {
		interface_index = 1
		network_id = metalcloud_network.san.id
	}
	
	interface_worker {
		interface_index = 0
		network_id = metalcloud_network.wan.id
	}
	
	interface_worker {
		interface_index = 1
		network_id = metalcloud_network.san.id
	}
	
	instance_array_network_profile_eksa_mgmt {
		network_id = metalcloud_network.wan.id
		network_profile_id = data.metalcloud_network_profile.eksa-mgmt.id
	}

	instance_array_network_profile_mgmt {
		network_id = metalcloud_network.wan.id
		network_profile_id = data.metalcloud_network_profile.eksa-control-plane.id
	}

	instance_array_network_profile_worker {
		network_id = metalcloud_network.wan.id
		network_profile_id = data.metalcloud_network_profile.eksa-workload.id
	}
}

```
## Argument Reference

* `cluster_label` (Required) *  **Cluster** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `instance_array_instance_count_eksa_mgmt` The number of instances in the eks_mgmt instance array.
* `instance_array_instance_count_mgmt` The number of instances in the mgmt instance array.
* `instance_array_instance_count_worker` The number of instances in the worker instance array.
* `instance_server_type_eksa_mgmt` The id of the server type to use for eks mgmt nodes for each instance (see example above)
* `instance_server_type_mgmt` The id of the server type to use for mgmt nodes for each instance (see example above)
* `instance_server_type_worker` The id of the server type to use for worker nodes for each instance (see example above)
* `interface_eksa_mgmt` The interface mapping to a network. (see example above)
* `interface_mgmt` The interface mapping to a network. (see example above)
* `interface_worker` The interface mapping to a network. (see example above)
* `instance_array_custom_variables_eks_mgmt` The instance array custom variables.
* `instance_array_custom_variables_mgmt` The instance array custom variables.
* `instance_array_custom_variables_worker` The instance array custom variables.
* `instance_custom_variables_eks_mgmt` instance level custom variables. (see example above)
* `instance_custom_variables_mgmt` instance level custom variables. (see example above)
* `instance_custom_variables_worker` instance level custom variables. (see example above)


## Expanding the EKS-A cluster

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
