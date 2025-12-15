data "netbox_fhrp_group_assignment" "test" {
  group_id     = 123
  interface_id = 456
}

output "example" {
  value = data.netbox_fhrp_group_assignment.test.id
}
