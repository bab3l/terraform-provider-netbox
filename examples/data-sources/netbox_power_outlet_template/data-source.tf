data "netbox_power_outlet_template" "test" {
  name           = "test-power-outlet-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_power_outlet_template.test.id
}
