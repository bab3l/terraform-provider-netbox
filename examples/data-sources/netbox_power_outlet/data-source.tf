data "netbox_power_outlet" "test" {
  name      = "test-power-outlet"
  device_id = 123
}

output "example" {
  value = data.netbox_power_outlet.test.id
}
