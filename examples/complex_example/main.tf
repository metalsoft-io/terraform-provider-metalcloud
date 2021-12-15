
terraform {
  required_providers {
    metalcloud = {
      source = "metalsoft-io/metalcloud"
      version = "2.1.2"
    }
  }
}

#
# #######################################
#
# Configure the metalcloud provider
#
# #######################################

provider "metalcloud" {
   user_email = var.user_email
   api_key = var.api_key
   endpoint = var.endpoint

}



module "tenancy" {
  source = "./modules/tenancy"
  tenancy_config = var.tenancy_config
}

