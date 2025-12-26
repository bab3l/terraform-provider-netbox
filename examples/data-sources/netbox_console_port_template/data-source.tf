data "netbox_console_port_template" "by_id" {
  id = "123"
}

data "netbox_console_port_template" "by_device_type_and_name" {
  device_type = "456"
  name        = "Con0"
}

output "by_id" {
  value = data.netbox_console_port_template.by_id.name
}

output "by_device_type_and_name" {
  value = data.netbox_console_port_template.by_device_type_and_name.id
}
