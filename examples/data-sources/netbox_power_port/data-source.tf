# Lookup by ID
data "netbox_power_port" "by_id" {
  id = 123
}

output "by_id" {
  value = data.netbox_power_port.by_id.name
}

# Lookup by device_id and name
data "netbox_power_port" "by_device_and_name" {
  device_id = 456
  name      = "PWR1"
}

output "by_device_and_name" {
  value = data.netbox_power_port.by_device_and_name.type
}
