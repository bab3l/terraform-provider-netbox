data "netbox_console_server_port" "test" {
  name      = "test-console-server-port"
  device_id = 123
}

output "example" {
  value = data.netbox_console_server_port.test.id
}
