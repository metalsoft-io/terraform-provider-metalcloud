---
page_title: "metalcloud_extension_instance Resource - terraform-provider-metalcloud"
description: |-
  Extension Instance resource for adding an instance of a MetalCloud extension to an infrastructure
---

# metalcloud_extension_instance (Resource)

Creates an instance of a MetalCloud **extension** inside an infrastructure and sets its input variables. Extensions are defined by the platform administrator; use the [`metalcloud_extension`](../data-sources/extension.md) data source to resolve an extension label to its Id.

Changes only take effect after a deploy. Use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) with a `depends_on` on this resource.

## Example Usage

```hcl
data "metalcloud_extension" "my_app" {
  label = "my-app"
}

data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "my-infra01"
  create_if_missing = true
}

resource "metalcloud_extension_instance" "app_inst1" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  extension_id      = data.metalcloud_extension.my_app.extension_id
  label             = "app_inst1"

  input_variables = [
    {
      label     = "var1"
      value_str = "test1"
    },
    {
      label     = "var2"
      value_int = 5
    },
    {
      label      = "var3"
      value_bool = true
    }
  ]
}
```

## Schema

### Required

- `extension_id` (String) Id of the extension to instantiate. Changing this value forces a new resource.
- `infrastructure_id` (String) Id of the infrastructure the instance belongs to. Changing this value forces a new resource.
- `label` (String) Extension instance label. Changing this value forces a new resource.

### Optional

- `input_variables` (Attributes Set) Input variables passed to the extension. Updated in place (see [below for nested schema](#nestedatt--input_variables)).

### Read-Only

- `extension_instance_id` (String) Id of the extension instance.

<a id="nestedatt--input_variables"></a>
### Nested Schema for `input_variables`

Required:

- `label` (String) Name of the input variable as defined by the extension.

Optional:

- `value_bool` (Boolean) Boolean value. Set exactly one of `value_str`, `value_bool` or `value_int`.
- `value_int` (Number) Integer value.
- `value_str` (String) String value.

## Import

```sh
terraform import metalcloud_extension_instance.app_inst1 <extension_instance_id>
```
