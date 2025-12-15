data "netbox_tunnel_group" "test" {
  name = "test-tunnel-group"
}

output "example" {
  value = data.netbox_tunnel_group.test.id
}
