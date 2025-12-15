data "netbox_l2vpn_termination" "test" {
  l2vpn_id           = 123
  assigned_object_id = 456
}

output "example" {
  value = data.netbox_l2vpn_termination.test.id
}
