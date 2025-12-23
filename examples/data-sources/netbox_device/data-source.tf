# Example: Look up a device by ID
data "netbox_device" "by_id" {
  id = "1"
}

# Example: Look up a device by name (may return multiple results)
data "netbox_device" "by_name" {
  name = "test-device"
}

# Example: Look up a device by serial number (unique, preferred for uniqueness)
data "netbox_device" "by_serial" {
  serial = "ABC123456789"
}

# Example: Use device data in other resources
output "device_id" {
  value = data.netbox_device.by_id.id
}

output "device_name" {
  value = data.netbox_device.by_name.name
}

output "device_type" {
  value = data.netbox_device.by_serial.device_type
}

output "device_site" {
  value = data.netbox_device.by_serial.site
}
