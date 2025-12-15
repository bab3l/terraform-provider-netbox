data "netbox_console_port" "test" {
  name      = "test-console-port"
  device_id = 123
}

output "example" {
  value = data.netbox_console_port.test.id
}
