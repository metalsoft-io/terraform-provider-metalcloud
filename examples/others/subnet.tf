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


data "metalcloud_infrastructure" "infra" {
   
		infrastructure_label = "test-subnet"
		datacenter_name = "us-chi-qts01-dc"
	 
		create_if_not_exists = true
}
 
resource metalcloud_subnet subnet01 {
	infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

	subnet_is_ip_range = false
	subnet_prefix_size = 27
	subnet_type = "ipv4"
}
