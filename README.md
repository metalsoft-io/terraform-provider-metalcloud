Bigstep Metalcloud Terraform Provider
==================
This is a terraform plugin for controlling Bigstep Metalcloud resources.

[![Build Status](https://travis-ci.org/bigstepinc/terraform-provider-metalcloud.svg?branch=master)](https://travis-ci.org/bigstepinc/terraform-provider-metalcloud)

Maintainers
-----------

This provider plugin is maintained by the Bigstep Team.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)


Installing The Provider
----------------------

While the provider is in alpha stage, it will need to be manually installed by:
```bash
go install github.com/bigstepinc/terraform-provider-metalcloud
cp $GOPATH/bin/terraform-provider-metalcloud ~/.terraform.d/plugins
terraform init
```

Using the Provider
------------------
A terraform `main.tf` template file, for an infrastructure with a single server would look something like this:

```terraform
provider "metalcloud" {
   user = var.user
   api_key = var.api_key 
   endpoint = var.endpoint
}

data "metalcloud_volume_template" "centos75" {
  volume_template_label = "centos7-5"
}

resource "metalcloud_infrastructure" "my-infra97" {
  
  infrastructure_label = "my-terraform-infra97"
  datacenter_name = "us-santaclara"

  instance_array {
        instance_array_label = "my-instance-array"
        instance_array_instance_count = 1

        instance_array_ram_gbytes = 8
        instance_array_processor_count = 1
        instance_array_processor_core_count = 8

        drive_array{
          drive_array_storage_type = "iscsi_hdd"
          drive_size_mbytes_default = 49000
          volume_template_id = tonumber(data.metalcloud_volume_template.centos75.id)
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
}
```

To deploy this infrastructure export the following variables (or use -var):

```bash
export TF_VAR_api_key="<yourkey>"
export TF_VAR_user="test@test.com"
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

Clone repository to: `$GOPATH/src/github.com/bigstepinc/terraform-provider-metalcloud`

```sh
$ mkdir -p $GOPATH/src/github.com/bigstepinc; cd $GOPATH/src/github.com/bigstepinc
$ git clone git@github.com/bigstepinc/terraform-provider-metalcloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/bigstepinc/terraform-provider-metalcloud
$ make build
```
Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-metalcloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
export METALCLOUD_DATACENTER="uk-reading"
export METALCLOUD_API_KEY="<api-key>"
export METALCLOUD_USER="user"
export METALCLOUD_ENDPOINT="https://your-endpoint"

make testacc
```
