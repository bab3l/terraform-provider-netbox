data "netbox_console_server_port_template" "test" {
  name           = "test-console-server-port-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_console_server_port_template.test.id
}
