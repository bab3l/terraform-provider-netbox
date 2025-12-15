data "netbox_virtual_machine" "test" {
  name = "test-vm"
}

output "example" {
  value = data.netbox_virtual_machine.test.id
}
