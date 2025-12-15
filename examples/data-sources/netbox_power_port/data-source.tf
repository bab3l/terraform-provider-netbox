data "netbox_power_port" "test" {
  name      = "test-power-port"
  device_id = 123
}

output "example" {
  value = data.netbox_power_port.test.id
}
