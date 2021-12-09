---
layout: "metalcloud"
page_title: "Template: infrastructure_reference"
description: |-
  Provides a reference to a infrastructure. 
---

# infrastructure_reference

This data source provides a mechanism to determine the `infrastructure_id` of an *Infrastructure* and to ensure that it is created if it does not exist.


## Example usage

The following example locates the volume_template_ID for 'Cenots 7.6'.

```hcl
data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test" 
    datacenter_name = "dc-1" 
    create_if_not_exists = true

}
```

## Arguments

* `infrastructure_label` - (Required) **Infrastructure** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `datacenter_name` - (Required) The name of the **Datacenter** where the provisioning will take place. Check the MetalCloud provider for available options.
* `create_if_not_exist` - (Optional) If set to true it will create the infrastructure if it does not exist. Defaults to `true`.

## Attributes

This resource exports the following attributes:

* `infrastructure_id` - The id of the infrastructure is used for many operations. It is also the ID of the data object.
