# Look up tunnel by ID
data "netbox_tunnel" "by_id" {
  id = "1"
}

# Look up tunnel by name
data "netbox_tunnel" "by_name" {
  name = "example-tunnel"
}

output "tunnel_by_id" {
  value = data.netbox_tunnel.by_id.status
}

output "tunnel_by_name" {
  value = data.netbox_tunnel.by_name.encapsulation
}
