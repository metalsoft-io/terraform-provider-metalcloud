---
layout: "metalcloud"
page_title: "Metalcloud: firmware_policy"
description: |-
  Controls a Metalcloud Firmware policy
---


# metalcloud_firmware_policy

This resource allows users to request certain firmware for a component to be present on the hardware. The process will start after servers will be allocated to a deploy (or they might have already been allocated. Then the components are matched against the rules in all configured policies. 

## Example usage

```hcl

data "metalcloud_server_type" "large" {

  server_type_name="M.40.32.v2"
  
}

resource "metalcloud_firmware_policy" "upgrade-raid-controller" {
  server_firmware_upgrade_policy_label = "upgrade-by-component-name-to-specific-version"

  //Possible values: accept, deny, accept_with_confirmation
  server_firmware_upgrade_policy_action = "accept"

  server_firmware_upgrade_policy_rule {
    operation = "string_equal"
    property = "server_type_id"
    value = data.metalcloud_server_type.large.server_type_id
  }

  server_firmware_upgrade_policy_rule {
    operation = "string_contains"
    property  = "server_component_name"
    value     = "PERC H330 Adapter"
  }
  server_firmware_upgrade_policy_rule {
    operation = "string_equal"
    property  = "server_component_target_version"
    value     = "25.5.9.0001"
  }

  instance_array_list = [metalcloud_instance_array.cluster.instance_array_id]
}

```
## Argument Reference

* `server_firmware_upgrade_policy_label` (Required) *  **Policy** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `server_firmware_upgrade_policy_action` (Required) Possible values: `accept`, `reject`,`accept_with_confirmation`. 
* `instance_array_list` (Optional, default: 40960) The list of instance array ids to which this policy applies
* `server_firmware_upgrade_policy_rule` (Required, default: []) An array of policy rules such as:
  ```
    server_firmware_upgrade_policy_rule {
        operation = "string_equal"
        property = "server_type_id"
        value = "1"
    }

    server_firmware_upgrade_policy_rule {
      operation = "string_contains"
      property = "server_component_name"
      value = "BIOS"
    }
  ```
  Which should be interpreted as a series of logic tests `<property> <operation> <value>` joined by `AND`. All of the rules need to match before the action of the policy is executed.

  Where: 
    * `operation` is one of:
      * `string_equal`
      * `string_contains`
    * `property` is one of:
      * `server_component_name`
      * `server_component_type`
      * `server_component_firmware_version`
      * `datacenter_name`
      * `server_vendor`
      * `server_id` (only string_equal operation is supported)
      * `server_tags_json`
      * `server_type_id` (only string_equal operation is supported)
      * `server_component_target_version`. This is a special rule that will instruct the system to set a particular version on the component rather than the latest available. It is not a condition.
    * `value` is the left hand operator, defined as a string value
 
  Work with your service provider to get a list of valid component names. This list depends on the hardware vendor and generation used, as does the firmware version strings. 

> Note that firmware downgrades might not be supported by all components. Check with your hardware vendor. 

For more information visit the [MetalSoft Documentation](https://docs.metalsoft.io/en/latest/advanced/managing_firmware.html)