data "netbox_front_port" "test" {
  name      = "test-front-port"
  device_id = 123
}

output "example" {
  value = data.netbox_front_port.test.id
}
