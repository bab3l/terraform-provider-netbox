data "netbox_interface_template" "test" {
  name           = "test-interface-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_interface_template.test.id
}
