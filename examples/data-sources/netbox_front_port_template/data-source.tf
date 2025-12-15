data "netbox_front_port_template" "test" {
  name           = "test-front-port-template"
  device_type_id = 123
}

output "example" {
  value = data.netbox_front_port_template.test.id
}
