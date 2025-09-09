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
