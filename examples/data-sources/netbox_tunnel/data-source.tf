data "netbox_tunnel" "test" {
  name = "test-tunnel"
}

output "example" {
  value = data.netbox_tunnel.test.id
}
