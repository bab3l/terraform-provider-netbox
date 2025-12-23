data "netbox_console_port" "by_id" {
  id = "789"
}

data "netbox_console_port" "by_device_and_name" {
  device_id = "456"
  name      = "con0"
}

output "by_id" {
  value = data.netbox_console_port.by_id.name
}

output "by_device_and_name" {
  value = data.netbox_console_port.by_device_and_name.id
}
