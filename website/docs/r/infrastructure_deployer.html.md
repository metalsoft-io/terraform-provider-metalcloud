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

  # Set this to false to trigger deploys.
  prevent_deploy = true

  # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  # instead of exiting while the deploy is ongoing
  await_deploy_finished = true
  await_delete_finished = true

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
* `await_delete_finished` (Optional, default false) - If **true**, the provider will wait for a deploy (involving delete) to finish before exiting. If **false**, the delete operation (really a deploy) will continue after the provider existed. This operation is generally quick.
* `keep_detaching_drives` (Optional, default true) - If **true**, the detaching Drive objects will not be deleted. If **false**, and the number of Instance objects is reduced, then the detaching Drive objects will be deleted.
* `infrastructure_custom_variables` (Optional, default []) - All of the variables specified as a map of *string* = *string* such as { var_a="var_value" } will be sent to the underlying deploy process and referenced in operating system templates and workflows. 


## Attributes

This resource exports the following attributes:

* `infrastructure_id` - The id of the infrastructure is used for many operations. It is also the ID of the resource object.
