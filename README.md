MetalSoft Terraform Provider
==================
This is a terraform plugin for controlling Metalcloud resources.

Maintainers
-----------

This provider plugin is maintained by the MetalSoft Team.

Using the Provider
------------------
A terraform `main.tf` template file, for an infrastructure with a single server would look something like this:

```terraform
/* Simple example of using metalcloud */
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}

provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key
   endpoint = var.endpoint

}

//this is an infrastructure reference. It is needed to avoid a cyclic dependency where the 
//infrastructure depends on the resoruces and vice-versa. This will create the infrastructure if it does not exist
//if the create_if_not_exists flag is set to true
data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "${var.datacenter}" 

    create_if_not_exists = true
}

data "metalcloud_volume_template" "esxi7" {
  volume_template_label = "esxi-700-uefi-v2"
}

resource "metalcloud_instance_array" "cluster" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    instance_array_label = "test-2"

    instance_array_instance_count = 1 //deprecated, keep equal to 1
    instance_array_ram_gbytes = "16"
    instance_array_processor_count = 1
    instance_array_processor_core_count = 1
    instance_array_boot_method = "local_drives"

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

}

resource "metalcloud_shared_drive" "datastore" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  
    shared_drive_label = "test-da"
    shared_drive_size_mbytes = 40966
    shared_drive_storage_type = "iscsi_hdd"

    shared_drive_attached_instance_arrays = [metalcloud_instance_array.cluster.instance_array_id]  //this will create a dependency on the instance array
}

resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  # Set this to false to actually trigger deploys.
  prevent_deploy = true

  #these options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
  #instead of exiting while the deploy is ongoing

  await_deploy_finished = true

  #this option disables a safety check that metalsoft performs to prevent accidental data loss
  #it is required when testing delete operations

  allow_data_loss = true

  //This is important to ensure that deploys happen after everything else. If you need to add or remove resources dinamically
  //use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
  depends_on = [
    metalcloud_instance_array.cluster,
    metalcloud_shared_drive.datastore
  ]

}
```

To deploy this infrastructure export the following variables (or use -var):

```bash
export TF_VAR_api_key="<yourkey>"
export TF_VAR_user_email="<your user email>"
export TF_VAR_endpoint="<your endpoint>"
export TF_VAR_datacenter="<your datacenter>"
```

The plan phase:
```bash
terraform plan
```

The apply phase:
```bash
terraform apply
```

To delete the infrastrucure:
```bash
terraform destroy
```

Building The Provider
---------------------

To build the provider:

```sh
make
```

To install the provider so that `terraform init` works:
```sh
make install
```

Enter the provider directory and build the provider

Testing the Provider
---------------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
export METALCLOUD_DATACENTER="uk-reading"
export METALCLOUD_API_KEY="<api-key>"
export METALCLOUD_USER_EMAIL="user"
export METALCLOUD_ENDPOINT="https://your-endpoint"

make testacc
```

Troubleshooting
```
export METALCLOUD_LOGGING_ENABLED=true 
```