data "netbox_console_port_template" "test" {
  name           = "test-console-port-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_console_port_template.test.id
}
