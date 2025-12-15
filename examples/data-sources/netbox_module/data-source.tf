data "netbox_module" "test" {
  device_id     = 123
  module_bay_id = 456
}

output "example" {
  value = data.netbox_module.test.id
}
