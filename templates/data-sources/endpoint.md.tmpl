---
page_title: "metalcloud_endpoint Data Source - terraform-provider-metalcloud"
description: |-
  Use this data source to look up a MetalCloud endpoint (an unmanaged node bound to switch interfaces) by its label.
---

# metalcloud_endpoint (Data Source)

Retrieves a MetalCloud **endpoint** by label. Endpoints are unmanaged nodes (for example HGX hosts) bound to switch interfaces. The returned `endpoint_id` is typically passed to [`metalcloud_endpoint_instance_group`](../resources/endpoint_instance_group.md) to attach the endpoint to a logical network.

## Example Usage

```hcl
data "metalcloud_endpoint" "hgx_h08" {
  label = "hgx-su00-h08"
}

output "hgx_h08_id" {
  value = data.metalcloud_endpoint.hgx_h08.endpoint_id
}
```

## Schema

### Required

- `label` (String) Label of the endpoint.

### Optional

- `site_id` (String) Id of the site the endpoint belongs to. Set it to narrow the search when labels are only unique within a site. Computed when omitted.

### Read-Only

- `endpoint_id` (String) Id of the endpoint.
- `name` (String) Name of the endpoint.

## Notes

- The endpoint listing has no server-side label filter, so the label is matched client-side. Set `site_id` to limit the search scope.
