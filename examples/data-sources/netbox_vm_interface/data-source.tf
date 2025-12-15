data "netbox_vm_interface" "test" {
  name               = "eth0"
  virtual_machine_id = 123
}

output "example" {
  value = data.netbox_vm_interface.test.id
}
