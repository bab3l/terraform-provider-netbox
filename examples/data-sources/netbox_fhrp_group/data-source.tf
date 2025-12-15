data "netbox_fhrp_group" "test" {
  group_id = 10
}

output "example" {
  value = data.netbox_fhrp_group.test.id
}
