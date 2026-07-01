---
page_title: "metalcloud_endpoint_instance_group Resource - terraform-provider-metalcloud"
description: |-
  Attaches selected endpoints to one or more logical networks within a MetalCloud infrastructure.
---

# metalcloud_endpoint_instance_group (Resource)

Attaches a set of **endpoints** (unmanaged nodes such as HGX hosts) to one or more **logical networks**. The resource creates an endpoint instance group in the given infrastructure, adds each selected endpoint to it as an endpoint instance, and connects the group to the logical network(s).

A deploy is required afterwards to apply the attachment — use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) with a `depends_on` this resource.

## Example Usage

```hcl
data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "tenant1-infra"
  create_if_missing = true
}

resource "metalcloud_logical_network" "tenant_l3" {
  infrastructure_id          = data.metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.l3.logical_network_profile_id
  name                       = "tenant1-l3"
  label                      = "tenant1-l3"
}

# Select the endpoints to attach, by label.
data "metalcloud_endpoint" "hgx" {
  for_each = toset(["hgx-su00-h08", "hgx-su00-h24"])
  label    = each.key
}

resource "metalcloud_endpoint_instance_group" "hgx_hosts" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  label             = "hgx-hosts"

  endpoint_ids = [for e in data.metalcloud_endpoint.hgx : e.endpoint_id]

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.tenant_l3.logical_network_id
      tagged             = false
      access_mode        = "l3"
      mtu                = 9000
    }
  ]

  depends_on = [metalcloud_logical_network.tenant_l3]
}
```

## Argument Reference

### Required

- `infrastructure_id` (String) The infrastructure the endpoint instance group belongs to. Changing this forces a new resource.
- `endpoint_ids` (Set of String) Ids of the endpoints to attach. Each becomes an endpoint instance in the group. Use the [`metalcloud_endpoint`](../data-sources/endpoint.md) data source to resolve labels to ids.

### Optional

- `label` (String) The endpoint instance group label. Assigned by the platform if omitted. Changing this forces a new resource.
- `network_connections` (Attributes List) The logical networks this group of endpoints connects to. See below.

### Nested Schema for `network_connections`

- `logical_network_id` (String, Required) The logical network to connect to.
- `tagged` (Boolean, Required) Whether the connection is VLAN-tagged.
- `access_mode` (String, Required) The access mode (e.g. `l2`, `l3`).
- `mtu` (Number, Optional) MTU for the connection (defaults to 1500 when omitted).

## Attributes Reference

- `endpoint_instance_group_id` (String) The Id of the created endpoint instance group.

## Import

```sh
terraform import metalcloud_endpoint_instance_group.hgx_hosts <endpoint_instance_group_id>
```

## Related Resources

- [`metalcloud_endpoint`](../data-sources/endpoint.md) - Look up endpoints by label
- [`metalcloud_logical_network`](logical_network.md) - Manage logical networks
- [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) - Trigger deploys
