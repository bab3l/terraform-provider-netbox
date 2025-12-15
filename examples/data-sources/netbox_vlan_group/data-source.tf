data "netbox_vlan_group" "test" {
  name = "test-vlan-group"
}

output "example" {
  value = data.netbox_vlan_group.test.id
}
