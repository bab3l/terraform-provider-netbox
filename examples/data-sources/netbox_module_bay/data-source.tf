data "netbox_module_bay" "test" {
  name      = "test-module-bay"
  device_id = 123
}

output "example" {
  value = data.netbox_module_bay.test.id
}
