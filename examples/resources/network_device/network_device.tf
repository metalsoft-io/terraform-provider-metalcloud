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

# The target fabric must already exist. The switches' site is set explicitly
# below (the fabric data source also exposes its site_id, so they can be wired
# to the same site).
data "metalcloud_fabric" "target" {
    site_id = data.metalcloud_site.dc.site_id
    label = "my-fabric"
}

locals {
    # Defaults shared by every switch (equivalent to the YAML 'defaults' block).
    switch_defaults = {
        driver          = "cumulus_linux"
        username        = "cumulus"
        management_port = 22
    }

    # Per-switch inventory (equivalent to the YAML 'switches' list, keyed by
    # identifier_string so for_each produces a stable address per switch).
    switches = {
        "leaf-su00-r0" = {
            position           = "leaf"
            management_address = "192.168.200.11"
            loopback_address   = "10.253.128.1"
            asn                = 4200000000
            tags_map = {
                "nvidia/scalability-unit-id" = "0"
                "nvidia/rail-group-id"       = "0"
            }
        }
        "leaf-su00-r1" = {
            position           = "leaf"
            management_address = "192.168.200.12"
            loopback_address   = "10.253.128.2"
            asn                = 4200000001
        }
        "spine-su00-s0" = {
            position           = "spine"
            management_address = "192.168.200.60"
        }
    }
}

resource "metalcloud_network_device" "switch" {
    for_each = local.switches

    site_id             = data.metalcloud_site.dc.site_id
    identifier_string   = each.key
    driver              = local.switch_defaults.driver
    username            = local.switch_defaults.username
    management_password = var.switch_password
    management_port     = local.switch_defaults.management_port

    position           = each.value.position
    management_address = each.value.management_address
    loopback_address   = try(each.value.loopback_address, null)
    asn                = try(each.value.asn, null)
    tags_map           = try(each.value.tags_map, null)

    # Optional and editable: setting this attaches the switch to the fabric,
    # changing it reassigns it, removing it detaches it. No deploy is triggered;
    # deploy the fabric separately to push the change to the hardware.
    fabric_id = data.metalcloud_fabric.target.fabric_id
}

variable "endpoint" {
    default = ""
}

variable "api_key" {
    default = ""
}

variable "logging" {
    default = "false"
}

variable "site" {
    default = ""
}

variable "switch_password" {
    default   = ""
    sensitive = true
}
