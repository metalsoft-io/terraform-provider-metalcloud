terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
    }
  }
}

provider "metalcloud" {
  endpoint = var.endpoint
  api_key  = var.api_key
  logging  = var.logging
}

data "metalcloud_site" "dc" {
  label = var.site
}

data "metalcloud_fabric" "wan" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "wan-fabric"
}

# An existing L3 logical-network profile (encodes kind/route-domain/L3VNI).
data "metalcloud_logical_network_profile" "l3" {
  label     = "tenant-l3"
  fabric_id = data.metalcloud_fabric.wan.fabric_id
}

data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "tenant1-infra"
  create_if_missing = true
}

# The tenant L3 logical network the endpoints attach to.
resource "metalcloud_logical_network" "tenant_l3" {
  infrastructure_id          = data.metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.l3.logical_network_profile_id

  name  = "tenant1-l3"
  label = "tenant1-l3"
}

# Select the endpoints (HGX nodes) to attach, by label.
data "metalcloud_endpoint" "hgx" {
  for_each = toset(["hgx-su00-h08", "hgx-su00-h24"])
  label    = each.key
}

# Attach the selected endpoints to the logical network.
resource "metalcloud_endpoint_instance_group" "hgx_hosts" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  label             = "hgx-hosts"

  endpoint_ids = [for e in data.metalcloud_endpoint.hgx : e.endpoint_id]

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.tenant_l3.logical_network_id
      tagged             = false
      access_mode        = "l3"
      mtu                = 9000
    }
  ]

  depends_on = [metalcloud_logical_network.tenant_l3]
}

# Trigger the deploy after the attachment is in place.
resource "metalcloud_infrastructure_deployer" "deploy" {
  infrastructure_id   = data.metalcloud_infrastructure.infra.infrastructure_id
  prevent_deploy      = true
  await_deploy_finish = false
  allow_data_loss     = true

  depends_on = [metalcloud_endpoint_instance_group.hgx_hosts]
}

variable "endpoint" { default = "" }
variable "api_key" { default = "" }
variable "logging" { default = "false" }
variable "site" { default = "" }
