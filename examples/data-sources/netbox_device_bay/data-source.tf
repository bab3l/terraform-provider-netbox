# Example: Look up a device bay by ID
data "netbox_device_bay" "by_id" {
  id = "1"
}

# Example: Look up a device bay by device and name
# This is useful when you have the parent device and bay name
data "netbox_device_bay" "by_device_and_name" {
  device = "5" # Device ID
  name   = "Bay 1"
}

# Example: Use device bay data in other resources
output "device_bay_id" {
  value = data.netbox_device_bay.by_id.id
}

output "device_bay_name" {
  value = data.netbox_device_bay.by_device_and_name.name
}

output "device_bay_device" {
  value = data.netbox_device_bay.by_device_and_name.device
}
