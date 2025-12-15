data "netbox_power_port_template" "test" {
  name           = "test-power-port-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_power_port_template.test.id
}
