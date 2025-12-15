data "netbox_rear_port_template" "test" {
  name           = "test-rear-port-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_rear_port_template.test.id
}
