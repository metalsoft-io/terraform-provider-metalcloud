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
  api_key    = var.api_key
  endpoint   = var.endpoint
}


variable "user_email" {
  type    = string
  default = ""
}

variable "api_key" {
  type    = string
  default = ""
}


variable "endpoint" {
  type    = string
  default = ""
}


variable "datacenter" {
  type    = string
  default = ""
}

data "metalcloud_external_connection" "uplink1" {
  external_connection_label = "fortigate1"
  datacenter_name = var.datacenter
}

data "metalcloud_subnet_pool" "primary" {
  subnet_pool_label = "n6"
}

resource "metalcloud_network_profile" "myprofile" {
    network_profile_label ="test-1"
    datacenter_name = "${var.datacenter}"
    network_type = "wan"

    network_profile_vlan {
        vlan_id = "101"
        port_mode = "trunk"
        provision_subnet_gateways = false
   }

   network_profile_vlan {
        vlan_id = "auto"
        port_mode = "native"
        provision_subnet_gateways = false
        provision_vxlan = true
        subnet_pool_ids = [
          data.metalcloud_subnet_pool.primary.id
          ]
   }

   network_profile_vlan {
        vlan_id = "66"
        port_mode = "trunk"
        provision_subnet_gateways = false
        external_connection_ids = [
          data.metalcloud_external_connection.uplink1.id
        ]
   }
}  
