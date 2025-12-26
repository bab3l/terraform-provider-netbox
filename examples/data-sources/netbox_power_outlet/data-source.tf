# Example 1: Lookup by ID
data "netbox_power_outlet" "by_id" {
  id = 1
}

output "power_outlet_by_id" {
  value       = data.netbox_power_outlet.by_id.name
  description = "Power outlet name when looked up by ID"
}

# Example 2: Lookup by device_id and name
data "netbox_power_outlet" "by_device_and_name" {
  device_id = 5
  name      = "PSU-1"
}

output "power_outlet_by_device_and_name" {
  value       = data.netbox_power_outlet.by_device_and_name.type
  description = "Power outlet type when looked up by device and name"
}
