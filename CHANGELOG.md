## 8.0.0 (Unreleased)

UPGRADE COMPATIBILITY NOTES:

* The provider is built for the MetalCloud v8.0 API. Use provider 8.x with MetalCloud v8.0 controllers and pin the version, for example `version = "~> 8.0"`.
* `metalcloud_drive` and `metalcloud_vm_instance_group` are imported by `<infrastructure_id>:<id>`, for example `terraform import metalcloud_drive.data01 1234:5678`. Importing by the plain Id never worked. (MS-10918)
* `raid_type` in `storage_controllers` accepts only `RAID0`, `RAID1`, `RAID5` and `RAID10`. Lowercase values such as `raid1`, which the API rejects, now fail at `terraform plan`. (MS-10915)
* `metalcloud_server_instance_group`, `metalcloud_vm_instance_group` and `metalcloud_infrastructure` read their pending (design) configuration instead of the deployed one. State and `terraform plan` now show what Terraform has written, not what is deployed. The first refresh after upgrading replaces deployed values in state with design values. (MS-10918)

DEPRECATIONS:

* `create_if_missing` on the `metalcloud_infrastructure` data source is deprecated. Data sources are read during `terraform plan`, so it created the infrastructure during `plan` and Terraform never tracked or destroyed it. Use the `metalcloud_infrastructure` resource and import an existing infrastructure instead; the resource documentation explains how to migrate. (MS-6131)

IMPROVEMENTS:

* `metalcloud_server_instance_group`, `metalcloud_vm_instance_group`, `metalcloud_endpoint_instance_group`: add `provision_instance_dns_records` and `provision_load_balancing_dns_record` to `network_connections`, to create DNS records for the instances and a DNS load balancing record for the group on a logical network. (MS-8395)
* `metalcloud_server_instance_group`: add `hostname`, which names the group's DNS load balancing record. (MS-8395)
* `metalcloud_server_type` data source: add `disk_count` and `storage_controllers`, with the RAID types, controller modes and disks of each controller. (MS-8468)
* Deploys started with `await_deploy_finish = true` are awaited for up to 2 hours instead of 30 minutes. (MS-5065)
* `terraform destroy` warns when `metalcloud_infrastructure` has `prevent_deploy = true`, because the infrastructure is then only marked for deletion and stays active until it is deployed. (MS-9909)
* Importing `metalcloud_drive` sets `storage_pool_id`, and importing `metalcloud_vm_instance_group` sets `label`, `vm_type_id` and `os_template_id`, so the first plan after an import does not replace the resource. (MS-10918)
* Importing `metalcloud_infrastructure` sets the deploy options to their defaults, so the first plan after an import shows no changes. (MS-6131)
* Documentation: fix the provider and getting started examples, document deploys and destroy, and document importing an infrastructure and the resources inside it. (MS-339, MS-6131)

BUG FIXES:

* Stop waiting for a deploy when one of its jobs fails, and report the failure, instead of waiting until the timeout. (MS-10320)
* `metalcloud_server_instance_group`: an update now converges. Before, every apply reissued the same update until the infrastructure was deployed. The same applies to `metalcloud_vm_instance_group` (`instance_count`, `label`, `custom_variables`) and to the `metalcloud_infrastructure` label. (MS-10918)
* `metalcloud_server_instance_group`: send `server_type_id` on update. Changing it had no effect. (MS-10916)
* `metalcloud_server_instance_group`: only send an empty storage profile when `storage_controllers` is removed from the configuration, instead of clearing the profile on every update. (MS-10919)
* `metalcloud_server_instance_group`, `metalcloud_vm_instance_group`: an omitted `network_connections` no longer shows a perpetual diff against an empty set. (MS-10917)
* `metalcloud_server_instance_group`, `metalcloud_vm_instance_group`, `metalcloud_endpoint_instance_group`: fix a crash when the API returns a null `interfaceCount` for a network connection. (MS-10914)
* `metalcloud_vm_instance_group`: fix `terraform destroy` failing with `412 Precondition Failed`. The delete sent the revision of the group's design configuration instead of the revision of the group. (MS-8649)
* `metalcloud_server_instance_group`, `metalcloud_vm_instance_group`, `metalcloud_endpoint_instance_group`: fix `Provider returned invalid result object after apply` when a network connection without `interface_count` is created or changed, and stop sending `interfaceCount: 0` for it on update.
* `metalcloud_server_instance_group`: fix a crash on update and delete when the API response has no `ETag` header. The object's revision is sent in `If-Match` instead.
* Fix the Ids sent when updating network connections after the switch to the v8.0 SDK.
