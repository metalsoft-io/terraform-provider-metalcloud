---
layout: "metalcloud"
page_title: "Metalcloud: instance array custom variables set"
description: |-
  Represents a set of custom variables that is applied on an instance array.
---

# metalcloud_infrastructure/instance_array/instance_array_custom_variables



## Example usage

The following is an example of instance array level custom variable set. These are variables that will be applied at the `instance array` level and will override any identical ones configured at the `infrastructure` level specified via the `infrastructure_custom_variables` property.

```hcl
resource "metalcloud_infrastructure" "foo" {
    ...
    
    instance_array {
        ...

        instance_array_custom_variables = {
            b = "c"
            d = "e"
            c = "f"
            r = "p"
        }
}
```


