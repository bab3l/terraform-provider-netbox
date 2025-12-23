data "netbox_virtual_machine" "by_name" {
  name = "test-vm"
}

data "netbox_virtual_machine" "by_id" {
  id = 123
}

output "by_name" {
  value = data.netbox_virtual_machine.by_name.id
}

output "by_id" {
  value = data.netbox_virtual_machine.by_id.name
}
