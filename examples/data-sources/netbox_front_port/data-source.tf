# Example: Look up a front port by ID
data "netbox_front_port" "by_id" {
  id = 1
}

# Example: Look up a front port by device_id and name
data "netbox_front_port" "by_device_and_name" {
  device_id = 5
  name      = "eth0"
}

# Example: Use front port data in other resources
output "front_port_id" {
  value = data.netbox_front_port.by_id.id
}

output "front_port_name" {
  value = data.netbox_front_port.by_device_and_name.name
}

output "front_port_device" {
  value = data.netbox_front_port.by_device_and_name.device
}
