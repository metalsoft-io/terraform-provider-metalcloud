---
layout: "metalcloud"
page_title: "Metalcloud: instance custom variables set"
description: |-
  Represents a set of custom variables that is applied on an specific instance of an instance array.
---

# instance_array/instance_custom_variables


## Example usage

These are variables that will be applied at the **instance** level and will override any identical ones configured at the **infrastructure** and **instance_array** level via the `infrastructure_custom_variables` and `instance_array_custom_variables` properties. Use the `instance_index` property to specify which from the instance array's instances this set of variables applies to. For example the variables for the second instance of an array would be:

```hcl
resource "metalcloud_instance_array" "cluster" {
    ...

    instance_custom_variables {
      instance_index = 1
      custom_variables={
        "test1":"test2"
        "test3":"test4"
      }
    }
}
```


