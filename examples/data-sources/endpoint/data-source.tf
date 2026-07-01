# Look up an endpoint (an unmanaged node bound to switch interfaces) by label.
data "metalcloud_endpoint" "hgx_h08" {
  label = "hgx-su00-h08"
}

output "hgx_h08_id" {
  value = data.metalcloud_endpoint.hgx_h08.endpoint_id
}
