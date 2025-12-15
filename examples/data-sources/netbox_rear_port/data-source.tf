data "netbox_rear_port" "test" {
  name      = "test-rear-port"
  device_id = 123
}

output "example" {
  value = data.netbox_rear_port.test.id
}
