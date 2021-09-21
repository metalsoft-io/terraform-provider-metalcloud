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
# List required providers
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
      version = "1.0.12"
    }
  }
}

# Configure the metalcloud provider
provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key 
   endpoint = var.endpoint
}

# Identity the ID of the volume template we want
data "metalcloud_volume_template" "centos76" {
  volume_template_label = "centos7-6"
}


resource "metalcloud_infrastructure" "my-infra92" {
  
  infrastructure_label = "my-terraform-infra92"
  datacenter_name = var.datacenter

  # Set this to false to actually deploy the changes, otherwise all changes will remain in edit mode only.
  prevent_deploy = true 

  shared_drive {
    shared_drive_label = "my-shared-drive"
    shared_drive_size_mbytes = 40965
    shared_drive_storage_type = "iscsi_ssd"
    shared_drive_attached_instance_arrays = ["web-servers","web-servers-2"]
  }
  
  instance_array {
    # Name of your cluster. Needs to obey DNS rules as it will translate into a DNS record.
    instance_array_label = "web-servers"

    instance_array_instance_count = 1
    instance_array_ram_gbytes = 2
    instance_array_processor_count = 1
    instance_array_processor_core_count = 2

    drive_array{
      drive_array_label = "web-servers-centos"
      drive_array_storage_type = "iscsi_hdd"

      # The size of the drive array in MBytes
      drive_size_mbytes_default = 49000

      # The id of the template we located earlier
      volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
    }

    #one or more FW rules. By default all traffic is denied so we need at least one entry.
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
    # Name of your cluster. Needs to obey DNS rules as it will translate into a DNS record.
    instance_array_label = "web-servers-2"

    instance_array_instance_count = 1
    instance_array_ram_gbytes = 2
    instance_array_processor_count = 1
    instance_array_processor_core_count = 2

    drive_array{
      drive_array_label = "web-servers-centos-2"
      drive_array_storage_type = "iscsi_hdd"

      # The size of the drive array in MBytes
      drive_size_mbytes_default = 49000

      # The id of the template we located earlier
      volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
    }

    #one or more FW rules. By default all traffic is denied so we need at least one entry.
    firewall_rule {
      firewall_rule_description = "test fw rule-2"
      firewall_rule_port_range_start = 22
      firewall_rule_port_range_end = 22
      firewall_rule_source_ip_address_range_start="0.0.0.0"
      firewall_rule_source_ip_address_range_end="0.0.0.0"
      firewall_rule_protocol="tcp"
      firewall_rule_ip_address_type="ipv4"
    }
  }
}
```

To deploy this infrastructure export the following variables (or use -var):

```bash
export TF_VAR_api_key="<yourkey>"
export TF_VAR_user_email="test@test.com"
export TF_VAR_endpoint="https://api.bigstep.com/metal-cloud"
export TF_VAR_datacenter="uk-reading"
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
