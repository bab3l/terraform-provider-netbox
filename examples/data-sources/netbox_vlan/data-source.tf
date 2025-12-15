data "netbox_vlan" "test" {
  vid = 100
}

output "example" {
  value = data.netbox_vlan.test.id
}
