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

data "metalcloud_extension" "my-app" {
    label = "my-app"
}

data "metalcloud_infrastructure" "infra" {
    site_id = data.metalcloud_site.dc.site_id
    label = "my-infra01"

    create_if_missing = true
}

resource "metalcloud_extension_instance" "app_inst1" {
    infrastructure_id =	 data.metalcloud_infrastructure.infra.infrastructure_id
    extension_id = data.metalcloud_extension.my-app.extension_id

    label = "app_inst1"

    input_variables = [
        {
            label = "var1"
            value_string = "test1"
        },
        {
            label = "var2"
            value_int = 5
        }
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
        metalcloud_extension_instance.app_inst1,
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
