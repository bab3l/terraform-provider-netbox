data "netbox_tunnel_termination" "test" {
  tunnel_id = 123
  role      = "peer"
}

output "example" {
  value = data.netbox_tunnel_termination.test.id
}
