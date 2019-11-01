---
layout: "metalcloud"
page_title: "Template: volume_template"
description: |-
  Provides a mechanism to search for template ids.
---

# volume_template

This data source provides a mechanism to identify the ID of a volume template based on it's name.


## Example usage

The following example locates the volume_template_ID for 'Cenots 7.6'.

```hcl
data "metalcloud_volume_template" "centos76" {
			volume_template_label = "centos7-6"
}
resource "metalcloud_infrastructure" "foo" {

  instance_array {

    volume_template_id = tonumber(data.metalcloud_volume_template.centos76.id)
  }
}
```

## Arguments

`volume_template_label` (Required) String used to locate the template. Values such as centos7-7, rhel7-6 etc. are permitted. If the provided name does not mach any valid templates,a list of possible templates is returned in the error message.