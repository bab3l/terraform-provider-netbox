data "netbox_virtual_chassis" "test" {
  name = "test-virtual-chassis"
}

output "example" {
  value = data.netbox_virtual_chassis.test.id
}
