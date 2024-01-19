---
layout: "metalcloud"
page_title: "Template: workflow_task"
description: |-
  Provides a mechanism to search for workflow task ids.
---

# volume_template

This data source provides a mechanism to identify the ID of a volume template based on it's name.


## Example usage

The following example locates the volume_template_ID for 'Cenots 7.6'.

```hcl
data "metalcloud_workflow_task" "PowerFlex" {
    stage_definition_label = "deploy-pf-hci"
}

//example usage
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {

    ...

    workflow_task {
      stage_definition_id = data.metalcloud_workflow_task.PowerFlex.id
      run_level = 0
      stage_run_group = "post_deploy"
    }
    ...
}
```

## Arguments

`stage_definition_label` (Required) String used to locate the workflow task.


## Attributes

This resource exports the following attributes:

* `stage_definition_id` - The id of the task.
* `id` - Same as `stage_definition_id`
