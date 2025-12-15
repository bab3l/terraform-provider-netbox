data "netbox_module_bay_template" "test" {
  name           = "test-module-bay-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_module_bay_template.test.id
}
