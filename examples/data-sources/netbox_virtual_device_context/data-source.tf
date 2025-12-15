data "netbox_virtual_device_context" "test" {
  name      = "test-vdc"
  device_id = 123
}

output "example" {
  value = data.netbox_virtual_device_context.test.id
}
