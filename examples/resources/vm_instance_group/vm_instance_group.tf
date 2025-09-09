terraform {
    required_providers {
        metalcloud = {
            source = "metalsoft-io/metalcloud"
        }
    }
}

check "vars" {
    assert {
        condition = var.site != ""
        error_message = "site variable cannot be empty"
    }
}

provider "metalcloud" {
    endpoint = var.endpoint
    api_key = var.api_key
    logging = var.logging
}

data "metalcloud_site" "dc" {
    label = "${var.site}"
}

data "metalcloud_fabric" "wan" {
    site_id = data.metalcloud_site.dc.site_id
    label = "wan-fabric"
}

data "metalcloud_logical_network_profile" "np01" {
    label = "np-01"
    fabric_id = data.metalcloud_fabric.wan.fabric_id
}

data "metalcloud_vm_type" "vm1" {
    label = "vm-medium"
}

data "metalcloud_os_template" "vmos1" {
    label = "vm-ubuntu-22-04-cloud"
}

data "metalcloud_infrastructure" "infra" {
    site_id = data.metalcloud_site.dc.site_id
    label = "my-infra01"

    create_if_missing = true
}

resource "metalcloud_logical_network" "net1" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
    logical_network_profile_id = data.metalcloud_logical_network_profile.np01.logical_network_profile_id

    name = "net01"
    label = "net01"
}

resource "metalcloud_vm_instance_group" "vminst01" {
    infrastructure_id =	 data.metalcloud_infrastructure.infra.infrastructure_id

    label = "vminst01"

    instance_count = 1
    vm_type_id = data.metalcloud_vm_type.vm1.vm_type_id
    disk_size_gbytes = 10
    os_template_id = data.metalcloud_os_template.vmos1.os_template_id

    network_connections = [
        {
            logical_network_id = metalcloud_logical_network.net1.logical_network_id
            tagged = true
            access_mode = "l2"
            mtu = 1500
        }
    ]

    custom_variables = [
        {
            name = "vkey1"
            value = "vtest1"
        },
        {
            name = "vkey2"
            value = "vtest2"
        }
    ]

    depends_on = [
        metalcloud_logical_network.net1,
    ]
}

# Use this resource to effect deploys of the above resources.
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    # Set this to false to actually trigger deploys.
    prevent_deploy = true

    # These options will make terraform apply operation will wait for the deploy to finish (when prevent_deploy is false)
    # instead of exiting while the deploy is ongoing
    await_deploy_finish = false

    # This option disables a safety check that metalsoft performs to prevent accidental data loss
    # It is required when testing delete operations
    allow_data_loss = true

    # IMPORTANT. This is important to ensure that deploys happen after everything else. If you need to add or remove resources dynamically
    # use either count or for_each in the resources or move everything that is dynamic into a module and make this depend on the module
    depends_on = [
        metalcloud_vm_instance_group.vminst01,
    ]
}

variable "endpoint" {
    default =""
}

variable "api_key" {
    default = ""
}

variable "logging" {
    default="false"
}

variable "site" {
    default=""
}
