---
page_title: "metalcloud_endpoint Data Source - terraform-provider-metalcloud"
description: |-
  Use this data source to look up a MetalCloud endpoint (an unmanaged node bound to switch interfaces) by its label.
---

# metalcloud_endpoint (Data Source)

Use this data source to retrieve a MetalCloud **endpoint** by label. Endpoints are unmanaged nodes (for example HGX hosts) bound to switch interfaces. The returned `endpoint_id` is typically fed into [`metalcloud_endpoint_instance_group`](../resources/endpoint_instance_group.md) to attach the endpoint to a logical network.

## Example Usage

```hcl
# Look up an endpoint by label
data "metalcloud_endpoint" "hgx_h08" {
  label = "hgx-su00-h08"
}

output "hgx_h08_id" {
  value = data.metalcloud_endpoint.hgx_h08.endpoint_id
}
```

## Argument Reference

### Required

- `label` (String) The label of the endpoint to retrieve.

### Optional

- `site_id` (String) The site the endpoint belongs to. Provide it to narrow the search when labels are only unique within a site.

## Attributes Reference

- `endpoint_id` (String) The endpoint's Id.
- `name` (String) The endpoint's name.
- `site_id` (String) The site the endpoint belongs to (computed when not supplied).

## Notes

- The endpoints listing has no server-side label filter, so the lookup matches the label client-side. Set `site_id` to limit the search scope.

## Related Resources

- [`metalcloud_endpoint_instance_group`](../resources/endpoint_instance_group.md) - Attach endpoints to logical networks
- [`metalcloud_logical_network`](../resources/logical_network.md) - Manage logical networks
